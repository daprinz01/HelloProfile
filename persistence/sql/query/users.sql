-- name: GetUsers :many
select * from user_details;

-- name: GetUser :one
select * from user_details 
where username = $1 or email = $1 limit 1;

-- name: CreateUser :one
insert into users ("firstname",
  "lastname",
  "username",
  "email",
  "is_email_confirmed",
  "password",
  "is_password_system_generated",
  "address",
  "city" ,
  "state" ,
  "country" ,
  "created_at",
  "is_locked_out",
  "image_url",
  "is_active")
  values ($1, $2, $3, $4, $5, $6, $7, $8,$9, $10, $11, $12, $13, $14, $15)
  returning *;

-- name: AddUserLanguage :one
insert into user_languages (
    user_id, language_id
) values ((select a.id from users a where a.username = $1 or a.email = $1 limit 1), (select b.id from languages b where  b.name = $2))
returning *;

-- name: UpdateUserLanguage :one
update user_languages set user_id = (select a.id from users a where a.username = $1 or a.email = $1 limit 1) , 
language_id = (select b.id from languages b where b.name = $2) 
where user_id = (select c.id from users c where c.username = $3 or c.email = $3 limit 1)
and language_id = (select d.id from languages d where  d.name = $4) returning *;

-- name: AddUserTimezone :one
insert into user_timezones (
    user_id, timezone_id
) values ((select a.id from users a where a.username = $1 or a.email = $1), (select b.id from timezones b where b.name = $2))
returning *;

-- name: UpdateUserTimezone :one
update user_timezones set user_id = (select a.id from users a where a.username = $1 or a.email = $1 limit 1) , 
timezone_id = (select b.id from timezones b where b.name = $2) where user_id = (select c.id from users c where c.username = $3 or c.email = $3 limit 1) returning *;

-- name: AddUserProvider :one
insert into user_providers (
    user_id, identity_provider_id
) values ((select a.id from users a where  a.username = $1 or a.email = $1), (select b.id from identity_providers b where  b.name = $2))
returning *;

-- name: UpdateUserProvider :one
update user_providers set user_id = (select id from users a where  a.username = $1 or a.email = $1 limit 1) , 
identity_provider_id = (select b.id from identity_providers b where  b.name = $2) where user_id = (select c.id from users c where c.username = $1 or c.email = $1 limit 1) and identity_provider_id = (select d.id from identity_providers d where  d.name = $2)  returning *;

-- name: GetUserRoles :many
select b.name from roles b where b.Id = (select a.role_id from user_roles a where a.user_id = $1);


-- name: AddUserRole :one
insert into user_roles (
    user_id, role_id
) values ((select d.id from users d where d.username = $1 or d.email = $1), (select a.id from roles a where  a.name = $2))
returning *;

-- name: UpdateUserRole :one
update user_roles set user_id = (select a.id from users a where a.username = $1 or a.email = $1 limit 1) , 
role_id = (select b.id from roles b where  b.name = $2) where user_id = (select c.id from users c where c.username = $3 or c.email = $3 limit 1) 
and role_id = (select d.id from roles d where  d.name = $4) returning *;


-- name: UpdateUser :one
  update users set "firstname" = $1,
  "lastname" = $2,
  "username" = $3,
  "email" = $4,
  "is_email_confirmed" = $5,
  "password" = $6,
  "is_password_system_generated" = $7,
  "address" = $8,
  "city" = $9,
  "state" = $10,
  "country" = $11,
  "created_at" = $12,
  "is_locked_out" = $13,
  "image_url" = $14,
  "is_active" = $15
  where "username" = $16 or "email" = $16
  returning *;

-- name: DeleteUser :exec
 update users set is_active = false;

 

-- create table "applications" (
--   "id" bigserial PRIMARY KEY,
--   name varchar NOT NULL,
--   "description" varchar NOT NULL,
--   "created_at" timestamptz NOT NULL DEFAULT (now()),

--   CONSTRAINT "uc_applications" UNIQUE ("id", name)
-- );

