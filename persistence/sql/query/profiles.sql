-- name: GetAllProfiles :many
select * from profiles;

-- name: GetProfiles :many
select * from profiles where user_id=$1;

-- name: GetProfile :one
select * from profiles where id=$1 limit 1;

-- name: AddProfile :one
insert into profiles(
    user_id
    status,
    profile_name,
    fullname,
    title,
    bio,
    company,
    company_address,
    image_url,
    phone,
    email,
    address_id uuid,
    website,
    is_default,
    color
) VALUES(
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
returning *;

-- name: UpdateProfile :exec
update profiles set
    status = $1,
    profile_name = $2,
    fullname = $3,
    title = $4,
    bio = $5,
    company = $6,
    company_address = $7,
    image_url = $8,
    phone = #9,
    email = $10,
    address_id uuid = $11,
    website = $12,
    is_default = $13,
    color = $14 where id=$15;

-- name: DeleteProfile :exec
delete profiles where id=$1;