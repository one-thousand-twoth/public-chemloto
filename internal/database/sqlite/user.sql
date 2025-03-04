-- name: InsertUser :one
INSERT INTO
    users (name, apikey, room, role)
VALUES
    (?, ?, ?, ?) RETURNING *;

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

-- name: GetUserSubsribtions :many
SELECT
    c.id,
    c.name
FROM
    channels c
    JOIN channel_subscribers cs ON c.id = cs.channel_id
WHERE
    cs.user_id = ?;