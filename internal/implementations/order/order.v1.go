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

	"product"
	"showtime"
)

type orderService struct {
	SqlStore       sqlc.IStore
	RedisClient    services.IRedis
	ProductClient  product.ProductClient
	ShowtimeClient showtime.ShowtimeClient
}

const (
	fifteen_minutes = 15
	sixty_minutes   = 60
)

// CreateOrder implements services.IOrder.
func (o *orderService) CreateOrder(ctx context.Context, arg request.Order) (int32, int, error) {
	if len(arg.FABs) != 0 {
		for _, fAB := range arg.FABs {
			// check if fab exists by calling product service (grpc) to check
			_, err := o.ProductClient.CheckFABExist(ctx,
				&product.CheckFABExistReq{
					FABId: fAB.FABId,
				})
			if err != nil {
				if err.Error() == fmt.Sprintf("rpc error: code = Unknown desc = fab with %d doesn't exist", fAB.FABId) {
					return -1, http.StatusNotFound, err
				}

				return -1, http.StatusInternalServerError, err
			}
		}
	}

	var showTimeSeatsRedisKey string = fmt.Sprintf("%s%d::%s", global.SHOWTIME_SEATS, arg.ShowtimeId, arg.ShowDate)
	var showtimeSeats response.ShowtimeSeats

	// get all showtime seats from cache
	// the reason this code is needed is because the user has to select a seat and then create an order
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

	// save this infor for payment service get this
	err = o.RedisClient.Save(
		fmt.Sprintf("%s%d", global.ORDER, orderId),
		request.SubOrder{
			OrderId:    orderId,
			ShowtimeId: arg.ShowtimeId,
			Seats:      arg.Seats,
			FABs:       arg.FABs,
		},
		fifteen_minutes,
	)
	if err != nil {
		return -1, http.StatusInternalServerError, err
	}

	err = o.RedisClient.Delete(showTimeSeatsRedisKey)
	if err != nil {
		return -1, http.StatusInternalServerError, err
	}

	go func() {
		seatIds := []int32{}

		for _, seat := range arg.Seats {
			seatIds = append(seatIds, seat.SeatId)
		}

		informationForTicket, _ := o.ShowtimeClient.GetSomeInformationForTicket(
			context.Background(),
			&showtime.GetSomeInformationForTicketReq{
				ShowtimeId: arg.ShowtimeId,
				SeatIds:    seatIds,
			})
		film, _ := o.ProductClient.GetFilm(context.Background(),
			&product.GetFilmReq{
				FilmId: informationForTicket.FilmId,
			})

		_ = o.RedisClient.Save(
			fmt.Sprintf("%s%d", global.TICKET_INFORMATION, orderId),
			ticketInformation{
				CinemaName: informationForTicket.CinemaName,
				City:       informationForTicket.City,
				Location:   informationForTicket.Location,
				RoomName:   informationForTicket.RoomName,
				ShowDate:   informationForTicket.ShowDate,
				StartTime:  informationForTicket.StartTime,
				Seats:      informationForTicket.Seats,
				Genres:     film.Genres,
				FilmPoster: film.PosterUrl,
				Title:      film.Title,
				Duration:   film.Duration,
			}, sixty_minutes)
	}()

	return orderId, http.StatusOK, nil
}

func NewOrderService(
	sqlStore sqlc.IStore,
	redisClient services.IRedis,
	productClient product.ProductClient,
	showtimeClient showtime.ShowtimeClient,
) services.IOrder {
	return &orderService{
		SqlStore:       sqlStore,
		RedisClient:    redisClient,
		ProductClient:  productClient,
		ShowtimeClient: showtimeClient,
	}
}
