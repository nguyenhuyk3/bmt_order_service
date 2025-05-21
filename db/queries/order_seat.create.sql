-- name: CreateOrderSeat :exec
INSERT INTO order_seats (
    order_id,
    seat_id
) VALUES (
    $1, $2
);
