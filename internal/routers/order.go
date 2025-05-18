package routers

import (
	"bmt_order_service/db/sqlc"
	"bmt_order_service/global"
	"bmt_order_service/internal/controllers"
	"bmt_order_service/internal/implementations/order"
	"bmt_order_service/internal/implementations/redis"
	"bmt_order_service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type OrderRouter struct {
}

func (o *OrderRouter) InitOrderRouter(router *gin.RouterGroup) {
	sqlStore := sqlc.NewStore(global.Postgresql)
	redisClient := redis.NewRedisClient()
	orderService := order.NewOrderService(sqlStore, redisClient)
	orderController := controllers.NewOrderController(orderService)
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
