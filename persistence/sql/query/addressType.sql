-- name: GetAllAddressTypes :many
select * from address_types;

-- name: GetAddressType :one
select * from address_types where id=$1 or name=$1;

-- name: AddAddressType :one
insert into address_types(
    name
)VALUES(
    $1
) returning *;

-- name: UpdateAddressType :exec
update address_types set name=$1 where id=$2;

-- name: DeleteAddressType :exec
delete from address_types where id=$1;