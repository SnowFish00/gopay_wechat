package router

import (
	test "pay/Test"
	"pay/router_basic/wxpay_web"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()
	r.MaxMultipartMemory = 5 << 20

	pay := r.Group("pay")
	{
		pay.POST("OrderBegin", wxpay_web.StartOrder)
		pay.POST("PayNotify", wxpay_web.PayNotify)
		pay.POST("ReduceNotify", wxpay_web.ReduceNotify)
		pay.POST("SearchOrder", wxpay_web.SearchOrder)

	}

	r.Run(":8080")
}

func NotrifyRouter() {
	r := gin.Default()
	r.MaxMultipartMemory = 5 << 20

	api := r.Group("api")
	{
		api.POST("nortify", test.TestPaysigin)
	}

	r.Run(":3636")

}
