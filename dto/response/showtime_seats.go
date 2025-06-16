package response

import "github.com/jackc/pgx/v5/pgtype"

type showtimeSeat struct {
	Id         int32  `json:"id"`
	ShowtimeID int32  `json:"showtime_id"`
	SeatID     int32  `json:"seat_id"`
	Status     string `json:"status"`
	BookedBy   string `json:"booked_by"`
	// CreatedAt  time.Time  `json:"created_at"`
	BookedAt pgtype.Timestamp `json:"booked_at"`
}

type ShowtimeSeats struct {
	Seats []showtimeSeat
}
