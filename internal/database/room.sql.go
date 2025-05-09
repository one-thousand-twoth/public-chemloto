// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: room.sql

package database

import (
	"context"
	"database/sql"
)

const deleteRoom = `-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE
    name = ?
`

func (q *Queries) DeleteRoom(ctx context.Context, name string) error {
	_, err := q.db.ExecContext(ctx, deleteRoom, name)
	return err
}

const getRoomByName = `-- name: GetRoomByName :one
SELECT
    name, engine
FROM
    rooms
WHERE
    name = ?
`

func (q *Queries) GetRoomByName(ctx context.Context, name string) (Room, error) {
	row := q.db.QueryRowContext(ctx, getRoomByName, name)
	var i Room
	err := row.Scan(&i.Name, &i.Engine)
	return i, err
}

const getRooms = `-- name: GetRooms :many
SELECT
    name, engine
FROM
    rooms
`

func (q *Queries) GetRooms(ctx context.Context) ([]Room, error) {
	rows, err := q.db.QueryContext(ctx, getRooms)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Room
	for rows.Next() {
		var i Room
		if err := rows.Scan(&i.Name, &i.Engine); err != nil {
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

const getUsersByRoom = `-- name: GetUsersByRoom :many
SELECT
    id, name, apikey, room, role
FROM
    users
WHERE
    room = ?
`

func (q *Queries) GetUsersByRoom(ctx context.Context, room sql.NullString) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsersByRoom, room)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Apikey,
			&i.Room,
			&i.Role,
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

const insertRoom = `-- name: InsertRoom :one
INSERT INTO
    rooms (name, engine)
VALUES
    (?, ?) RETURNING name, engine
`

type InsertRoomParams struct {
	Name   string
	Engine string
}

func (q *Queries) InsertRoom(ctx context.Context, arg InsertRoomParams) (Room, error) {
	row := q.db.QueryRowContext(ctx, insertRoom, arg.Name, arg.Engine)
	var i Room
	err := row.Scan(&i.Name, &i.Engine)
	return i, err
}
