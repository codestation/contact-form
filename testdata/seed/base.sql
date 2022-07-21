create table test_profiles
(
    id          integer generated always as identity,
    created_at  timestamptz not null,
    updated_at  timestamptz not null,
    deleted_at  timestamptz,
    external_id uuid        not null,
    avatar text,
    primary key (id)
);

create table test_users
(
    id          integer generated always as identity,
    created_at  timestamptz not null,
    updated_at  timestamptz not null,
    deleted_at  timestamptz,
    external_id uuid        not null,
    name        text        not null,
    profile_id  integer not null,
    primary key (id),
    unique (name),
    CONSTRAINT fk_users_profile FOREIGN KEY (profile_id) REFERENCES test_profiles (id)
);



insert into test_profiles (created_at, updated_at, external_id, avatar)
values (now(), now(), '00000000-0000-0000-0000-000000000001'::uuid, 'https://example.com/avatar.jpg');
insert into test_profiles (created_at, updated_at, external_id, avatar)
values (now(), now(), '00000000-0000-0000-0000-000000000002'::uuid, 'https://example.com/avatar.jpg');
insert into test_profiles (created_at, updated_at, external_id, avatar)
values (now(), now(), '00000000-0000-0000-0000-000000000003'::uuid, 'https://example.com/avatar.jpg');
insert into test_profiles (created_at, updated_at, external_id, avatar)
values (now(), now(), '00000000-0000-0000-0000-000000000004'::uuid, 'https://example.com/avatar.jpg');
insert into test_profiles (created_at, updated_at, external_id, avatar)
values (now(), now(), '00000000-0000-0000-0000-000000000005'::uuid, 'https://example.com/avatar.jpg');

insert into test_users (created_at, updated_at, name, external_id, profile_id)
values (now(), now(), 'John Doe 1', '00000000-0000-0000-0000-000000000001'::uuid, 1);
insert into test_users (created_at, updated_at, name, external_id, profile_id)
values (now(), now(), 'John Doe 2', '00000000-0000-0000-0000-000000000002'::uuid, 2);
insert into test_users (created_at, updated_at, name, external_id, profile_id)
values (now(), now(), 'John Doe 3', '00000000-0000-0000-0000-000000000003'::uuid, 3);
insert into test_users (created_at, updated_at, name, external_id, profile_id)
values (now(), now(), 'John Doe 4', '00000000-0000-0000-0000-000000000004'::uuid, 4);
insert into test_users (created_at, updated_at, name, external_id, profile_id)
values (now(), now(), 'John Doe 5', '00000000-0000-0000-0000-000000000005'::uuid, 5);
