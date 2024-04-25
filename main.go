package main

import (
	"fmt"
	test "pay/Test"
	config_viper "pay/config/vipper"
	"pay/global"
	log "pay/log/zap"
	"pay/mysql"
)

func init() {
	global.SetCfg(config_viper.Config())
	global.SetDB(mysql.Mysql())
	global.SetLogger(log.InitZap())
	//等待业务接通
	global.SetClient(global.ReturnCfg().NewClientV3Engine())
}

func main() {
	//离线充值测试
	// test.AddTest()
	//离线扣费测试
	// test.ReduceTest()

	//在线充值订单构造测试
	payID, err := test.TestAppltPay()
	if err != nil {
		fmt.Println("错误请求")
		fmt.Println(err.Error())
	}
	fmt.Println(payID)

	//在线充值参数构造测试
	parms, err := global.ReturnClient().PaySignOfApp(global.ReturnCfg().WxClient.AppID, payID)
	if err != nil {
		fmt.Println("错误请求")
		fmt.Println(err.Error())
	}
	fmt.Println(parms)

	//路由启动
	// router.InitRouter()
}
