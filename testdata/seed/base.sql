create table test_profiles
(
    id          integer generated always as identity,
    created_at  timestamptz not null,
    updated_at  timestamptz not null,
    deleted_at  timestamptz,
    avatar text,
    primary key (id)
);

create table test_users
(
    id          integer generated always as identity,
    created_at  timestamptz not null,
    updated_at  timestamptz not null,
    deleted_at  timestamptz,
    name        text        not null,
    profile_id  integer not null,
    primary key (id),
    unique (name),
    CONSTRAINT fk_users_profile FOREIGN KEY (profile_id) REFERENCES test_profiles (id)
);



insert into test_profiles (created_at, updated_at, avatar)
values (now(), now(), 'https://example.com/avatar.jpg');
insert into test_profiles (created_at, updated_at, avatar)
values (now(), now(), 'https://example.com/avatar.jpg');
insert into test_profiles (created_at, updated_at, avatar)
values (now(), now(), 'https://example.com/avatar.jpg');
insert into test_profiles (created_at, updated_at, avatar)
values (now(), now(), 'https://example.com/avatar.jpg');
insert into test_profiles (created_at, updated_at, avatar)
values (now(), now(), 'https://example.com/avatar.jpg');

insert into test_users (created_at, updated_at, name, profile_id)
values (now(), now(), 'John Doe 1', 1);
insert into test_users (created_at, updated_at, name, profile_id)
values (now(), now(), 'John Doe 2', 2);
insert into test_users (created_at, updated_at, name, profile_id)
values (now(), now(), 'John Doe 3', 3);
insert into test_users (created_at, updated_at, name, profile_id)
values (now(), now(), 'John Doe 4', 4);
insert into test_users (created_at, updated_at, name, profile_id)
values (now(), now(), 'John Doe 5', 5);
