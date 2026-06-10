package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config") // 配置文件在 config 目录下
	viper.AddConfigPath(".")        // 也可以同时添加根目录作为备选
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if loadErr := godotenv.Load(".env"); loadErr != nil {
		log.Printf("config: .env not loaded (%v), using process environment only", loadErr)
	}
	viper.AutomaticEnv()
}
