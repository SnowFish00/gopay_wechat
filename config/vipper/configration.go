package config_viper

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	//读取pem
	filePath, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前文件路径时出错:", err)
		return cfg.Config{}
	}

	// 读取.pem文件
	pemFile := filePath + "/config/com.pem"
	pemData, err := ioutil.ReadFile(pemFile)
	if err != nil {
		fmt.Println("读取.pem文件时出错:", err)
		return cfg.Config{}
	}

	config.WxClient.PrivateKey = string(pemData)

	return config

}
