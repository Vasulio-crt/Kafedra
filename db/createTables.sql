CREATE TABLE users (
	"id"	INTEGER PRIMARY KEY,
	"fio"	TEXT NOT NULL,
	"email"	TEXT NOT NULL UNIQUE,
	"password"	TEXT NOT NULL,
	"avatar" TEXT NOT NULL
);

CREATE TABLE products (
	"idP"	INTEGER,
	"name"	TEXT NOT NULL,
	"description"	TEXT NOT NULL,
	"price"	INTEGER NOT NULL,
	PRIMARY KEY("idP")
);

CREATE TABLE cart (
	"idC"	INTEGER,
	"id"	INTEGER,
	"idP"	INTEGER,
	FOREIGN KEY("id") REFERENCES "users"("id"),
	FOREIGN KEY("idP") REFERENCES "products"("idP")
	PRIMARY KEY("idC")
);