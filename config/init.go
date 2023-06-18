package config

import (
	"os"

	"github.com/spf13/viper"
)

func init() {
	if err := LoadConfig(); err != nil {
		panic(err)
	} else {
		if err := InitFS(); err != nil {
			panic(err)
		}
	}
}

func LoadConfig() error {
	viper.SetConfigName("settings")
	viper.SetConfigType("toml")
	viper.AddConfigPath("..")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func InitFS() error {
	folders := []string{
		viper.GetString("resources.user_contents"),
	}

	for _, folder := range folders {
		stat, err := os.Stat(folder)
		if os.IsNotExist(err) || !stat.IsDir() {
			if err := os.MkdirAll(folder, 0755); err != nil {
				return err
			}
		}
	}

	return nil
}
