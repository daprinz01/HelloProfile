-- name: GetTimezones :many
select * from timezones;

-- name: GetTimezone :one
select * from timezones where 
id == $1 or "name" == $1 limit 1;

-- name: CreateTimezone :one
insert into timezones ("name", "zone")
values ($1, $2)
returning *;

-- name: UpdateTimezone :one
update timezones set "name" = $1, "zone" = $2 where "name" == $3
returning *;

-- name: DeleteTimezone :exec
delete from timezones where id == $1 or "name" == $1;