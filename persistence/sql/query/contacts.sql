-- name: GetAllContacts :many
SELECT
    *
FROM
    contacts;

-- name: GetContacts :many
SELECT
    *
FROM
    contacts
WHERE
    user_id = $1;

-- name: AddContacts :one
INSERT INTO contacts (user_id, profile_id, contact_category_id)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: UpdateContact :exec
update contacts set contact_category_id=$1 where user_id=$2 and profile_id=$3;

-- name: DeleteContact :exec
delete from contacts where user_id=$1 and profile_id=$2;
