// Code generated by sqlc. DO NOT EDIT.
// source: users.sql

package authdb

import (
	"context"
	"database/sql"
	"time"
)

const addUserLanguage = `-- name: AddUserLanguage :one
insert into user_languages (
    user_id, language_id, proficiency
) values ((select a.id from users a where a.username = $1 or a.email = $1 limit 1), (select b.id from languages b where  b.name = $2), $3)
returning id, user_id, language_id, proficiency
`

type AddUserLanguageParams struct {
	Username    sql.NullString `json:"username"`
	Name        string         `json:"name"`
	Proficiency sql.NullString `json:"proficiency"`
}

func (q *Queries) AddUserLanguage(ctx context.Context, arg AddUserLanguageParams) (UserLanguage, error) {
	row := q.queryRow(ctx, q.addUserLanguageStmt, addUserLanguage, arg.Username, arg.Name, arg.Proficiency)
	var i UserLanguage
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.LanguageID,
		&i.Proficiency,
	)
	return i, err
}

const addUserProvider = `-- name: AddUserProvider :one
insert into user_providers (
    user_id, identity_provider_id
) values ((select a.id from users a where  a.username = $1 or a.email = $1), (select b.id from identity_providers b where  b.name = $2))
returning id, user_id, identity_provider_id
`

type AddUserProviderParams struct {
	Username sql.NullString `json:"username"`
	Name     string         `json:"name"`
}

func (q *Queries) AddUserProvider(ctx context.Context, arg AddUserProviderParams) (UserProvider, error) {
	row := q.queryRow(ctx, q.addUserProviderStmt, addUserProvider, arg.Username, arg.Name)
	var i UserProvider
	err := row.Scan(&i.ID, &i.UserID, &i.IdentityProviderID)
	return i, err
}

const addUserRole = `-- name: AddUserRole :one
insert into user_roles (
    user_id, role_id
) values ((select d.id from users d where d.username = $1 or d.email = $1), (select a.id from roles a where  a.name = $2))
returning id, user_id, role_id
`

type AddUserRoleParams struct {
	Username sql.NullString `json:"username"`
	Name     string         `json:"name"`
}

func (q *Queries) AddUserRole(ctx context.Context, arg AddUserRoleParams) (UserRole, error) {
	row := q.queryRow(ctx, q.addUserRoleStmt, addUserRole, arg.Username, arg.Name)
	var i UserRole
	err := row.Scan(&i.ID, &i.UserID, &i.RoleID)
	return i, err
}

const addUserTimezone = `-- name: AddUserTimezone :one
insert into user_timezones (
    user_id, timezone_id
) values ((select a.id from users a where a.username = $1 or a.email = $1), (select b.id from timezones b where b.name = $2))
returning id, user_id, timezone_id
`

type AddUserTimezoneParams struct {
	Username sql.NullString `json:"username"`
	Name     string         `json:"name"`
}

func (q *Queries) AddUserTimezone(ctx context.Context, arg AddUserTimezoneParams) (UserTimezone, error) {
	row := q.queryRow(ctx, q.addUserTimezoneStmt, addUserTimezone, arg.Username, arg.Name)
	var i UserTimezone
	err := row.Scan(&i.ID, &i.UserID, &i.TimezoneID)
	return i, err
}

const createUser = `-- name: CreateUser :one
insert into users ("firstname",
  "lastname",
  "username",
  "email",
  "phone",
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
  values ($1, $2, $3, $4, $5, $6, $7, $8,$9, $10, $11, $12, $13, $14,$15, $16)
  returning id, firstname, lastname, username, email, phone, is_email_confirmed, password, is_password_system_generated, address, city, state, country, created_at, is_locked_out, image_url, is_active
`

