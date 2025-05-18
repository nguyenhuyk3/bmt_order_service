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

	err := o.RedisClient.Get(showTimeSeatsRedisKey, &showtimeSeats.Seats)
	if err != nil {
		if err.Error() == fmt.Sprintf("key %s does not exist", showTimeSeatsRedisKey) {
			return http.StatusNotFound, fmt.Errorf("showtime id (%d) or show date (%s) not foud", arg.ShowtimeId, arg.ShowDate)
		}

		return http.StatusInternalServerError, fmt.Errorf("failed to get value from redis with err: %w", err)
	}

	if status, err := o.validateSeats(arg.Seats, showtimeSeats); err != nil && status != http.StatusOK {
		return status, err
	}

	err = o.SqlStore.CreateOrderTran(ctx, arg)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = o.RedisClient.Delete(showTimeSeatsRedisKey)
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
