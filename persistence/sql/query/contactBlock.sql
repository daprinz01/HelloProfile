-- name: GetContactBlock :one
select * from contact_block where id=$1 limit 1;

-- name: AddContactBlock :one
insert into contact_block(
    phone,
    email,
    "address",
    website 
)VALUES(
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateContactBlock :exec
update contact_block set phone=$1,
    email=$2,
    "address"=$3,
    website=$4 where id=$5;

-- name: DeleteContactBlock :exec
delete from contact_block where id=$1;

-- name: UpdateProfileWithContactBlockId :exec
update profiles set contact_block_id=$1 where id=$2;