-- name: CreateOrder :one
INSERT INTO orders (
    ordered_by,
    showtime_id,
    show_date,
    status,
    note
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id;
