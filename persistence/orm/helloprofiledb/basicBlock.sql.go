// Code generated by sqlc. DO NOT EDIT.
// source: basicBlock.sql

package helloprofiledb

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const addBasicBlock = `-- name: AddBasicBlock :one
insert into basic_block(
    profile_photo_url,
    cover_photo_url,
    cover_colour,
    fullname,
    title,
    bio
)VALUES(
    $1, $2, $3, $4, $5, $6
) RETURNING id, profile_photo_url, cover_photo_url, cover_colour, fullname, title, bio
`

type AddBasicBlockParams struct {
	ProfilePhotoUrl sql.NullString `json:"profile_photo_url"`
	CoverPhotoUrl   sql.NullString `json:"cover_photo_url"`
	CoverColour     sql.NullString `json:"cover_colour"`
	Fullname        string         `json:"fullname"`
	Title           string         `json:"title"`
	Bio             string         `json:"bio"`
}

func (q *Queries) AddBasicBlock(ctx context.Context, arg AddBasicBlockParams) (BasicBlock, error) {
	row := q.queryRow(ctx, q.addBasicBlockStmt, addBasicBlock,
		arg.ProfilePhotoUrl,
		arg.CoverPhotoUrl,
		arg.CoverColour,
		arg.Fullname,
		arg.Title,
		arg.Bio,
	)
	var i BasicBlock
	err := row.Scan(
		&i.ID,
		&i.ProfilePhotoUrl,
		&i.CoverPhotoUrl,
		&i.CoverColour,
		&i.Fullname,
		&i.Title,
		&i.Bio,
	)
	return i, err
}

const deleteBasicBlock = `-- name: DeleteBasicBlock :exec
delete from basic_block where id=$1
`

func (q *Queries) DeleteBasicBlock(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteBasicBlockStmt, deleteBasicBlock, id)
	return err
}

const getBasicBlock = `-- name: GetBasicBlock :one
select id, profile_photo_url, cover_photo_url, cover_colour, fullname, title, bio from basic_block where id=$1 limit 1
`

func (q *Queries) GetBasicBlock(ctx context.Context, id uuid.UUID) (BasicBlock, error) {
	row := q.queryRow(ctx, q.getBasicBlockStmt, getBasicBlock, id)
	var i BasicBlock
	err := row.Scan(
		&i.ID,
		&i.ProfilePhotoUrl,
		&i.CoverPhotoUrl,
		&i.CoverColour,
		&i.Fullname,
		&i.Title,
		&i.Bio,
	)
	return i, err
}

const updateBasicBlock = `-- name: UpdateBasicBlock :exec
update basic_block set profile_photo_url=$1,
    cover_photo_url=$2,
    cover_colour=$3,
    fullname=$4,
    title=$5,
    bio=$6 where id=$7
`

type UpdateBasicBlockParams struct {
	ProfilePhotoUrl sql.NullString `json:"profile_photo_url"`
	CoverPhotoUrl   sql.NullString `json:"cover_photo_url"`
	CoverColour     sql.NullString `json:"cover_colour"`
	Fullname        string         `json:"fullname"`
	Title           string         `json:"title"`
	Bio             string         `json:"bio"`
	ID              uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateBasicBlock(ctx context.Context, arg UpdateBasicBlockParams) error {
	_, err := q.exec(ctx, q.updateBasicBlockStmt, updateBasicBlock,
		arg.ProfilePhotoUrl,
		arg.CoverPhotoUrl,
		arg.CoverColour,
		arg.Fullname,
		arg.Title,
		arg.Bio,
		arg.ID,
	)
	return err
}

const updateProfileWithBasicBlockId = `-- name: UpdateProfileWithBasicBlockId :exec
update profiles set basic_block_id=$1 where id=$2
`

type UpdateProfileWithBasicBlockIdParams struct {
	BasicBlockID uuid.NullUUID `json:"basic_block_id"`
	ID           uuid.UUID     `json:"id"`
}

func (q *Queries) UpdateProfileWithBasicBlockId(ctx context.Context, arg UpdateProfileWithBasicBlockIdParams) error {
	_, err := q.exec(ctx, q.updateProfileWithBasicBlockIdStmt, updateProfileWithBasicBlockId, arg.BasicBlockID, arg.ID)
	return err
}
