// Code generated by sqlc. DO NOT EDIT.
// source: user_login.sql

package helloprofiledb

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUserLogin = `-- name: CreateUserLogin :one
insert into user_login (
    user_id,
    login_status,
    response_code,
    response_description,
    device,
    ip_address,
    longitude,
    latitude,
    resolved
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning id, user_id, login_time, login_status, response_code, response_description, device, ip_address, longitude, latitude, resolved
`

type CreateUserLoginParams struct {
	UserID              uuid.UUID      `json:"user_id"`
	LoginStatus         bool           `json:"login_status"`
	ResponseCode        sql.NullString `json:"response_code"`
	ResponseDescription sql.NullString `json:"response_description"`
	Device              sql.NullString `json:"device"`
	IpAddress           sql.NullString `json:"ip_address"`
	Longitude           sql.NullString `json:"longitude"`
	Latitude            sql.NullString `json:"latitude"`
	Resolved            bool           `json:"resolved"`
}

func (q *Queries) CreateUserLogin(ctx context.Context, arg CreateUserLoginParams) (UserLogin, error) {
	row := q.queryRow(ctx, q.createUserLoginStmt, createUserLogin,
		arg.UserID,
		arg.LoginStatus,
		arg.ResponseCode,
		arg.ResponseDescription,
		arg.Device,
		arg.IpAddress,
		arg.Longitude,
		arg.Latitude,
		arg.Resolved,
	)
	var i UserLogin
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.LoginTime,
		&i.LoginStatus,
		&i.ResponseCode,
		&i.ResponseDescription,
		&i.Device,
		&i.IpAddress,
		&i.Longitude,
		&i.Latitude,
		&i.Resolved,
	)
	return i, err
}

const deleteUserLogin = `-- name: DeleteUserLogin :exec
delete from user_login where user_id = $1
`

func (q *Queries) DeleteUserLogin(ctx context.Context, userID uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteUserLoginStmt, deleteUserLogin, userID)
	return err
}

const getUnResoledLogins = `-- name: GetUnResoledLogins :many
select id, user_id, login_time, login_status, response_code, response_description, device, ip_address, longitude, latitude, resolved from user_login where user_id = $1 and resolved = false
`

func (q *Queries) GetUnResoledLogins(ctx context.Context, userID uuid.UUID) ([]UserLogin, error) {
	rows, err := q.query(ctx, q.getUnResoledLoginsStmt, getUnResoledLogins, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserLogin
	for rows.Next() {
		var i UserLogin
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.LoginTime,
			&i.LoginStatus,
			&i.ResponseCode,
			&i.ResponseDescription,
			&i.Device,
			&i.IpAddress,
			&i.Longitude,
			&i.Latitude,
			&i.Resolved,
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

const getUserLogin = `-- name: GetUserLogin :many
select id, user_id, login_time, login_status, response_code, response_description, device, ip_address, longitude, latitude, resolved from user_login where user_id = $1
`

func (q *Queries) GetUserLogin(ctx context.Context, userID uuid.UUID) ([]UserLogin, error) {
	rows, err := q.query(ctx, q.getUserLoginStmt, getUserLogin, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserLogin
	for rows.Next() {
		var i UserLogin
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.LoginTime,
			&i.LoginStatus,
			&i.ResponseCode,
			&i.ResponseDescription,
			&i.Device,
			&i.IpAddress,
			&i.Longitude,
			&i.Latitude,
			&i.Resolved,
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

const getUserLogins = `-- name: GetUserLogins :many
select id, user_id, login_time, login_status, response_code, response_description, device, ip_address, longitude, latitude, resolved from user_login
`

func (q *Queries) GetUserLogins(ctx context.Context) ([]UserLogin, error) {
	rows, err := q.query(ctx, q.getUserLoginsStmt, getUserLogins)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserLogin
	for rows.Next() {
		var i UserLogin
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.LoginTime,
			&i.LoginStatus,
			&i.ResponseCode,
			&i.ResponseDescription,
			&i.Device,
			&i.IpAddress,
			&i.Longitude,
			&i.Latitude,
			&i.Resolved,
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

const updateResolvedLogin = `-- name: UpdateResolvedLogin :exec
update user_login set resolved = true where user_id = $1 and resolved = false
`

func (q *Queries) UpdateResolvedLogin(ctx context.Context, userID uuid.UUID) error {
	_, err := q.exec(ctx, q.updateResolvedLoginStmt, updateResolvedLogin, userID)
	return err
}
