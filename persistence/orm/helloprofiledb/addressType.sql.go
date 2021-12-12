// Code generated by sqlc. DO NOT EDIT.
// source: addressType.sql

package helloprofiledb

import (
	"context"

	"github.com/google/uuid"
)

const addAddressType = `-- name: AddAddressType :one
insert into address_types(
    name
)VALUES(
    $1
) returning id, name
`

func (q *Queries) AddAddressType(ctx context.Context, name string) (AddressType, error) {
	row := q.queryRow(ctx, q.addAddressTypeStmt, addAddressType, name)
	var i AddressType
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteAddressType = `-- name: DeleteAddressType :exec
delete from address_types where id=$1
`

func (q *Queries) DeleteAddressType(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteAddressTypeStmt, deleteAddressType, id)
	return err
}

const getAddressType = `-- name: GetAddressType :one
select id, name from address_types where id=$1 or name=$1
`

func (q *Queries) GetAddressType(ctx context.Context, id uuid.UUID) (AddressType, error) {
	row := q.queryRow(ctx, q.getAddressTypeStmt, getAddressType, id)
	var i AddressType
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getAllAddressTypes = `-- name: GetAllAddressTypes :many
select id, name from address_types
`

func (q *Queries) GetAllAddressTypes(ctx context.Context) ([]AddressType, error) {
	rows, err := q.query(ctx, q.getAllAddressTypesStmt, getAllAddressTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AddressType
	for rows.Next() {
		var i AddressType
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

const updateAddressType = `-- name: UpdateAddressType :exec
update address_types set name=$1 where id=$2
`

type UpdateAddressTypeParams struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

func (q *Queries) UpdateAddressType(ctx context.Context, arg UpdateAddressTypeParams) error {
	_, err := q.exec(ctx, q.updateAddressTypeStmt, updateAddressType, arg.Name, arg.ID)
	return err
}
