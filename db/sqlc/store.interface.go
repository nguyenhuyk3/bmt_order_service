package sqlc

import (
	"bmt_order_service/dto/request"
	"context"
)

type IStore interface {
	CreateOrderTran(ctx context.Context, arg request.Order) (int32, error)
	CreateSubOrderTran(ctx context.Context, arg request.SubOrder, isSuccess bool) error
}
