// Code generated by sqlc. DO NOT EDIT.
// source: profileContent.sql

package helloprofiledb

import (
	"context"

	"github.com/google/uuid"
)

const addProfileContent = `-- name: AddProfileContent :one
insert into profile_contents(title, display_title, "description", "url", profile_id, call_to_action_id, "order", content_id)
values($1, $2, $3, $4, $5, $6, $7, $8) returning id, title, display_title, description, url, profile_id, call_to_action_id, content_id, "order"
`

type AddProfileContentParams struct {
	Title          string    `json:"title"`
	DisplayTitle   string    `json:"display_title"`
	Description    string    `json:"description"`
	Url            string    `json:"url"`
	ProfileID      uuid.UUID `json:"profile_id"`
	CallToActionID uuid.UUID `json:"call_to_action_id"`
	Order          int32     `json:"order"`
	ContentID      uuid.UUID `json:"content_id"`
}

func (q *Queries) AddProfileContent(ctx context.Context, arg AddProfileContentParams) (ProfileContent, error) {
	row := q.queryRow(ctx, q.addProfileContentStmt, addProfileContent,
		arg.Title,
		arg.DisplayTitle,
		arg.Description,
		arg.Url,
		arg.ProfileID,
		arg.CallToActionID,
		arg.Order,
		arg.ContentID,
	)
	var i ProfileContent
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.DisplayTitle,
		&i.Description,
		&i.Url,
		&i.ProfileID,
		&i.CallToActionID,
		&i.ContentID,
		&i.Order,
	)
	return i, err
}

const deleteProfileContent = `-- name: DeleteProfileContent :exec
delete from profile_contents where id=$1
`

func (q *Queries) DeleteProfileContent(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteProfileContentStmt, deleteProfileContent, id)
	return err
}

const getAllContentTypes = `-- name: GetAllContentTypes :many
select id, type, image_url from content
`

func (q *Queries) GetAllContentTypes(ctx context.Context) ([]Content, error) {
	rows, err := q.query(ctx, q.getAllContentTypesStmt, getAllContentTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Content
	for rows.Next() {
		var i Content
		if err := rows.Scan(&i.ID, &i.Type, &i.ImageUrl); err != nil {
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

const getProfileContent = `-- name: GetProfileContent :one
select id, title, display_title, description, url, profile_id, call_to_action_id, content_id, "order" from profile_contents where id=$1 limit 1
`

func (q *Queries) GetProfileContent(ctx context.Context, id uuid.UUID) (ProfileContent, error) {
	row := q.queryRow(ctx, q.getProfileContentStmt, getProfileContent, id)
	var i ProfileContent
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.DisplayTitle,
		&i.Description,
		&i.Url,
		&i.ProfileID,
		&i.CallToActionID,
		&i.ContentID,
		&i.Order,
	)
	return i, err
}

const getProfileContents = `-- name: GetProfileContents :many
select id, title, display_title, description, url, profile_id, call_to_action_id, content_id, "order" from profile_contents where profile_id = $1
`

func (q *Queries) GetProfileContents(ctx context.Context, profileID uuid.UUID) ([]ProfileContent, error) {
	rows, err := q.query(ctx, q.getProfileContentsStmt, getProfileContents, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ProfileContent
	for rows.Next() {
		var i ProfileContent
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.DisplayTitle,
			&i.Description,
			&i.Url,
			&i.ProfileID,
			&i.CallToActionID,
			&i.ContentID,
			&i.Order,
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

const updateProfileContent = `-- name: UpdateProfileContent :exec
update profile_contents set title=$1, display_title=$2, "description"=$3, "url"=$4, call_to_action_id=$5, "order"=$6
where id=$7
`

type UpdateProfileContentParams struct {
	Title          string    `json:"title"`
	DisplayTitle   string    `json:"display_title"`
	Description    string    `json:"description"`
	Url            string    `json:"url"`
	CallToActionID uuid.UUID `json:"call_to_action_id"`
	Order          int32     `json:"order"`
	ID             uuid.UUID `json:"id"`
}

func (q *Queries) UpdateProfileContent(ctx context.Context, arg UpdateProfileContentParams) error {
	_, err := q.exec(ctx, q.updateProfileContentStmt, updateProfileContent,
		arg.Title,
		arg.DisplayTitle,
		arg.Description,
		arg.Url,
		arg.CallToActionID,
		arg.Order,
		arg.ID,
	)
	return err
}
