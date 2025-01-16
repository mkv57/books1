CREATE TABLE users
(
    user_id serial PRIMARY KEY,
    password text NOT NULL,
    email text UNIQUE NOT NULL
);
CREATE TABLE books
(
id serial PRIMARY KEY,
title VARCHAR (50) UNIQUE NOT NULL,
year INTEGER NOT NULL,
created_at time,
updated_at time,
)
CREATE TABLE session
(
    id serial primary key.
    user_id int references users not null,
    token TEXT not null,
    ip TEXT not null,
    user_agent TEXT not null,
    created_at timestamp not null default now()

)



