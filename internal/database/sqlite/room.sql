-- name: InsertRoom :one
INSERT INTO
    rooms (name, engine)
VALUES
    (?, ?) RETURNING *;

-- name: GetRoomByName :one
SELECT
    *
FROM
    rooms
WHERE
    name = ?;

-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE
    name = ?;

-- name: GetRooms :many
SELECT
    *
FROM
    rooms;