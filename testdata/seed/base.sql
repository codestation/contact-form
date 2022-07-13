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

insert into test_profiles (id, created_at, updated_at, external_id, avatar)
    overriding system value
values (1, now(), now(), '00000000-0000-0000-0000-000000000000'::uuid, 'https://example.com/avatar.jpg')
on conflict do nothing;

insert into test_users (created_at, updated_at, name, external_id, profile_id)
values (now(), now(), 'John Doe', '00000000-0000-0000-0000-000000000000'::uuid, 1);
