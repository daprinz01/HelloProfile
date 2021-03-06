// Code generated by sqlc. DO NOT EDIT.
// source: socials.sql

package helloprofiledb

import (
	"context"

	"github.com/google/uuid"
)

const addProfileSocial = `-- name: AddProfileSocial :one
insert into profile_socials(username, socials_id, profile_id, "order")VALUES
($1, $2, $3, $4) returning id, username, socials_id, profile_id, "order"
`

type AddProfileSocialParams struct {
	Username  string    `json:"username"`
	SocialsID uuid.UUID `json:"socials_id"`
	ProfileID uuid.UUID `json:"profile_id"`
	Order     int32     `json:"order"`
}

func (q *Queries) AddProfileSocial(ctx context.Context, arg AddProfileSocialParams) (ProfileSocial, error) {
	row := q.queryRow(ctx, q.addProfileSocialStmt, addProfileSocial,
		arg.Username,
		arg.SocialsID,
		arg.ProfileID,
		arg.Order,
	)
	var i ProfileSocial
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.SocialsID,
		&i.ProfileID,
		&i.Order,
	)
	return i, err
}

const addSocial = `-- name: AddSocial :one
insert into socials(name, placeholder, image_url)
VALUES
($1, $2, $3) returning id, name, placeholder, image_url
`

type AddSocialParams struct {
	Name        string `json:"name"`
	Placeholder string `json:"placeholder"`
	ImageUrl    string `json:"image_url"`
}

func (q *Queries) AddSocial(ctx context.Context, arg AddSocialParams) (Social, error) {
	row := q.queryRow(ctx, q.addSocialStmt, addSocial, arg.Name, arg.Placeholder, arg.ImageUrl)
	var i Social
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Placeholder,
		&i.ImageUrl,
	)
	return i, err
}

const deleteProfileSocial = `-- name: DeleteProfileSocial :exec
delete from profile_socials where id=$1
`

func (q *Queries) DeleteProfileSocial(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteProfileSocialStmt, deleteProfileSocial, id)
	return err
}

const deleteSocial = `-- name: DeleteSocial :exec
delete from socials where id=$1
`

func (q *Queries) DeleteSocial(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteSocialStmt, deleteSocial, id)
	return err
}

const getProfileSocial = `-- name: GetProfileSocial :one
select id, username, socials_id, profile_id, "order" from profile_socials where id=$1
`

func (q *Queries) GetProfileSocial(ctx context.Context, id uuid.UUID) (ProfileSocial, error) {
	row := q.queryRow(ctx, q.getProfileSocialStmt, getProfileSocial, id)
	var i ProfileSocial
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.SocialsID,
		&i.ProfileID,
		&i.Order,
	)
	return i, err
}

const getProfileSocials = `-- name: GetProfileSocials :many
select a.username, b.name, b.placeholder, b.image_url, a."order", a.socials_id, a.profile_id, a.id from profile_socials a 
inner join socials b on a.socials_id = b.id and a.profile_id = $1
`

type GetProfileSocialsRow struct {
	Username    string    `json:"username"`
	Name        string    `json:"name"`
	Placeholder string    `json:"placeholder"`
	ImageUrl    string    `json:"image_url"`
	Order       int32     `json:"order"`
	SocialsID   uuid.UUID `json:"socials_id"`
	ProfileID   uuid.UUID `json:"profile_id"`
	ID          uuid.UUID `json:"id"`
}

func (q *Queries) GetProfileSocials(ctx context.Context, profileID uuid.UUID) ([]GetProfileSocialsRow, error) {
	rows, err := q.query(ctx, q.getProfileSocialsStmt, getProfileSocials, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProfileSocialsRow
	for rows.Next() {
		var i GetProfileSocialsRow
		if err := rows.Scan(
			&i.Username,
			&i.Name,
			&i.Placeholder,
			&i.ImageUrl,
			&i.Order,
			&i.SocialsID,
			&i.ProfileID,
			&i.ID,
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

const getSocial = `-- name: GetSocial :one
select id, name, placeholder, image_url from socials where id=$1 limit 1
`

func (q *Queries) GetSocial(ctx context.Context, id uuid.UUID) (Social, error) {
	row := q.queryRow(ctx, q.getSocialStmt, getSocial, id)
	var i Social
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Placeholder,
		&i.ImageUrl,
	)
	return i, err
}

const getSocials = `-- name: GetSocials :many
select id, name, placeholder, image_url from socials
`

func (q *Queries) GetSocials(ctx context.Context) ([]Social, error) {
	rows, err := q.query(ctx, q.getSocialsStmt, getSocials)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Social
	for rows.Next() {
		var i Social
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Placeholder,
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

const updateProfileSocial = `-- name: UpdateProfileSocial :exec
update profile_socials set username=$1, "order"=$2 where id = $3
`

type UpdateProfileSocialParams struct {
	Username string    `json:"username"`
	Order    int32     `json:"order"`
	ID       uuid.UUID `json:"id"`
}

func (q *Queries) UpdateProfileSocial(ctx context.Context, arg UpdateProfileSocialParams) error {
	_, err := q.exec(ctx, q.updateProfileSocialStmt, updateProfileSocial, arg.Username, arg.Order, arg.ID)
	return err
}

const updateSocial = `-- name: UpdateSocial :exec
update socials set name=$1, placeholder=$2, image_url = $3 where id=$4
`

type UpdateSocialParams struct {
	Name        string    `json:"name"`
	Placeholder string    `json:"placeholder"`
	ImageUrl    string    `json:"image_url"`
	ID          uuid.UUID `json:"id"`
}

func (q *Queries) UpdateSocial(ctx context.Context, arg UpdateSocialParams) error {
	_, err := q.exec(ctx, q.updateSocialStmt, updateSocial,
		arg.Name,
		arg.Placeholder,
		arg.ImageUrl,
		arg.ID,
	)
	return err
}
