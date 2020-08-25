-- name: GetIdentityProviders :many
select * from identity_providers;

-- name: GetIdentityProvider :one
select * from identity_providers where name = $1  limit 1;

-- name: CreateIdentityProvider :one
insert into identity_providers (
  name,
  "client_id",
  "client_secret",
  "image_url")
  values($1, $2, $3, $4)
  returning *;

-- name: UpdateIdentityProvider :one
update identity_providers set name = $1, "client_id" = $2, client_secret = $3, image_url = $4 where name = $5
returning *;

-- name: DeleteIdentityProvider :exec
delete from identity_providers where name = $1 ;