-- create table "users" (
--   "id" bigserial primary key,
--   "firstname" varchar  null,
--   "lastname" varchar  null,
--   "username" varchar  null,
--   "email" varchar not null,
--   "is_email_confirmed" BOOLEAN not null DEFAULT FALSE,
--   "password" varchar  null,
--   "is_password_system_generated" BOOLEAN not null DEFAULT FALSE,
--   "address" varchar  null,
--   "city" VARCHAR null,
--   "state" VARCHAR null,
--   "country" VARCHAR null,
--   "created_at" timestamptz NOT NULL DEFAULT (now()),
--   "is_locked_out" BOOLEAN not null default FALSE,
--   "image_url" varchar  null,
--   "is_active" BOOLEAN not null DEFAULT TRUE,
--   CONSTRAINT "uc_users" UNIQUE ("id", "username", "email")
-- );

-- create table "languages" (
--     "id" bigserial primary key,
--     name varchar not null,
--     CONSTRAINT "uc_languages" UNIQUE ("id", name)
-- );

-- create table "user_languages" (
--     "id" bigserial primary key,
--     "user_id" bigserial ,
--     "language_id" bigserial,
--     CONSTRAINT "uc_user_languages" UNIQUE ("id")
-- );

-- create table "timezones" (
--     "id" bigserial primary key,
--     name varchar not null,
--     "zone" varchar not null,
--     CONSTRAINT "uc_timezones" UNIQUE ("id", name)
-- );

-- create table "user_timezones" (
--     "id" bigserial primary key,
--     "user_id" bigserial,
--     "timezone_id" bigserial,
--     CONSTRAINT "uc_user_timezones" UNIQUE ("id", "user_id")
-- );

-- create table "roles" (
--     "id" bigserial primary key,
--     name varchar not null,
--     "description" varchar not null,
--     CONSTRAINT "uc_roles" UNIQUE ("id", name)
-- );

-- create table "user_roles" (
--     "id" bigserial primary key,
--     "user_id" bigserial,
--     "role_id" bigserial,
--     CONSTRAINT "uc_user_roles" UNIQUE ("id")
-- );

-- create table "applications_roles" (
--     "id" bigserial primary key,
--     "applications_id" bigserial,
--     "roles_id" bigserial,
--     CONSTRAINT "uc_applications_roles" UNIQUE ("id")
-- );

-- create table "identity_providers" (
--     "id" bigserial primary key,
--     name varchar not null,
--     "client_id" varchar not null,
--     "client_secret" varchar not null,
--     "image_url" varchar not null,
--     CONSTRAINT "uc_identity_providers" UNIQUE ("id", name)
-- );

-- create table "user_providers" (
--     "id" bigserial primary key,
--     "user_id" bigserial,
--     "identity_provider_id" bigserial,
--     CONSTRAINT "uc_user_providers" UNIQUE ("id")
-- );
-- create TABLE "countries" (
--     "id" bigserial PRIMARY KEY,
--     name VARCHAR not NULL,
--     "flag_image_url" VARCHAR null,

--     CONSTRAINT "uc_countries" UNIQUE ("id", name, "flag_image_url")
-- );
-- create TABLE "states" (
--     "id" bigserial PRIMARY KEY,
--     name VARCHAR not NULL,
--     "country_id" bigserial,
--     CONSTRAINT "uc_states" UNIQUE ("id", name)
-- );


-- CREATE  VIEW user_details as
-- SELECT b.firstname, b.lastname, b.email, b.username, b."password", b.address, b.city, b.state, b.country, b.image_url as profile_picture, b.is_email_confirmed, b.is_locked_out, b.is_password_system_generated, b.created_at, b.is_active, d.name as language_name, f.name as role_name, k.name as timezone_name, k.zone, m.name as provider_name, m.client_id, m.client_secret, m.image_url as provider_logo
-- from users b
-- full join languages d on d.id = (select language_id from user_languages e where e.user_id = b.id)
-- full join roles f on f.id = (select role_id from user_roles j where j.user_id = b.id)
-- full join timezones k on k.id = (select timezone_id from user_timezones l where l.user_id = b.id)
-- full join identity_providers m on m.id = (select identity_provider_id from user_providers n where n.user_id = b.id);