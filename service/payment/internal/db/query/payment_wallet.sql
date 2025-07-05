-- name: CreatePaymentWallet :one
INSERT INTO payment_wallets (
  network,
  wallet_address
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetWalletAddressByNetwork :one
SELECT * FROM payment_wallets
WHERE network = $1
LIMIT 1;

-- name: ListPaymentWallets :many
SELECT * FROM payment_wallets
ORDER BY network;

-- name: UpdatePaymentWallet :one
UPDATE payment_wallets
SET wallet_address = $2
WHERE network = $1
RETURNING *;

-- name: DeletePaymentWallet :exec
DELETE FROM payment_wallets
WHERE network = $1;
