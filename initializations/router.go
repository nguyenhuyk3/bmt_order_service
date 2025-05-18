package initializations

import (
	"bmt_order_service/internal/routers"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	r := gin.Default()

	// Routers
	orderRouter := routers.OrderServiceRouterGroup.Order

	mainGroup := r.Group("/v1")
	{
		orderRouter.InitOrderRouter(mainGroup)
	}

	return r
}
