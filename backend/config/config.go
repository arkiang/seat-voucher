package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Env            string
	Port           string
	FrontendURL    string
	DBPath         string
	SeatLayoutPath string
}

func LoadConfig() Config {
	root := findProjectRootWithEnv()
	viper.SetConfigFile(filepath.Join(root, "app.env"))
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No app.env file found or failed to load, using system env if available.")
	}

	return Config{
		Env:            viper.GetString("ENV"),
		Port:           viper.GetString("PORT"),
		FrontendURL:    viper.GetString("FRONTEND_URL"),
		DBPath:         filepath.Join(root, viper.GetString("DB_PATH")),
		SeatLayoutPath: filepath.Join(root, viper.GetString("SEAT_LAYOUT_PATH")),
	}
}

func findProjectRootWithEnv() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("cannot get working directory:", err)
	}

	for {
		envPath := filepath.Join(dir, "app.env")
		if _, err := os.Stat(envPath); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			log.Fatal("app.env not found in any parent directories")
		}
		dir = parent
	}
}
