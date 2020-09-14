-- name: GetLanguageProficiencies :many
select * from language_proficiency;

-- name: GetLanguageProficiency :one
select * from language_proficiency where proficiency = $1  limit 1;

-- name: CreateLanguageProficiency :one
insert into language_proficiency (proficiency) values ($1)
returning *;

-- name: UpdateLanguageProficiency :one
update language_proficiency set proficiency = $1 where proficiency = $2
returning *;

-- name: DeleteLanguageProficiency :exec
delete from language_proficiency where proficiency = $1 ;