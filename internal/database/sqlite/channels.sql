-- name: InsertRegularChannel :one
INSERT INTO
    channels (name, type, room_name)
VALUES
    (?, 'channel', NULL) RETURNING *;

-- name: InsertRoomChannel :one
INSERT INTO
    channels (name, type, room_name)
VALUES
    (?, 'room', ?) RETURNING *;

-- name: InsertChannelSubscribe :one
INSERT INTO
    channel_subscribers (channel_id, user_id)
VALUES
    (?, ?) RETURNING *;

-- name: GetChannels :many
SELECT
    *
FROM
    channels;

-- name: GetChannel :one
SELECT
    *
FROM
    channels
WHERE
    type = ?
    AND name = ?;

-- name: GetChannelByID :one
SELECT
    *
FROM
    channels
WHERE
    id = ?;

-- name: GetChannelSubscribers :many
SELECT
    u.id,
    u.name,
    u.apikey,
    u.room,
    u.role
FROM
    users u
    JOIN channel_subscribers cs ON u.id = cs.user_id
    JOIN channels c ON cs.channel_id = c.id
WHERE
    c.name = ?;