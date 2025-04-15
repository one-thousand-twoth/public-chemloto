CREATE TABLE
    users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        apikey TEXT NOT NULL UNIQUE,
        room TEXT,
        role INTEGER NOT NULL,
        FOREIGN KEY (room) REFERENCES rooms (name) ON DELETE SET NULL
    );

CREATE TABLE
    rooms (
        name TEXT PRIMARY KEY,
        engine TEXT NOT NULL
    );

CREATE TABLE
    channels (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE
    );

CREATE TABLE
    channel_subscribers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        channel_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        UNIQUE (channel_id, user_id),
        FOREIGN KEY (channel_id) REFERENCES channels (id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );