-- name: GetUserLogins :many
select * from user_login;

-- name: GetUserLogin :many
select * from user_login where user_id = $1;

-- name: GetUnResoledLogins :many
select * from user_login where user_id = $1 and resolved = false;

-- name: CreateUserLogin :one
insert into user_login (
    user_id,
    application_id,
    login_status,
    response_code,
    response_description,
    device,
    ip_address,
    longitude,
    latitude,
    resolved
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
returning *;

-- name: UpdateResolvedLogin :exec
update user_login set resolved = true where user_id = $1 and resolved = false;

-- name: DeleteUserLogin :exec
delete from user_login where user_id = $1 and application_id = $2;