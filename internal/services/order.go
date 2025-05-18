package services

import (
	"bmt_order_service/dto/request"
	"context"
)

type IOrder interface {
	CreateOrder(ctx context.Context, arg request.Order) (int, error)
}
