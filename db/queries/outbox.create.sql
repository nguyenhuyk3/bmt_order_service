-- name: CreateOutbox :exec
INSERT INTO outboxes (
    aggregated_type,
    aggregated_id,
    event_type,
    payload
) VALUES (
    $1, $2, $3, $4
);