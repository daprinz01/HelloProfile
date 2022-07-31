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
INSERT INTO contacts (user_id, profile_id)
    VALUES ($1, $2)
RETURNING
    *;



-- name: DeleteContact :exec
delete from contacts where user_id=$1 and profile_id=$2;

-- name: DeleteProfileContact :exec
delete from contacts where profile_id=$1;
