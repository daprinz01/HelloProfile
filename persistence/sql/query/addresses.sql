-- name: GetAllAddresses :many
SELECT
    *
FROM
    addresses;

-- name: GetUserAddresses :many
SELECT
    *
FROM
    addresses
WHERE
    user_id = $1;

-- name: GetAddress :one
SELECT
    *
FROM
    addresses
WHERE
    id = $1
LIMIT 1;

-- name: AddAddress :one
INSERT INTO addresses (user_id, street, city, state, country)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: UpdateAddress :exec
UPDATE
    addresses
SET
    street = $2,
    city = $3,
    state = $4,
    country = $5
WHERE
    id = $1;

-- name: DeleteAddress :exec
DELETE FROM addresses
WHERE id = $1;