type CreateUserParams struct {
	Firstname                 sql.NullString `json:"firstname"`
	Lastname                  sql.NullString `json:"lastname"`
	Username                  sql.NullString `json:"username"`
	Email                     string         `json:"email"`
	Phone                     sql.NullString `json:"phone"`
	IsEmailConfirmed          bool           `json:"is_email_confirmed"`
	Password                  sql.NullString `json:"password"`
	IsPasswordSystemGenerated bool           `json:"is_password_system_generated"`
	Address                   sql.NullString `json:"address"`
	City                      sql.NullString `json:"city"`
	State                     sql.NullString `json:"state"`
	Country                   sql.NullString `json:"country"`
	CreatedAt                 time.Time      `json:"created_at"`
	IsLockedOut               bool           `json:"is_locked_out"`
	ImageUrl                  sql.NullString `json:"image_url"`
	IsActive                  bool           `json:"is_active"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.queryRow(ctx, q.createUserStmt, createUser,
		arg.Firstname,
		arg.Lastname,
		arg.Username,
		arg.Email,
		arg.Phone,
		arg.IsEmailConfirmed,
		arg.Password,
		arg.IsPasswordSystemGenerated,
		arg.Address,
		arg.City,
		arg.State,
		arg.Country,
		arg.CreatedAt,
		arg.IsLockedOut,
		arg.ImageUrl,
		arg.IsActive,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Username,
		&i.Email,
		&i.Phone,
		&i.IsEmailConfirmed,
		&i.Password,
		&i.IsPasswordSystemGenerated,
		&i.Address,
		&i.City,
		&i.State,
		&i.Country,
		&i.CreatedAt,
		&i.IsLockedOut,
		&i.ImageUrl,
		&i.IsActive,
	)
	return i, err
}

const deleteProviders = `-- name: DeleteProviders :exec
delete from user_providers a where a.user_id = (select b.id from users b where b.username = $1 or b.email = $1) and a.identity_provider_id = (select c.id from identity_providers c where c.name = $2)
`

type DeleteProvidersParams struct {
	Username sql.NullString `json:"username"`
	Name     string         `json:"name"`
}

func (q *Queries) DeleteProviders(ctx context.Context, arg DeleteProvidersParams) error {
	_, err := q.exec(ctx, q.deleteProvidersStmt, deleteProviders, arg.Username, arg.Name)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
 update users set is_active = false where email = $1 or username = $1
`

func (q *Queries) DeleteUser(ctx context.Context, email string) error {
	_, err := q.exec(ctx, q.deleteUserStmt, deleteUser, email)
	return err
}

const deleteUserLanguage = `-- name: DeleteUserLanguage :exec
delete from user_languages a where a.user_id = (select b.id from users b where b.username = $1 or b.email = $1) and a.language_id = (select c.id from languages c where c.name = $2)
`

type DeleteUserLanguageParams struct {
	Username sql.NullString `json:"username"`
	Name     string         `json:"name"`
}

func (q *Queries) DeleteUserLanguage(ctx context.Context, arg DeleteUserLanguageParams) error {
	_, err := q.exec(ctx, q.deleteUserLanguageStmt, deleteUserLanguage, arg.Username, arg.Name)
	return err
}

const getUser = `-- name: GetUser :one
select id, firstname, lastname, username, email, phone, is_email_confirmed, password, is_password_system_generated, address, city, state, country, created_at, is_locked_out, profile_picture, is_active, timezone_name, zone from user_details 
where username = $1 or email = $1 limit 1
`

func (q *Queries) GetUser(ctx context.Context, username sql.NullString) (UserDetail, error) {
	row := q.queryRow(ctx, q.getUserStmt, getUser, username)
	var i UserDetail
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Username,
		&i.Email,
		&i.Phone,
		&i.IsEmailConfirmed,
		&i.Password,
		&i.IsPasswordSystemGenerated,
		&i.Address,
		&i.City,
		&i.State,
		&i.Country,
		&i.CreatedAt,
		&i.IsLockedOut,
		&i.ProfilePicture,
		&i.IsActive,
		&i.TimezoneName,
		&i.Zone,
	)
	return i, err
}

const getUserLanguages = `-- name: GetUserLanguages :many
select a.id, a.name, d.proficiency from languages a inner join user_languages d on a.id = d.language_id inner join users e on e.id = d.user_id inner join users f on f.username = $1 or f.email = $1
`

type GetUserLanguagesRow struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Proficiency sql.NullString `json:"proficiency"`
}

