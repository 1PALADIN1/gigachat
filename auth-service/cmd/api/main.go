package main

import (
	"log"
	"os"

	app "github.com/1PALADIN1/gigachat_server/auth/internal"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	config := new(app.Config)

	// server
	config.Server.GRPCPort = viper.GetInt("server.grpc-port")
	// auth
	config.Auth.SigningKey = os.Getenv("SINGING_KEY")
	config.Auth.PasswordHashSalt = os.Getenv("PASSWORD_HASH_SALT")
	config.Auth.TokenTTL = viper.GetInt("auth.token-ttl")
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
