-- +goose Up
CREATE TABLE auth_users(  
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    password_confirm VARCHAR(255),
    is_admin BOOLEAN,
    create_time TIMESTAMP NOT NULL DEFAULT now(),
    update_time TIMESTAMP
);

-- +goose Down
drop table auth_users
