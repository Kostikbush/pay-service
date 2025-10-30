package ports

import (
	"github.com/gin-gonic/gin"
)

func Routers(routerGroup *gin.RouterGroup, handler *Handler) {
	routerGroup.POST("/pay/init/", handler.InitPay)
}