package config_viper

import (
	"log"
	cfg "pay/model/config_model"

	"github.com/spf13/viper"
)

func Config() cfg.Config {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("读取配置文件错误", err)
	}

	var config cfg.Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Unmarshal config failed: ", err)
	}

	return config

}
