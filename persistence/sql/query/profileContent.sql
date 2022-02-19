-- name: GetProfileContents :many
select * from profile_contents where profile_id = $1;

-- name: GetProfileContent :one
select * from profile_contents where id=$1 limit 1;

-- name: AddProfileContent :one
insert into profile_contents(title, display_title, "description", "url", profile_id, call_to_action_id, "order", content_id)
values($1, $2, $3, $4, $5, $6, $7, $8) returning *;

-- name: UpdateProfileContent :exec
update profile_contents set title=$1, display_title=$2, "description"=$3, "url"=$4, call_to_action_id=$5, "order"=$6
where id=$7;

-- name: DeleteProfileContent :exec
delete from profile_contents CASCADE where id=$1;

-- name: GetAllContentTypes :many
select * from content;