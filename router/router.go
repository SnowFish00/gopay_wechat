package router

import (
	test "pay/Test"
	"pay/router_basic/ping"
	"pay/router_basic/wxpay_web"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()
	r.MaxMultipartMemory = 5 << 20

	r.GET("ping", ping.Ping)

	pay := r.Group("pay")
	{
		pay.POST("OrderBegin", wxpay_web.StartOrder)
		pay.POST("PayNotify", wxpay_web.PayNotify)
		pay.POST("AddNotrify", wxpay_web.AddNotrify)
		pay.POST("ReduceNotify", wxpay_web.ReduceNotify)
		pay.POST("SearchOrder", wxpay_web.SearchOrder)

	}

	api := r.Group("api")
	{
		api.POST("nortify", test.TestPaysigin)
	}

	r.Run(":3636")
}
