-- name: CreatePayment :one
INSERT INTO payments (
  id,
  order_id,
  username,
  network,
  currency,
  amount,
  wallet_address,
  tx_hash,
  status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1;

-- name: GetPaymentByTxHash :one
SELECT * FROM payments
WHERE tx_hash = $1;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET status = $2,
    tx_hash = COALESCE($3, tx_hash),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: ListPaymentsByUser :many
SELECT * FROM payments
WHERE username = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
