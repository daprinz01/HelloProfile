-- name: GetSavedProfiles :many
select * from saved_profiles;

-- name: GetSavedProfilesByEmail :many
select * from saved_profiles where email=$1 and is_added = $2;


-- name: GetSavedProfilesByProfileId :many
select * from saved_profiles where profile_id = $1 and is_added = $2;

-- name: GetSavedProfile :one
select * from saved_profiles where 
id = $1  limit 1;

-- name: CreateSavedProfile :one
insert into saved_profiles (first_name, last_name, email, is_added, profile_id)
values ($1, $2, $3, $4, $5)
returning *;


-- name: UpdateSavedProfile :one
update saved_profiles set first_name = $1, last_name = $2, is_added = $3 where id = $3
returning *;

-- name: DeleteSavedProfile :exec
delete from saved_profiles where id = $1 ;