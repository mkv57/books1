CREATE TABLE users
(
    name_id serial PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(50) NOT NULL,
    email VARCHAR(300) UNIQUE NOT NULL
);
CREATE TABLE books
(
id serial PRIMARY KEY,
title VARCHAR (50) UNIQUE NOT NULL,
year INTEGER NOT NULL,
created_at time,
updated_at time
);
CREATE TABLE books1
(
id_id serial PRIMARY KEY,
title_id VARCHAR (50) UNIQUE NOT NULL,
year_id INTEGER NOT NULL
);
CREATE TABLE books4
(
id_id_4 serial PRIMARY KEY,
title_id_4 VARCHAR (50) UNIQUE NOT NULL,
year_id_4 INTEGER NOT NULL
);
CREATE TABLE books3
(
id_id_2 serial PRIMARY KEY,
title_id_2 VARCHAR (50) UNIQUE NOT NULL,
year_id_2 INTEGER NOT NULL
);


