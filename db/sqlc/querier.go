// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"context"
)

type Querier interface {
	CreateOrder(ctx context.Context, arg CreateOrderParams) (int32, error)
	CreateOrderFAB(ctx context.Context, arg CreateOrderFABParams) error
	CreateOrderSeat(ctx context.Context, arg CreateOrderSeatParams) error
	CreateOutbox(ctx context.Context, arg CreateOutboxParams) error
	GetOrderByTicketBooker(ctx context.Context, orderedBy string) ([]Orders, error)
	UpdateOrderStatusByOrderId(ctx context.Context, arg UpdateOrderStatusByOrderIdParams) error
}

var _ Querier = (*Queries)(nil)
