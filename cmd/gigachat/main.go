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
	config.DB.Host = viper.GetString("db.host")
	config.DB.Port = viper.GetInt("db.port")
	config.DB.User = viper.GetString("db.user")
	config.DB.Password = os.Getenv("DB_PASSWORD")
	config.DB.DBName = viper.GetString("db.db-name")
	config.DB.SSLMode = viper.GetString("db.ssl-mode")

	app.Run(config)
}

func initConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
