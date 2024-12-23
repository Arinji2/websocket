CREATE TABLE users (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    session_id VARCHAR(10) NOT NULL
);

CREATE TABLE rooms (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_by VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE CASCADE
);
CREATE TABLE players (
    id VARCHAR(10) PRIMARY KEY,
    room_id VARCHAR(10) NOT NULL,
    player_id VARCHAR(10) NOT NULL
);

