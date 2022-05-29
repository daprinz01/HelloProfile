-- name: GetAllProfiles :many
select * from profiles;

-- name: GetProfiles :many
select * from profiles where user_id=$1;

-- name: GetProfile :one
select * from profiles where id=$1 limit 1;

-- name: IsProfileExist :one
select exists(select 1 from profiles where id=$1) AS "exists";

-- name: AddProfile :one
insert into profiles(
    user_id,
    "status",
    profile_name,
    basic_block_id,
    contact_block_id,
    page_color,
    font,
    is_default
) VALUES(
    $1, $2, $3, $4, $5, $6, $7, $8
)
returning *;

-- name: UpdateProfile :exec
update profiles set
   user_id = $1,
    "status" = $2,
    profile_name = $3,
    basic_block_id = $4,
    contact_block_id = $5,
    page_color = $6,
    font = $7,
    is_default= $8 where id=$9;

-- name: DeleteProfile :exec
delete from profiles where id=$1;

-- name: UpdateProfileUrl :exec
update profiles set url=$1 where id=$2;

-- name: IsUrlExists :one
select exists(select 1 from profiles where "url"=$1) AS "exists";