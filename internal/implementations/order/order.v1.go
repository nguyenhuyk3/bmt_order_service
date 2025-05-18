package order

import (
	"bmt_order_service/db/sqlc"
	"bmt_order_service/dto/request"
	"bmt_order_service/dto/response"
	"bmt_order_service/global"
	"bmt_order_service/internal/services"
	"context"
	"fmt"
	"net/http"
)

type orderService struct {
	SqlStore    sqlc.IStore
	RedisClient services.IRedis
}

// CreateOrder implements services.IOrder.
func (o *orderService) CreateOrder(ctx context.Context, arg request.Order) (int, error) {
	var showTimeSeatsRedisKey string = fmt.Sprintf("%s%d::%s", global.SHOWTIME_SEATS, arg.ShowtimeId, arg.ShowDate)
	var showtimeSeats response.ShowtimeSeats

	err := o.RedisClient.Get(showTimeSeatsRedisKey, &showtimeSeats)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("showtime id (%d) and show date (%s) not foud with err: %v", arg.ShowtimeId, arg.ShowDate, err)
	}

	if err = o.validateSeats(arg.Seats, showtimeSeats); err != nil {
		return http.StatusBadRequest, err
	}

	err = o.SqlStore.CreateOrderTran(ctx, arg)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func NewOrderService(
	sqlStore sqlc.IStore,
	redisClient services.IRedis,
) services.IOrder {
	return &orderService{
		SqlStore:    sqlStore,
		RedisClient: redisClient,
	}
}
