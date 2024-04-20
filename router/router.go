package router

import (
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
		pay.POST("SearchOrder", wxpay_web.SearchOrder)
	}

	r.Run(":8080")
}
