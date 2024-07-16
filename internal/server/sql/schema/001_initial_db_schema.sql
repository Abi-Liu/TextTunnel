-- +goose Up
CREATE TABLE users (
	id UUID PRIMARY KEY,
	username VARCHAR(32) UNIQUE NOT NULL,
	password VARCHAR(64) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

CREATE TABLE rooms (
	id UUID PRIMARY KEY,
	name VARCHAR(32) UNIQUE NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	creator_id UUID NOT NULL,
	owner_id UUID NOT NULL,
	FOREIGN KEY(creator_id) REFERENCES users(id),
	FOREIGN KEY(owner_id) REFERENCES users(id)
);

CREATE TABLE messages (
	id UUID PRIMARY KEY,
	content VARCHAR(1500) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	sender_id UUID NOT NULL,
	room_id UUID NOT NULL,
	FOREIGN KEY(sender_id) REFERENCES users(id),
	FOREIGN KEY(room_id) REFERENCES rooms(id)
);


-- +goose Down
DROP TABLE messages;
DROP TABLE rooms;
DROP TABLE users;
