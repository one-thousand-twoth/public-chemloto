-- name: InsertUser :one
INSERT INTO
    users (name, apikey, room, role)
VALUES
    (?, ?, ?, ?) RETURNING *;

-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = ?;

-- name: GetUserByName :one
SELECT
    *
FROM
    users
WHERE
    name = ?;


-- name: UpdateUserByID :exec
UPDATE users
SET
    name = ?,
    apikey = ?,
    room = ?,
    role = ?
WHERE id = ?;

-- name: UpdateUserRoom :exec
UPDATE users
SET
    room = ?
WHERE
    id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
    id = ?;

-- name: PatchUserRole :one
UPDATE users
SET
    role = ?
WHERE
    name = ? RETURNING *;

-- name: GetUsers :many
SELECT
    *
FROM
    users;

-- name: GetUserByApikey :one
SELECT
    *
FROM
    users
WHERE
    apikey = ?;
