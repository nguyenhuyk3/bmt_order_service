-- name: UpdateOrderStatusByOrderId :exec
UPDATE orders
SET 
    status = $2,
    updated_at = NOW()
WHERE id = $1;
