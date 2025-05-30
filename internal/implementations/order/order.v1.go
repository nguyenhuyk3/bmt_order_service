package order

import (
	"bmt_order_service/db/sqlc"
	"bmt_order_service/dto/request"
	"bmt_order_service/dto/response"
	"bmt_order_service/global"
	"bmt_order_service/internal/services"
	"context"
	"errors"
	"fmt"
	"net/http"

	"product"
)

type orderService struct {
	SqlStore      sqlc.IStore
	RedisClient   services.IRedis
	ProductClient product.ProductClient
}

const (
	fifteen_minutes = 15
)

// CreateOrder implements services.IOrder.
func (o *orderService) CreateOrder(ctx context.Context, arg request.Order) (int32, int, error) {
	if len(arg.FABs) != 0 {
		for _, fAB := range arg.FABs {
			_, err := o.ProductClient.CheckFABExist(ctx, &product.CheckFABExistReq{FABId: fAB.FABId})
			if err != nil {
				if errors.Is(err, fmt.Errorf("fab with %d doesn't exist", fAB.FABId)) {
					return -1, http.StatusNotFound, err
				}

				return -1, http.StatusInternalServerError, err
			}
		}
	}

	var showTimeSeatsRedisKey string = fmt.Sprintf("%s%d::%s", global.SHOWTIME_SEATS, arg.ShowtimeId, arg.ShowDate)
	var showtimeSeats response.ShowtimeSeats

	err := o.RedisClient.Get(showTimeSeatsRedisKey, &showtimeSeats.Seats)
	if err != nil {
		if err.Error() == fmt.Sprintf("key %s does not exist", showTimeSeatsRedisKey) {
			return -1, http.StatusNotFound, fmt.Errorf("showtime id (%d) or show date (%s) not foud", arg.ShowtimeId, arg.ShowDate)
		}

		return -1, http.StatusInternalServerError, fmt.Errorf("failed to get value from redis with err: %w", err)
	}

	if status, err := o.validateSeats(arg.Seats, showtimeSeats); err != nil && status != http.StatusOK {
		return -1, status, err
	}

	orderId, err := o.SqlStore.CreateOrderTran(ctx, arg)
	if err != nil {
		return -1, http.StatusInternalServerError, err
	}

	var orderRedisKey string = fmt.Sprintf("%s%d", global.ORDER, orderId)

	err = o.RedisClient.Save(orderRedisKey,
		request.SubOrder{
			OrderId:    orderId,
			ShowtimeId: arg.ShowtimeId,
			Seats:      arg.Seats,
			FABs:       arg.FABs,
		},
		fifteen_minutes)
	if err != nil {
		return -1, http.StatusInternalServerError, err
	}

	err = o.RedisClient.Delete(showTimeSeatsRedisKey)
	if err != nil {
		return -1, http.StatusInternalServerError, err
	}

	return orderId, http.StatusOK, nil
}

func NewOrderService(
	sqlStore sqlc.IStore,
	redisClient services.IRedis,
	productClient product.ProductClient,
) services.IOrder {
	return &orderService{
		SqlStore:    sqlStore,
		RedisClient: redisClient,
	}
}
