
-- name: InsertChannel :one
INSERT INTO
    channels (name)
VALUES
    (?) RETURNING *;

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

