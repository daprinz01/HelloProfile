-- name: GetLanguages :many
select * from languages;

-- name: GetLanguage :one
select * from languages where name = $1  limit 1;

-- name: CreateLanguage :one
insert into languages (name) values ($1)
returning *;

-- name: UpdateLanguage :one
update languages set name = $1 where name = $2
returning *;

-- name: DeleteLanguage :exec
delete from languages where name = $1 ;