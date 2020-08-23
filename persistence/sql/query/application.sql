-- name: GetApplications :many
select * from applications;

-- name: GetApplication :one
select * from applications where id == $1 or "name" == $1 limit 1;

-- name: CreateApplication :one
insert into applications (
  "name",
  "description",
  "created_at")
  values($1, $2, $3)
  returning *;

-- name: AddApplicationRole :one
insert into applications_roles (
    applications_id, roles_id) 
values ((select a.id from applications a where a.id == $1 or a."name" == $1), 
(select b.id from roles b where b.id == $2 or b."name" == $2))
returning *;

-- name: UpdateApplicationRole :one
update applications_roles set applications_id = (select a.id from applications a where a.id == $1 or a."name" == $1 limit 1) , 
roles_id = (select b.id from roles b where b.id == $2 or b."name" == $2) 
where applications_id == (select c.id from applications c where c.id == $3 or c."name" == $3 limit 1) 
and roles_id == (select d.id from roles d where d.id == $4 or d.name == $4 limit 1) 
returning *;


-- name: UpdateApplication :one
update applications set "name" = $1, "description" = $2, created_at = $3 
where "name" == $4
returning *;

-- name: DeleteApplication :exec
delete from applications where id == $1 or "name" == $1;

