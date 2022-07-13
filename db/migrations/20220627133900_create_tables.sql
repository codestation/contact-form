-- +migrate Up
create table if not exists contacts
(
    id          integer generated always as identity,
    created_at  timestamptz not null,
    updated_at  timestamptz not null,
    deleted_at  timestamptz,
    external_id uuid        not null,
    first_name  text        not null,
    last_name   text        not null,
    email       text        not null,
    message     text        not null,
    company     text,
    phone       text,
    subject     text,
    primary key (id),
    check (char_length(first_name) <= 255),
    check (char_length(last_name) <= 255),
    check (char_length(email) <= 255),
    check (char_length(subject) <= 512),
    check (char_length(message) <= 8192),
    check (char_length(phone) <= 64),
    check (char_length(company) <= 128)
);

-- +migrate Down
drop table if exists contacts;
