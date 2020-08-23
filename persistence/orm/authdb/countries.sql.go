// Code generated by sqlc. DO NOT EDIT.
// source: countries.sql

package authdb

import (
	"context"
	"database/sql"
)

const deleteCountry = `-- name: DeleteCountry :exec
delete from countries where id == $1 or "name" == $1 or flag_image_url == $1
`

func (q *Queries) DeleteCountry(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteCountryStmt, deleteCountry, id)
	return err
}

const getCountries = `-- name: GetCountries :many
select id, name, flag_image_url from countries
`

func (q *Queries) GetCountries(ctx context.Context) ([]Country, error) {
	rows, err := q.query(ctx, q.getCountriesStmt, getCountries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Country
	for rows.Next() {
		var i Country
		if err := rows.Scan(&i.ID, &i.Name, &i.FlagImageUrl); err != nil {
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

const getCountry = `-- name: GetCountry :one
select id, name, flag_image_url from countries where id == $1 or "name" == $1 limit 1
`

func (q *Queries) GetCountry(ctx context.Context, id int64) (Country, error) {
	row := q.queryRow(ctx, q.getCountryStmt, getCountry, id)
	var i Country
	err := row.Scan(&i.ID, &i.Name, &i.FlagImageUrl)
	return i, err
}

const updateCountry = `-- name: UpdateCountry :one
update countries set "name" = $1, flag_image_url = $2 where "name" == $3
returning id, name, flag_image_url
`

type UpdateCountryParams struct {
	Name         string         `json:"name"`
	FlagImageUrl sql.NullString `json:"flag_image_url"`
	Name_2       string         `json:"name_2"`
}

func (q *Queries) UpdateCountry(ctx context.Context, arg UpdateCountryParams) (Country, error) {
	row := q.queryRow(ctx, q.updateCountryStmt, updateCountry, arg.Name, arg.FlagImageUrl, arg.Name_2)
	var i Country
	err := row.Scan(&i.ID, &i.Name, &i.FlagImageUrl)
	return i, err
}
