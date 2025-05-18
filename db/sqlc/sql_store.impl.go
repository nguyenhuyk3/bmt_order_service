package sqlc

import (
	"bmt_order_service/dto/request"
	"bmt_order_service/utils/convertors"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SqlStore struct {
	connPool *pgxpool.Pool
}

// CreateOrderTran implements IStore.
func (s *SqlStore) CreateOrderTran(ctx context.Context, arg request.Order) error {
	showDate, err := convertors.ConvertDateStringToTime(arg.ShowDate)
	if err != nil {
		return err
	}

	err = s.execTran(ctx, func(q *Queries) error {
		orderId, err := q.CreateOrder(ctx,
			CreateOrderParams{
				OrderedBy:  arg.OrderedBy,
				ShowtimeID: arg.ShowtimeId,
				ShowDate: pgtype.Date{
					Time:  showDate,
					Valid: true,
				},
				Status: OrderStatusesPending,
				Note:   arg.Note,
			})
		if err != nil {
			return fmt.Errorf("failed to create order with showtime id (%d): %w", arg.ShowtimeId, err)
		}

		for _, seat := range arg.Seats {
			err = q.CreateOrderSeat(ctx,
				CreateOrderSeatParams{
					OrderID: orderId,
					SeatID:  seat.SeatId,
				})
			if err != nil {
				return fmt.Errorf("failed to create seat order with id (%d): %v", seat.SeatId, err)
			}
		}

		if len(arg.FABs) != 0 {
			for _, fab := range arg.FABs {
				err = q.CreateOrderFAB(ctx,
					CreateOrderFABParams{
						OrderID:  orderId,
						FabID:    fab.FABId,
						Quantity: int32(fab.Quantity),
					})
				if err != nil {
					return fmt.Errorf("failed to create fab order with id (%d): %v", fab.FABId, err)
				}
			}
		}

		return nil
	})

	return err
}

func (s *SqlStore) execTran(ctx context.Context, fn func(*Queries) error) error {
	// Start transaction
	tran, err := s.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tran)
	// fn performs a series of operations down the db
	err = fn(q)
	if err != nil {
		// If an error occurs, rollback the transaction
		if rbErr := tran.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tran err: %v, rollback err: %v", err, rbErr)
		}

		return err
	}

	return tran.Commit(ctx)
}

func NewStore(connPool *pgxpool.Pool) IStore {
	return &SqlStore{
		connPool: connPool,
	}
}
