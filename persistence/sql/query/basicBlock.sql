-- name: GetBasicBlock :one
select * from basic_block where id=$1 limit 1;

-- name: AddBasicBlock :one
insert into basic_block(
    profile_photo_url,
    cover_photo_url,
    cover_colour,
    fullname,
    title,
    bio
)VALUES(
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateBasicBlock :exec
update basic_block set profile_photo_url=$1,
    cover_photo_url=$2,
    cover_colour=$3,
    fullname=$4,
    title=$5,
    bio=$6 where id=$7;

-- name: DeleteBasicBlock :exec
delete from basic_block where id=$1;
