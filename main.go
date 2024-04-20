package main

import (
	config_viper "pay/config/vipper"
	"pay/global"
	log "pay/log/zap"
	"pay/mysql"
	"pay/router"
)

func init() {
	global.SetCfg(config_viper.Config())
	global.SetDB(mysql.Mysql())
	global.SetLogger(log.InitZap())
	//等待业务接通
	// global.SetClient(global.ReturnCfg().NewClientV3Engine())
}

func main() {
	// test.AddTest()
	router.InitRouter()
}
