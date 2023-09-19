-- Initial public schema relates to Library 0.x

SET
statement_timeout = 0;
SET
lock_timeout = 0;
SET
idle_in_transaction_session_timeout = 0;
SET
client_encoding = 'UTF8';
SET
standard_conforming_strings = on;
SET
client_min_messages = warning;
SET
row_security = off;

CREATE
EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;

SET
search_path = public, pg_catalog;
SET
default_tablespace = '';

-- posts
CREATE TABLE posts
(
    id         bigint,
    title      text        NOT NULL,
    content    text        NOT NULL,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz null,
    deleted_at timestamptz null,
    CONSTRAINT posts_pk PRIMARY KEY (id)
);


-- comments
CREATE TABLE comments
(
    id      bigint,
    post_id uuid NOT NULL,
    text    text NOT NULL,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz null,
    deleted_at timestamptz null,
    CONSTRAINT comments_pk PRIMARY KEY (id),
    CONSTRAINT fk_comments_post_id FOREIGN KEY (post_id)
        REFERENCES posts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);


CREATE
EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE users
(
    id            bigint,
    username      text NOT NULL UNIQUE,
    user_password text NOT NULL,
    user_role     text NOT NULL,
    access_token  text,
    CONSTRAINT users_pk PRIMARY KEY (id)
);
CREATE INDEX user_access_token
    ON users (access_token);
INSERT INTO users(username, user_password, user_role)
VALUES
    ('admin', crypt('admin', gen_salt('bf')), 'admin');

DROP TABLE comments, users,posts;

INSERT INTO comments(text, post_id)
VALUES
    ('fitst comment', 1);

INSERT INTO comments(text, post_id)
VALUES
    ('second comment', 1);

INSERT INTO comments(text, post_id)
VALUES
    ('third comment', 2);
