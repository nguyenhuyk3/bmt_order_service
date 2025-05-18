package request

type OrderSeatReq struct {
	SeatId int32 `json:"seat_id" binding:"required"`
}

type OrderFAB struct {
	FABId    int32 `json:"fab_id" binding:"required"`
	Quantity int   `json:"quantity" binding:"required"`
}

type Order struct {
	OrderedBy  string         `json:"ordered_by" binding:"required"`
	ShowtimeId int32          `json:"showtime_id" binding:"required"`
	ShowDate   string         `json:"show_date" binding:"required"`
	Note       string         `json:"note"`
	Seats      []OrderSeatReq `json:"seats" binding:"required"`
	FABs       []OrderFAB     `json:"fab"`
}
