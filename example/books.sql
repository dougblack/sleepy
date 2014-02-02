PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE books(
	id integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	title varchar(255),
	books_url varchar(255),
	have boolean NOT NULL,
	read boolean NOT NULL,
	created_at datetime,
	updated_at datetime
);
COMMIT;
