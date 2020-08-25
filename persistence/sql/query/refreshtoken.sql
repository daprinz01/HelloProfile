-- name: GetRefreshTokens :many
select * from refresh_token;

-- name: GetRefreshToken :one
select * from refresh_token where token = $1;

-- name: CreateRefreshToken :one
insert into refresh_token (
    "user_id",
    "token"
    )
    values
    ($1, $2)
    returning *;

-- name: UpdateRefreshToken :one
update refresh_token set "user_id" = $1, "token" = $2 where "token" = $3 returning *;

-- name: DeleteRefreshToken :exec
delete from refresh_token where "token" = $1;

