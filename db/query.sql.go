// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package db

import (
	"context"
	"database/sql"
	"strings"
)

const getExamples = `-- name: GetExamples :many
SELECT id, mean_id, example FROM examples e
WHERE e.mean_id IN (/*SLICE:mean_ids*/?)
`

func (q *Queries) GetExamples(ctx context.Context, meanIds []sql.NullInt32) ([]Example, error) {
	query := getExamples
	var queryParams []interface{}
	if len(meanIds) > 0 {
		for _, v := range meanIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:mean_ids*/?", strings.Repeat(",?", len(meanIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:mean_ids*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Example
	for rows.Next() {
		var i Example
		if err := rows.Scan(&i.ID, &i.MeanID, &i.Example); err != nil {
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

const getMeans = `-- name: GetMeans :many
SELECT id, part_id, meaning, level FROM means m
WHERE m.part_id IN (/*SLICE:part_ids*/?)
`

func (q *Queries) GetMeans(ctx context.Context, partIds []sql.NullInt32) ([]Mean, error) {
	query := getMeans
	var queryParams []interface{}
	if len(partIds) > 0 {
		for _, v := range partIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:part_ids*/?", strings.Repeat(",?", len(partIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:part_ids*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Mean
	for rows.Next() {
		var i Mean
		if err := rows.Scan(
			&i.ID,
			&i.PartID,
			&i.Meaning,
			&i.Level,
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

const getVobsCollection = `-- name: GetVobsCollection :many
SELECT v.id, v.word FROM
collections c, collection_words cw, vobs v
WHERE c.id = cw.collection_id 
AND cw.vob_id = v.id
AND c.id = ? ORDER BY RAND()
`

func (q *Queries) GetVobsCollection(ctx context.Context, id int32) ([]Vob, error) {
	rows, err := q.db.QueryContext(ctx, getVobsCollection, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Vob
	for rows.Next() {
		var i Vob
		if err := rows.Scan(&i.ID, &i.Word); err != nil {
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

const getVobsRandom = `-- name: GetVobsRandom :one
SELECT GetVobsRandom(?, ?)
`

type GetVobsRandomParams struct {
	Getvobsrandom   interface{}
	Getvobsrandom_2 interface{}
}

func (q *Queries) GetVobsRandom(ctx context.Context, arg GetVobsRandomParams) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getVobsRandom, arg.Getvobsrandom, arg.Getvobsrandom_2)
	var getvobsrandom interface{}
	err := row.Scan(&getvobsrandom)
	return getvobsrandom, err
}

const getWord = `-- name: GetWord :one
SELECT id, word FROM vobs WHERE word = ?
`

func (q *Queries) GetWord(ctx context.Context, word sql.NullString) (Vob, error) {
	row := q.db.QueryRowContext(ctx, getWord, word)
	var i Vob
	err := row.Scan(&i.ID, &i.Word)
	return i, err
}

const setWord = `-- name: SetWord :exec
CALL SetWord(?,?)
`

type SetWordParams struct {
	Setword   interface{}
	Setword_2 interface{}
}

func (q *Queries) SetWord(ctx context.Context, arg SetWordParams) error {
	_, err := q.db.ExecContext(ctx, setWord, arg.Setword, arg.Setword_2)
	return err
}
