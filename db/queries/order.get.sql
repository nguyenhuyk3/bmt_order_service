-- name: GetOrderByTicketBooker :many
SELECT * 
FROM orders
WHERE ordered_by = $1;