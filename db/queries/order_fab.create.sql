-- name: CreateOrderFAB :exec
INSERT INTO order_fabs (
    order_id,
    fab_id,
    quantity,
    price
) VALUES (
    $1, $2, $3, $4
);
