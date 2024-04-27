package ping

import (
	responses "pay/response"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	responses.OkWithMessage("pong", c)
}
