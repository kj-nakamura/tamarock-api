CREATE TABLE IF NOT EXISTS artist_infos(
	id serial PRIMARY KEY,
	artist_id VARCHAR (50),
	name VARCHAR (255) UNIQUE NOT NULL,
	url VARCHAR (255),
	twitter_id VARCHAR (255),
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at DATETIME DEFAULT NULL
)