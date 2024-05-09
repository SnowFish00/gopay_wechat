package config_viper

import (
	"fmt"
	"io/ioutil"
	"log"
	cfg "pay/model/config_model"

	"github.com/spf13/viper"
)

func ConfigWithDocker() cfg.Config {
	viper.AddConfigPath("/root/config")
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

	// 读取.pem文件
	pemFile := "/root/config/com.pem"
	pemData, err := ioutil.ReadFile(pemFile)
	if err != nil {
		fmt.Println("读取.pem文件时出错:", err)
		return cfg.Config{}
	}

	config.WxClient.PrivateKey = string(pemData)

	return config

}
