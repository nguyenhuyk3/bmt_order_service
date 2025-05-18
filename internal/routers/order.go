package routers

import (
	"bmt_order_service/internal/injectors"
	"bmt_order_service/internal/middlewares"
	"log"

	"github.com/gin-gonic/gin"
)

type OrderRouter struct {
}

func (o *OrderRouter) InitOrderRouter(router *gin.RouterGroup) {
	orderController, err := injectors.InitOrderController()
	if err != nil {
		log.Fatalf("failed to initialize OrderController: %v", err)
		return
	}

	getFromHeaderMiddleware := middlewares.NewGetFromHeaderMiddleware()

	orderRouter := router.Group("/order")
	{
		orderPublicRouter := orderRouter.Group("/public")
		{
			orderPublicRouter.POST("/create",
				getFromHeaderMiddleware.GetEmailFromHeader(),
				orderController.CreateOrder)
		}
	}
}
