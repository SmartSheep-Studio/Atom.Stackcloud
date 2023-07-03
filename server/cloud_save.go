package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"repo.smartsheep.studio/atom/matrix/datasource/models"
)

type CloudSaveController struct {
	db *gorm.DB
}

func NewCloudSaveController(db *gorm.DB) *CloudSaveController {
	return CloudSaveController{db}
}

func (ctrl *CloudSaveController) Map(router *fiber.App) {

}

func (ctrl *CloudSaveController) get(c *fiber.Ctx) error {
	u := c.Locals("matrix-id").(*models.MatrixProfile)
}
