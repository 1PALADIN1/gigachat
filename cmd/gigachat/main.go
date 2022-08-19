package main

import (
	"log"
	"os"

	app "github.com/1PALADIN1/gigachat_server/internal"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	config := app.Config{}
	config.Server.Port = viper.GetInt("server.port")
	config.Server.ReadTimeout = viper.GetInt("server.read-timeout")
	config.Server.WriteTimeout = viper.GetInt("server.write-timeout")
	config.Auth.SigningKey = os.Getenv("SINGING_KEY")

	app.Run(config)
}

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
