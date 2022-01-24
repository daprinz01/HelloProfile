// Code generated by sqlc. DO NOT EDIT.
// source: contactCatgeories.sql

package helloprofiledb

import (
	"context"

	"github.com/google/uuid"
)

const addContactCategory = `-- name: AddContactCategory :one
insert into contact_categories(
    name
)VALUES(
    $1
) returning id, name
`

func (q *Queries) AddContactCategory(ctx context.Context, name string) (ContactCategory, error) {
	row := q.queryRow(ctx, q.addContactCategoryStmt, addContactCategory, name)
	var i ContactCategory
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteContactCategory = `-- name: DeleteContactCategory :exec
delete from contact_categories where name=$1
`

func (q *Queries) DeleteContactCategory(ctx context.Context, name string) error {
	_, err := q.exec(ctx, q.deleteContactCategoryStmt, deleteContactCategory, name)
	return err
}

const getAllContactCategories = `-- name: GetAllContactCategories :many
select id, name from contact_categories
`

func (q *Queries) GetAllContactCategories(ctx context.Context) ([]ContactCategory, error) {
	rows, err := q.query(ctx, q.getAllContactCategoriesStmt, getAllContactCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ContactCategory
	for rows.Next() {
		var i ContactCategory
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

const getContactCategory = `-- name: GetContactCategory :one
select id, name from contact_categories where name=$1
`

func (q *Queries) GetContactCategory(ctx context.Context, name string) (ContactCategory, error) {
	row := q.queryRow(ctx, q.getContactCategoryStmt, getContactCategory, name)
	var i ContactCategory
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const updateContactCategory = `-- name: UpdateContactCategory :exec
update contact_categories set name=$1 where id=$2
`

type UpdateContactCategoryParams struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

func (q *Queries) UpdateContactCategory(ctx context.Context, arg UpdateContactCategoryParams) error {
	_, err := q.exec(ctx, q.updateContactCategoryStmt, updateContactCategory, arg.Name, arg.ID)
	return err
}
