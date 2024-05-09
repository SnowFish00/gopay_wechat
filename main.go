package main

import (
	"encoding/json"
	"fmt"
	test "pay/Test"
	config_viper "pay/config/vipper"
	"pay/global"
	log "pay/log/zap"
	"pay/mysql"
	"pay/router"
)

func init() {
	//使用docker时解除这个⬇️
	// global.SetCfg(config_viper.ConfigWithDocker())
	//使用docker时注释掉这个⬇️
	global.SetCfg(config_viper.Config())
	global.SetDB(mysql.Mysql())
	global.SetLogger(log.InitZap())
	global.SetClient(global.ReturnCfg().NewClientV3Engine())
}

func Test() {
	//离线充值测试
	test.AddTest()
	//离线扣费测试
	test.ReduceTest()
}

func PayOnlineTest() {
	//在线充值订单构造测试
	payID, out, err := test.TestAppltPay()
	if err != nil {
		fmt.Println("错误请求")
		fmt.Println(err.Error())
	}
	fmt.Println(payID)

	//在线充值参数构造测试
	parms, err := global.ReturnClient().PaySignOfApplet(global.ReturnCfg().WxClient.AppID, payID)
	if err != nil {
		fmt.Println("错误请求")
		fmt.Println(err.Error())
	}

	result := make(map[string]interface{})
	result["AppletParams"] = parms
	result["outTradeNo"] = out

	// 将结构体实例转换为JSON格式的字符串
	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Println("转换为JSON时出错:", err)
		return
	}
	fmt.Println(string(jsonData))

}

func main() {
	//路由启动
	router.InitRouter()
}
