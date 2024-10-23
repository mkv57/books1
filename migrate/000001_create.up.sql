CREATE TABLE users
(
    name_id serial PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(50)  NOT NULL,
    email VARCHAR(300) UNIQUE NOT NULL,
);
CREATE TABLE books
(
id serial PRIMARY KEY,
title VARCHAR (50) UNIQUE NOT NULL,
year INTEGER NOT NULL
);
CREATE TABLE books1
(
id_id serial PRIMARY KEY,
title_id VARCHAR (50) UNIQUE NOT NULL,
year_id INTEGER NOT NULL
);

