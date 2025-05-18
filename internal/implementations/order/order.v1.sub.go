package order

import (
	"bmt_order_service/dto/request"
	"bmt_order_service/dto/response"
	"fmt"
)

func (o *orderService) validateSeats(orderSeats []request.OrderSeatReq, availableSeats response.ShowtimeSeats) error {
	seatStatusMap := make(map[int32]string)
	for _, s := range availableSeats.Seats {
		seatStatusMap[s.SeatID] = s.Status
	}

	for _, seat := range orderSeats {
		status, exists := seatStatusMap[seat.SeatId]
		if !exists {
			return fmt.Errorf("seat_id %d does not exist in the showtime", seat.SeatId)
		}
		if status != "available" {
			return fmt.Errorf("seat_id %d is not available (current status: %s)", seat.SeatId, status)
		}
	}

	return nil
}
