-- name: CreateOrderFAB :exec
INSERT INTO order_fabs (
    order_id,
    fab_id,
    quantity
) VALUES (
    $1, $2, $3
);
