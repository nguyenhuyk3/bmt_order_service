//go:build wireinject

package injectors

import (
	"bmt_order_service/internal/controllers"
	"bmt_order_service/internal/implementations/order"

	"github.com/google/wire"
)

func InitOrderController() (*controllers.OrderController, error) {
	wire.Build(
		dbSet,
		redisSet,

		order.NewOrderService,
		controllers.NewOrderController,
	)

	return nil, nil
}
