-- name: GetAllContactCategories :many
select * from contact_categories;

-- name: GetContactCategory :one
select * from contact_categories where id=$1 or name=$1;

-- name: AddContactCategory :one
insert into contact_categories(
    name
)VALUES(
    $1
) returning *;

-- name: UpdateContactCategory :exec
update contact_categories set name=$1 where id=$2;

-- name: DeleteContactCategory :exec
delete from contact_categories where id=$1 or name=$1;