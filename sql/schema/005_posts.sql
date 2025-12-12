-- +goose Up
CREATE TABLE posts(
	id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	title TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL,
	description TEXT NOT NULL,
	published_at TIMESTAMP NOT NULL,
	feed_id UUID NOT NULL,
	CONSTRAINT fk_feed FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts; 
