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
	config.Auth.Addr = viper.GetString("auth.addr")
	config.Auth.ConnTimeout = viper.GetInt("auth.conn-timeout")
	// log
	config.Log.Addr = viper.GetString("log.addr")
	config.Log.ConnTimeout = viper.GetInt("log.conn-timeout")
	config.Log.Source = viper.GetString("log.source")
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
