-- name: CreateOrder :one
INSERT INTO orders (
 id,
 username,
 car_id,
 status,
 amount
) VALUES (
  $1, $2 ,$3, $4, $5
)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: ListOrder :many
SELECT * FROM orders
ORDER BY id
LIMIT $1
OFFSET $2; 

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: UpdateOrder :one
UPDATE orders
SET 
  status = COALESCE(sqlc.narg(status),status)
WHERE id = sqlc.arg(id)
RETURNING *;