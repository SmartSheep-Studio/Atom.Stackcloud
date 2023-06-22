package server

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"repo.smartsheep.studio/atom/matrix/datasource/models"
	"repo.smartsheep.studio/atom/matrix/server/middleware"
	"repo.smartsheep.studio/atom/nucleus/toolbox"
	"repo.smartsheep.studio/atom/nucleus/utils"
)

type LibraryController struct {
	db   *gorm.DB
	conn *toolbox.ExternalServiceConnection
	auth middleware.AuthHandler
}

func NewLibraryController(db *gorm.DB, conn *toolbox.ExternalServiceConnection, auth middleware.AuthHandler) *LibraryController {
	return &LibraryController{db, conn, auth}
}

func (ctrl *LibraryController) Map(router *fiber.App) {
	router.Get("/api/library", ctrl.auth(true), ctrl.list)
	router.Post("/api/library/add", ctrl.auth(true), ctrl.add)
	router.Get("/api/library/add/callback", ctrl.callback)
}

func (ctrl *LibraryController) list(c *fiber.Ctx) error {
	u := c.Locals("matrix-prof").(*models.MatrixProfile)

	var items []models.MatrixLibraryItem
	if err := ctrl.db.Where("profile_id = ?", u.ID).Find(&items).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *LibraryController) add(c *fiber.Ctx) error {
	u := c.Locals("matrix-prof").(*models.MatrixProfile)

	var req struct {
		App string `json:"app" validate:"required"`
	}

	if err := utils.ParseRequestBody(c, &req); err != nil {
		return err
	}

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ?", req.App).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	url, err := ctrl.conn.GetExternalServiceEndpointPath(
		"repo.smartsheep.studio/atom/quarkpay",
		fmt.Sprintf("/api/shops/%s/transactions", app.PriceOptions.Data().Shop),
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	client := resty.New()
	res, err := client.R().SetAuthToken(app.PriceOptions.Data().ApiToken).SetBody(fiber.Map{
		"description": "Matrix Marketplace Checkout",
		"items": []fiber.Map{
			{"id": app.PriceOptions.Data().ProductID, "amount": 1},
		},
	}).SetResult(&models.QuarkTransaction{}).Post(url)

	if err != nil || res.StatusCode() != fiber.StatusOK {
		fmt.Println(err, string(res.Body()))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create order on quarkpay")
	}

	transaction := res.Result().(*models.QuarkTransaction)
	url, _ = ctrl.conn.GetExternalServiceEndpointPath(
		"repo.smartsheep.studio/atom/quarkpay",
		fmt.Sprintf("/checkout?id=%d&redirect_uri=%s/api/library/add/callback", transaction.ID, viper.GetString("general.base_url")),
	)

	order := models.MatrixTransaction{
		AppID:      app.ID,
		ProfileID:  u.ID,
		QuarkpayID: transaction.ID,
	}

	if err := ctrl.db.Save(&order).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	return c.JSON(fiber.Map{
		"url":         url,
		"order":       order,
		"transaction": transaction,
	})
}

func (ctrl *LibraryController) callback(c *fiber.Ctx) error {
	id := c.Query("transaction_id")

	var order models.MatrixTransaction
	if err := ctrl.db.Where("quarkpay_id = ?", id).First(&order).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	var app models.MatrixApp
	if err := ctrl.db.Where("id = ?", order.AppID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	url, err := ctrl.conn.GetExternalServiceEndpointPath(
		"repo.smartsheep.studio/atom/quarkpay",
		fmt.Sprintf("/api/transactions/%s", id),
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	client := resty.New()
	res, err := client.R().SetAuthToken(app.PriceOptions.Data().ApiToken).SetBody(fiber.Map{
		"description": "Matrix Marketplace Checkout",
		"items": []fiber.Map{
			{"id": app.PriceOptions.Data().ProductID, "amount": 1},
		},
	}).SetResult(&models.QuarkTransaction{}).Get(url)

	if err != nil || res.StatusCode() != fiber.StatusOK {
		fmt.Println(err, string(res.Body()))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get order from quarkpay")
	}

	transaction := res.Result().(*models.QuarkTransaction)

	if !transaction.IsFinished {
		return fiber.NewError(fiber.StatusForbidden, "order in quarkpay isn't finish yet!")
	} else {
		item := models.MatrixLibraryItem{
			ProfileID: order.ProfileID,
			AppID:     order.AppID,
		}

		if err := ctrl.db.Save(&item).Error; err != nil {
			return utils.ParseDataSourceError(err)
		} else {
			return c.Redirect(fmt.Sprintf("%s/library", viper.GetString("general.base_url")))
		}
	}
}