func (q *Queries) GetUserLanguages(ctx context.Context, username sql.NullString) ([]GetUserLanguagesRow, error) {
	rows, err := q.query(ctx, q.getUserLanguagesStmt, getUserLanguages, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserLanguagesRow
	for rows.Next() {
		var i GetUserLanguagesRow
		if err := rows.Scan(&i.ID, &i.Name, &i.Proficiency); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserProviders = `-- name: GetUserProviders :many
select id, name, client_id, client_secret, image_url from identity_providers a where a.id = (select b.identity_provider_id from user_providers b where b.user_id = (select c.id from users c where c.username = $1 or c.email = $1))
`

func (q *Queries) GetUserProviders(ctx context.Context, username sql.NullString) ([]IdentityProvider, error) {
	rows, err := q.query(ctx, q.getUserProvidersStmt, getUserProviders, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []IdentityProvider
	for rows.Next() {
		var i IdentityProvider
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ClientID,
			&i.ClientSecret,
			&i.ImageUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserRoles = `-- name: GetUserRoles :many
select b.name from roles b where b.Id = (select a.role_id from user_roles a where a.user_id = $1)
`

func (q *Queries) GetUserRoles(ctx context.Context, userID sql.NullInt64) ([]string, error) {
	rows, err := q.query(ctx, q.getUserRolesStmt, getUserRoles, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserTimezones = `-- name: GetUserTimezones :many
select id, name, zone from timezones a where a.id = (select b.timezone_id from user_timezones b where b.user_id = (select c.id from users c where c.username = $1 or c.email = $1))
`

func (q *Queries) GetUserTimezones(ctx context.Context, username sql.NullString) ([]Timezone, error) {
	rows, err := q.query(ctx, q.getUserTimezonesStmt, getUserTimezones, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Timezone
	for rows.Next() {
		var i Timezone
		if err := rows.Scan(&i.ID, &i.Name, &i.Zone); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsers = `-- name: GetUsers :many
select id, firstname, lastname, username, email, phone, is_email_confirmed, password, is_password_system_generated, address, city, state, country, created_at, is_locked_out, profile_picture, is_active, timezone_name, zone from user_details
`

func (q *Queries) GetUsers(ctx context.Context) ([]UserDetail, error) {
	rows, err := q.query(ctx, q.getUsersStmt, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserDetail
	for rows.Next() {
		var i UserDetail
		if err := rows.Scan(
			&i.ID,
			&i.Firstname,
			&i.Lastname,
			&i.Username,
			&i.Email,
			&i.Phone,
			&i.IsEmailConfirmed,
			&i.Password,
			&i.IsPasswordSystemGenerated,
			&i.Address,
			&i.City,
			&i.State,
			&i.Country,
			&i.CreatedAt,
			&i.IsLockedOut,
			&i.ProfilePicture,
			&i.IsActive,
			&i.TimezoneName,
			&i.Zone,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
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
  "is_active" = $15,
    "phone" = $17
  where "username" = $16 or "email" = $16
  returning id, firstname, lastname, username, email, phone, is_email_confirmed, password, is_password_system_generated, address, city, state, country, created_at, is_locked_out, image_url, is_active
`

type UpdateUserParams struct {
	Firstname                 sql.NullString `json:"firstname"`
	Lastname                  sql.NullString `json:"lastname"`
	Username                  sql.NullString `json:"username"`
	Email                     string         `json:"email"`
	IsEmailConfirmed          bool           `json:"is_email_confirmed"`
	Password                  sql.NullString `json:"password"`
	IsPasswordSystemGenerated bool           `json:"is_password_system_generated"`
	Address                   sql.NullString `json:"address"`
	City                      sql.NullString `json:"city"`
	State                     sql.NullString `json:"state"`
	Country                   sql.NullString `json:"country"`
	CreatedAt                 time.Time      `json:"created_at"`
	IsLockedOut               bool           `json:"is_locked_out"`
	ImageUrl                  sql.NullString `json:"image_url"`
	IsActive                  bool           `json:"is_active"`
	Username_2                sql.NullString `json:"username_2"`
	Phone                     sql.NullString `json:"phone"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.queryRow(ctx, q.updateUserStmt, updateUser,
		arg.Firstname,
		arg.Lastname,
		arg.Username,
		arg.Email,
		arg.IsEmailConfirmed,
		arg.Password,
		arg.IsPasswordSystemGenerated,
		arg.Address,
		arg.City,
		arg.State,
		arg.Country,
		arg.CreatedAt,
		arg.IsLockedOut,
		arg.ImageUrl,
		arg.IsActive,
		arg.Username_2,
		arg.Phone,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Username,
		&i.Email,
		&i.Phone,
		&i.IsEmailConfirmed,
		&i.Password,
		&i.IsPasswordSystemGenerated,
		&i.Address,
		&i.City,
		&i.State,
		&i.Country,
		&i.CreatedAt,
		&i.IsLockedOut,
		&i.ImageUrl,
		&i.IsActive,
	)
	return i, err
}

const updateUserLanguage = `-- name: UpdateUserLanguage :one
update user_languages set user_id = (select a.id from users a where a.username = $1 or a.email = $1 limit 1) , 
language_id = (select b.id from languages b where b.name = $2) 
where user_id = (select c.id from users c where c.username = $3 or c.email = $3 limit 1)
and language_id = (select d.id from languages d where  d.name = $4) returning id, user_id, language_id, proficiency
`

type UpdateUserLanguageParams struct {
	Username   sql.NullString `json:"username"`
	Name       string         `json:"name"`
	Username_2 sql.NullString `json:"username_2"`
	Name_2     string         `json:"name_2"`
}

func (q *Queries) UpdateUserLanguage(ctx context.Context, arg UpdateUserLanguageParams) (UserLanguage, error) {
	row := q.queryRow(ctx, q.updateUserLanguageStmt, updateUserLanguage,
		arg.Username,
		arg.Name,
		arg.Username_2,
		arg.Name_2,
	)
	var i UserLanguage
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.LanguageID,
		&i.Proficiency,
	)
	return i, err
}

const updateUserProvider = `-- name: UpdateUserProvider :one
update user_providers set user_id = (select id from users a where  a.username = $1 or a.email = $1 limit 1) , 
identity_provider_id = (select b.id from identity_providers b where  b.name = $2) where user_id = (select c.id from users c where c.username = $1 or c.email = $1 limit 1) and identity_provider_id = (select d.id from identity_providers d where  d.name = $2)  returning id, user_id, identity_provider_id
`

type UpdateUserProviderParams struct {
	Username sql.NullString `json:"username"`
	Name     string         `json:"name"`
}

func (q *Queries) UpdateUserProvider(ctx context.Context, arg UpdateUserProviderParams) (UserProvider, error) {
	row := q.queryRow(ctx, q.updateUserProviderStmt, updateUserProvider, arg.Username, arg.Name)
	var i UserProvider
	err := row.Scan(&i.ID, &i.UserID, &i.IdentityProviderID)
	return i, err
}

const updateUserRole = `-- name: UpdateUserRole :one
update user_roles set user_id = (select a.id from users a where a.username = $1 or a.email = $1 limit 1) , 
role_id = (select b.id from roles b where  b.name = $2) where user_id = (select c.id from users c where c.username = $3 or c.email = $3 limit 1) 
and role_id = (select d.id from roles d where  d.name = $4) returning id, user_id, role_id
`

type UpdateUserRoleParams struct {
	Username   sql.NullString `json:"username"`
	Name       string         `json:"name"`
	Username_2 sql.NullString `json:"username_2"`
	Name_2     string         `json:"name_2"`
}

func (q *Queries) UpdateUserRole(ctx context.Context, arg UpdateUserRoleParams) (UserRole, error) {
	row := q.queryRow(ctx, q.updateUserRoleStmt, updateUserRole,
		arg.Username,
		arg.Name,
		arg.Username_2,
		arg.Name_2,
	)
	var i UserRole
	err := row.Scan(&i.ID, &i.UserID, &i.RoleID)
	return i, err
}

const updateUserTimezone = `-- name: UpdateUserTimezone :one
update user_timezones set user_id = (select a.id from users a where a.username = $1 or a.email = $1 limit 1) , 
timezone_id = (select b.id from timezones b where b.name = $2) where user_id = (select c.id from users c where c.username = $3 or c.email = $3 limit 1) returning id, user_id, timezone_id
`

type UpdateUserTimezoneParams struct {
	Username   sql.NullString `json:"username"`
	Name       string         `json:"name"`
	Username_2 sql.NullString `json:"username_2"`
}

func (q *Queries) UpdateUserTimezone(ctx context.Context, arg UpdateUserTimezoneParams) (UserTimezone, error) {
	row := q.queryRow(ctx, q.updateUserTimezoneStmt, updateUserTimezone, arg.Username, arg.Name, arg.Username_2)
	var i UserTimezone
	err := row.Scan(&i.ID, &i.UserID, &i.TimezoneID)
	return i, err
}
