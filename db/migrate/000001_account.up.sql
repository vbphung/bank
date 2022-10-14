create table accounts (
    id bigserial primary key,
    full_name varchar(255) not null,
    password varchar not null,
    balance bigint not null default 0,
    password_changed_at timestamptz not null default (now()),
    created_at timestamptz not null default (now())
);
