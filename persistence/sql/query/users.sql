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
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users ("firstname", "lastname", "username", "email", "phone", "is_email_confirmed", "password", "is_password_system_generated", "created_at", "is_locked_out", "image_url", "is_active")
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING
  *;

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
  "is_locked_out" = $8,
  "image_url" = $9,
  "is_active" = $10,
  "phone" = $11,
  "country" = $12,
  "city" = $13
WHERE
  "username" = $14
  OR "email" = $14
RETURNING
  *;

-- name: DeleteUser :exec
Delete from users CASCADE
WHERE
  email = $1
  OR username = $1;
