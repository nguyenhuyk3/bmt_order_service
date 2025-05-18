-- name: CreateOrderSeat :exec
INSERT INTO order_seats (
    order_id,
    seat_id,
    price
) VALUES (
    $1, $2, $3
);
