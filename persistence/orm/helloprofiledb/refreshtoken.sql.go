// Code generated by sqlc. DO NOT EDIT.
// source: refreshtoken.sql

package helloprofiledb

import (
	"context"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
insert into refresh_token (
    "user_id",
    "token"
    )
    values
    ($1, $2)
    returning id, user_id, token, created_at
`

type CreateRefreshTokenParams struct {
	UserID uuid.UUID `json:"user_id"`
	Token  string    `json:"token"`
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.queryRow(ctx, q.createRefreshTokenStmt, createRefreshToken, arg.UserID, arg.Token)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.CreatedAt,
	)
	return i, err
}

const deleteRefreshToken = `-- name: DeleteRefreshToken :exec
delete from refresh_token where "token" = $1
`

func (q *Queries) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := q.exec(ctx, q.deleteRefreshTokenStmt, deleteRefreshToken, token)
	return err
}

const getRefreshToken = `-- name: GetRefreshToken :one
select id, user_id, token, created_at from refresh_token where token = $1
`

func (q *Queries) GetRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.queryRow(ctx, q.getRefreshTokenStmt, getRefreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.CreatedAt,
	)
	return i, err
}

const getRefreshTokens = `-- name: GetRefreshTokens :many
select id, user_id, token, created_at from refresh_token
`

func (q *Queries) GetRefreshTokens(ctx context.Context) ([]RefreshToken, error) {
	rows, err := q.query(ctx, q.getRefreshTokensStmt, getRefreshTokens)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RefreshToken
	for rows.Next() {
		var i RefreshToken
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Token,
			&i.CreatedAt,
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

const updateRefreshToken = `-- name: UpdateRefreshToken :one
update refresh_token set "user_id" = $1, "token" = $2 where "token" = $3 returning id, user_id, token, created_at
`

type UpdateRefreshTokenParams struct {
	UserID  uuid.UUID `json:"user_id"`
	Token   string    `json:"token"`
	Token_2 string    `json:"token_2"`
}

func (q *Queries) UpdateRefreshToken(ctx context.Context, arg UpdateRefreshTokenParams) (RefreshToken, error) {
	row := q.queryRow(ctx, q.updateRefreshTokenStmt, updateRefreshToken, arg.UserID, arg.Token, arg.Token_2)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.CreatedAt,
	)
	return i, err
}
