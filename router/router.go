package router

import (
	gopayapi "pay/gopay_api"
	"pay/router_basic/ping"
	"pay/router_basic/wxpay_web"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()
	r.MaxMultipartMemory = 5 << 20

	r.GET("ping", ping.Ping) //通过

	pay := r.Group("pay")
	{
		pay.POST("OrderBegin", wxpay_web.StartOrder)   //通过
		pay.POST("SearchOrder", wxpay_web.SearchOrder) //通过

	}

	wxOptions := r.Group("options")
	{
		wxOptions.POST("GetOpenId", gopayapi.GetOpenIDBycode2Session)
	}

	api := r.Group("api")
	{
		api.POST("PayNotify", wxpay_web.PayNotify) //通过
	}

	syn := r.Group("syn")
	{
		syn.POST("AddNotrify", wxpay_web.AddNotrify)         //通过
		syn.POST("SnowReduceNotify", wxpay_web.ReduceNotify) //通过
	}

	r.Run(":3636")
}
