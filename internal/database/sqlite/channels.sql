-- name: InsertGroup :one
INSERT INTO
    channels (name)
VALUES
    (?) RETURNING *;

-- name: GetGroupByID :one
SELECT
    *
FROM
    channels
WHERE
    id = ?;

-- name: DeleteGroup :exec
DELETE FROM channels
WHERE
    id = ?;

-- name: SubscribeToGroup :exec
INSERT INTO
    channel_subscribers (channel_id, user_id)
VALUES
    (?, ?) ON CONFLICT DO NOTHING;

-- name: GetSubscribersByGroupID :many
SELECT
    users.*
FROM
    users
    JOIN channel_subscribers ON users.id = channel_subscribers.user_id
WHERE
    channel_subscribers.channel_id = ?;

-- name: UnsubscribeFromGroup :exec
DELETE FROM channel_subscribers
WHERE
    channel_id = ?
    AND user_id = ?;

-- name: GetGroupByUserID :many
SELECT
    channels.*
FROM
    channels
    JOIN channel_subscribers ON channels.id = channel_subscribers.channel_id
WHERE
    channel_subscribers.user_id = ?;