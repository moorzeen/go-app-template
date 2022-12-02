create table if not exists users
(
    user_id       uuid default uuid_generate_v4() unique not null primary key,
    login         varchar unique                         not null,
    password_hash bytea                                  not null,
    active        bool default false                     not null,
    created_at    timestamp                              not null,
    updated_at    timestamp                              not null,
    last_login_at timestamp                              not null
);