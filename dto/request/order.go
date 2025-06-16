package request

type OrderSeatReq struct {
	SeatId int32 `json:"seat_id" binding:"required"`
}

type OrderFABReq struct {
	FABId    int32 `json:"fab_id" binding:"required"`
	Quantity int   `json:"quantity" binding:"required"`
}

type Order struct {
	OrderedBy  string
	ShowtimeId int32          `json:"showtime_id" binding:"required"`
	ShowDate   string         `json:"show_date" binding:"required"`
	Note       string         `json:"note"`
	Seats      []OrderSeatReq `json:"seats" binding:"required"`
	FABs       []OrderFABReq  `json:"fabs"`
}

type SubOrder struct {
	OrderId    int32          `json:"order_id" binding:"required"`
	ShowtimeId int32          `json:"showtime_id" binding:"required"`
	Seats      []OrderSeatReq `json:"seats" binding:"required"`
	FABs       []OrderFABReq  `json:"fabs"`
}
