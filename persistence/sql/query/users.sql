-- name: GetUsers :many
SELECT
  *
FROM
  user_details;

-- name: GetUser :one
SELECT
  *
FROM
  user_details
WHERE
  username = $1
  OR email = $1
  or id = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users ("firstname", "lastname", "username", "email", "phone", "is_email_confirmed", "password", "is_password_system_generated", "created_at", "is_locked_out", "image_url", "is_active")
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING
  *;

-- name: AddUserLanguage :one
INSERT INTO user_languages (user_id, language_id, proficiency)
  VALUES ((
      SELECT
        a.id
      FROM
        users a
      WHERE
        a.username = $1
        OR a.email = $1
      LIMIT 1),
    (
      SELECT
        b.id
      FROM
        languages b
      WHERE
        b.name = $2), $3)
RETURNING
  *;

-- name: GetUserLanguages :many
SELECT
  a.id,
  a.name,
  d.proficiency
FROM
  languages a
  INNER JOIN user_languages d ON a.id = d.language_id
  INNER JOIN users e ON e.id = d.user_id
  INNER JOIN users f ON f.username = $1
    OR f.email = $1;

-- name: DeleteUserLanguage :exec
DELETE FROM user_languages a
WHERE a.user_id = (
    SELECT
      b.id
    FROM
      users b
    WHERE
      b.username = $1
      OR b.email = $1)
  AND a.language_id = (
    SELECT
      c.id
    FROM
      languages c
    WHERE
      c.name = $2);

-- name: UpdateUserLanguage :one
UPDATE
  user_languages
SET
  user_id = (
    SELECT
      a.id
    FROM
      users a
    WHERE
      a.username = $1
      OR a.email = $1
    LIMIT 1),
language_id = (
  SELECT
    b.id
  FROM
    languages b
  WHERE
    b.name = $2)
WHERE
  user_id = (
    SELECT
      c.id
    FROM
      users c
    WHERE
      c.username = $3
      OR c.email = $3
    LIMIT 1)
AND language_id = (
  SELECT
    d.id
  FROM
    languages d
  WHERE
    d.name = $4)
RETURNING
  *;

-- name: GetUserTimezones :many
SELECT
  *
FROM
  timezones a
WHERE
  a.id = (
    SELECT
      b.timezone_id
    FROM
      user_timezones b
    WHERE
      b.user_id = (
        SELECT
          c.id
        FROM
          users c
        WHERE
          c.username = $1
          OR c.email = $1));

-- name: AddUserTimezone :one
INSERT INTO user_timezones (user_id, timezone_id)
  VALUES ((
      SELECT
        a.id
      FROM
        users a
      WHERE
        a.username = $1
        OR a.email = $1),
      (
        SELECT
          b.id
        FROM
          timezones b
        WHERE
          b.name = $2))
  RETURNING
    *;

-- name: UpdateUserTimezone :one
UPDATE
  user_timezones
SET
  user_id = (
    SELECT
      a.id
    FROM
      users a
    WHERE
      a.username = $1
      OR a.email = $1
    LIMIT 1),
timezone_id = (
  SELECT
    b.id
  FROM
    timezones b
  WHERE
    b.name = $2)
WHERE
  user_id = (
    SELECT
      c.id
    FROM
      users c
    WHERE
      c.username = $3
      OR c.email = $3
    LIMIT 1)
RETURNING
  *;

-- name: GetUserProviders :many
SELECT
  *
FROM
  identity_providers a
WHERE
  a.id = (
    SELECT
      b.identity_provider_id
    FROM
      user_providers b
    WHERE
      b.user_id = (
        SELECT
          c.id
        FROM
          users c
        WHERE
          c.username = $1
          OR c.email = $1));

-- name: AddUserProvider :one
INSERT INTO user_providers (user_id, identity_provider_id)
  VALUES ((
      SELECT
        a.id
      FROM
        users a
      WHERE
        a.username = $1
        OR a.email = $1),
      (
        SELECT
          b.id
        FROM
          identity_providers b
        WHERE
          b.name = $2))
  RETURNING
    *;

-- name: UpdateUserProvider :one
UPDATE
  user_providers
SET
  user_id = (
    SELECT
      id
    FROM
      users a
    WHERE
      a.username = $1
      OR a.email = $1
    LIMIT 1),
identity_provider_id = (
  SELECT
    b.id
  FROM
    identity_providers b
  WHERE
    b.name = $2)
WHERE
  user_id = (
    SELECT
      c.id
    FROM
      users c
    WHERE
      c.username = $1
      OR c.email = $1
    LIMIT 1)
AND identity_provider_id = (
  SELECT
    d.id
  FROM
    identity_providers d
  WHERE
    d.name = $2)
RETURNING
  *;

-- name: DeleteProviders :exec
DELETE FROM user_providers a
WHERE a.user_id = (
    SELECT
      b.id
    FROM
      users b
    WHERE
      b.username = $1
      OR b.email = $1)
  AND a.identity_provider_id = (
    SELECT
      c.id
    FROM
      identity_providers c
    WHERE
      c.name = $2);

-- name: GetUserRoles :many
SELECT
  b.name
FROM
  roles b
  INNER JOIN user_roles a ON b.Id = a.role_id
    AND a.user_id = (
      SELECT
        c.id
      FROM
        users c
    WHERE
      c.username = $1
      OR c.email = $1
    LIMIT 1);

-- select b.name from roles b where b.Id = (select a.role_id from user_roles a where a.user_id = $1);
-- name: AddUserRole :one
INSERT INTO user_roles (user_id, role_id)
  VALUES ((
      SELECT
        d.id
      FROM
        users d
      WHERE
        d.username = $1
        OR d.email = $1),
      (
        SELECT
          a.id
        FROM
          roles a
        WHERE
          a.name = $2))
  RETURNING
    *;

-- name: UpdateUserRole :one
UPDATE
  user_roles
SET
  user_id = (
    SELECT
      a.id
    FROM
      users a
    WHERE
      a.username = $1
      OR a.email = $1
    LIMIT 1),
role_id = (
  SELECT
    b.id
  FROM
    roles b
  WHERE
    b.name = $2)
WHERE
  user_id = (
    SELECT
      c.id
    FROM
      users c
    WHERE
      c.username = $3
      OR c.email = $3
    LIMIT 1)
AND role_id = (
  SELECT
    d.id
  FROM
    roles d
  WHERE
    d.name = $4)
RETURNING
  *;

-- name: UpdateUser :one
UPDATE
  users
SET
  "firstname" = $1,
  "lastname" = $2,
  "username" = $3,
  "email" = $4,
  "is_email_confirmed" = $5,
  "password" = $6,
  "is_password_system_generated" = $7,
  "created_at" = $8,
  "is_locked_out" = $9,
  "image_url" = $10,
  "is_active" = $11,
  "phone" = $12
WHERE
  "username" = $13
  OR "email" = $13
  or id=$13
RETURNING
  *;

-- name: DeleteUser :exec
UPDATE
  users
SET
  is_active = FALSE
WHERE
  email = $1
  OR username = $1;

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
