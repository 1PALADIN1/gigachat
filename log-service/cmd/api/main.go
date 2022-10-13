package main

import (
	"log"
	"os"

	app "github.com/1PALADIN1/gigachat_server/log/internal"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	config := new(app.Config)
	// server
	config.Server.GRPCPort = viper.GetInt("server.grpc-port")
	// db
	config.DB.DSN = os.Getenv("DSN")
	config.DB.ConnectionTimeout = viper.GetInt("db.connection-timeout")

	app.Run(config)
}

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
