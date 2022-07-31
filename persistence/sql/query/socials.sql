-- name: GetSocials :many
select * from socials;

-- name: GetSocial :one
select * from socials where id=$1 limit 1;

-- name: AddSocial :one
insert into socials(name, placeholder, image_url)
VALUES
($1, $2, $3) returning *;

-- name: UpdateSocial :exec
update socials set name=$1, placeholder=$2, image_url = $3 where id=$4;

-- name: DeleteSocial :exec
delete from socials where id=$1;

-- name: GetProfileSocials :many
select a.username, b.name, b.placeholder, b.image_url, a."order", a.socials_id, a.profile_id, a.id from profile_socials a 
inner join socials b on a.socials_id = b.id and a.profile_id = $1;

-- name: GetProfileSocial :one
select * from profile_socials where id=$1;

-- name: AddProfileSocial :one
insert into profile_socials(username, socials_id, profile_id, "order")VALUES
($1, $2, $3, $4) returning *;

-- name: UpdateProfileSocial :exec
update profile_socials set username=$1, "order"=$2 where id = $3;

-- name: DeleteProfileSocial :exec
delete from profile_socials where id=$1;

-- name: DeleteProfileSocials :exec
delete from profile_socials where profile_id=$1;