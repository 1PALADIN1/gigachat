package main

import (
	"log"
	"os"

	app "github.com/1PALADIN1/gigachat_server/internal"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	config := new(app.Config)
	//server
	config.Server.Port = viper.GetInt("server.port")
	config.Server.ReadTimeout = viper.GetInt("server.read-timeout")
	config.Server.WriteTimeout = viper.GetInt("server.write-timeout")
	//auth
	config.Auth.SigningKey = os.Getenv("SINGING_KEY")
	config.Auth.PasswordHashSalt = os.Getenv("PASSWORD_HASH_SALT")
	config.Auth.TokenTTL = viper.GetInt("auth.token-ttl")
	//db
	config.DB.DSN = os.Getenv("DSN")
	config.DB.ConnectionTimeout = viper.GetInt("db.connection-timeout")
	//app
	config.App.MinSearchSymbols = viper.GetInt("app.min-search-symb")

	app.Run(config)
}

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
