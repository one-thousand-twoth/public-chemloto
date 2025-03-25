-- name: InsertRoom :one
INSERT INTO
    rooms (name, engine)
VALUES
    (?, ?) RETURNING *;

-- name: GetRooms :many
SELECT
    *
FROM
    rooms;

-- name: GetRoom :one
SELECT
    *
FROM
    rooms
WHERE
    name = ?;