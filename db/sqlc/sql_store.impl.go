package sqlc

import (
	"bmt_order_service/dto/request"
	"bmt_order_service/global"
	"bmt_order_service/utils/convertors"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SqlStore struct {
	connPool *pgxpool.Pool
}

// CreateOrderTran implements IStore.
func (s *SqlStore) CreateOrderTran(ctx context.Context, arg request.Order) (int32, error) {
	var finalOrderId int32 = -1

	showDate, err := convertors.ConvertDateStringToTime(arg.ShowDate)
	if err != nil {
		return finalOrderId, err
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
				Status: OrderStatusesCreated,
				Note:   arg.Note,
			})
		if err != nil {
			return fmt.Errorf("failed to create order with showtime id (%d): %w", arg.ShowtimeId, err)
		} else {
			finalOrderId = orderId
		}

		payloadBytes, err := json.Marshal(gin.H{
			"order_id":    orderId,
			"showtime_id": arg.ShowtimeId,
			"ordered_by":  arg.OrderedBy,
			"seats":       arg.Seats,
			"fabs":        arg.FABs,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		/**
			this message will be received by Showtime Service
		to change seat status available -> reserved
		*/
		err = q.CreateOutbox(ctx,
			CreateOutboxParams{
				AggregatedType: "ORDER_ID",
				AggregatedID:   orderId,
				EventType:      global.ORDER_CREATED,
				Payload:        payloadBytes,
			})
		if err != nil {
			return fmt.Errorf("failed to create outbox (create order): %w", err)
		}

		return nil
	})

	if err != nil {
		return -1, err
	}

	return finalOrderId, nil
}

// CreateSubOrderTran implements IStore.
func (s *SqlStore) CreateSubOrderTran(ctx context.Context, arg request.SubOrder, isSuccess bool) error {
	return s.execTran(ctx, func(q *Queries) error {
		eventType := global.ORDER_FAILED
		status := OrderStatusesFailed
		if isSuccess {
			status = OrderStatusesSuccess
			eventType = global.ORDER_SUCCESS
		}

		if err := q.UpdateOrderStatusByOrderId(ctx,
			UpdateOrderStatusByOrderIdParams{
				ID:     arg.OrderId,
				Status: status,
			}); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		// if payment failed then just perform this step
		if !isSuccess {
			payloadBytes, err := json.Marshal(arg)
			if err != nil {
				return fmt.Errorf("failed to marshal payload: %w", err)
			}

			// write messagw to outboxes table for showtime service reading this
			/*
					because in this situation due to payment failure
				then we will send ORDER_FAILED event type for showtime service to showtime service change seat status
				from 'reserved' -> 'available'
			*/
			err = q.CreateOutbox(ctx,
				CreateOutboxParams{
					AggregatedType: "ORDER_ID",
					AggregatedID:   arg.OrderId,
					EventType:      eventType,
					Payload:        payloadBytes,
				})
			if err != nil {
				return fmt.Errorf("failed to create outbox (create sub order): %w", err)
			}

			return nil
		}

		// create order for seats
		for _, seat := range arg.Seats {
			if err := q.CreateOrderSeat(ctx,
				CreateOrderSeatParams{
					OrderID: arg.OrderId,
					SeatID:  seat.SeatId,
				}); err != nil {
				return fmt.Errorf("failed to create seat order with id (%d): %w", seat.SeatId, err)
			}
		}

		// create order for fabs
		for _, fab := range arg.FABs {
			if err := q.CreateOrderFAB(ctx, CreateOrderFABParams{
				OrderID:  arg.OrderId,
				FabID:    fab.FABId,
				Quantity: int32(fab.Quantity),
			}); err != nil {
				return fmt.Errorf("failed to create fab order with id (%d): %w", fab.FABId, err)
			}
		}

		payloadBytes, err := json.Marshal(arg)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		/**
			this message will be received by Showtime Service
		to change seat status reserved -> available or booked based on envetType
		*/
		err = q.CreateOutbox(ctx,
			CreateOutboxParams{
				AggregatedType: "ORDER_ID",
				AggregatedID:   arg.OrderId,
				EventType:      eventType,
				Payload:        payloadBytes,
			})
		if err != nil {
			return fmt.Errorf("failed to create outbox (create sub order): %w", err)
		}

		return nil
	})
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
