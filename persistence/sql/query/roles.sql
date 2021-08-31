-- name: GetRoles :many
select * from roles;

-- name: GetRole :one
select * from roles where 
name = $1  limit 1;

-- name: GetRolesByApplication :many
select c.id, c.name, c.description from roles c
inner join applications_roles b on c.id = b.roles_id
inner join applications a on b.applications_id = a.id and a.name = $1;

-- name: CreateRole :one
insert into roles (name, "description")
values ($1, $2)
returning *;


-- name: UpdateRole :one
update roles set name = $1, "description" = $2 where name = $3
returning *;

-- name: DeleteRoles :exec
delete from roles where name = $1 ;