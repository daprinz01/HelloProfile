-- name: GetCountries :many
select * from countries;

-- name: GetCountry :one
select * from countries where name = $1  limit 1;

--- name: CreateCountry :one
insert into countries (name, "flag_image_url")
values ($1, $2)
returning *;

-- name: UpdateCountry :one
update countries set name = $1, flag_image_url = $2 where name = $3
returning *;

-- name: DeleteCountry :exec
delete from countries where name = $1  or flag_image_url = $1;