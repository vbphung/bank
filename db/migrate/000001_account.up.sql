create table accounts (
    id bigserial primary key,
    full_name varchar(255) not null,
    balance bigint not null default 0,
    created_at timestamptz not null default (now())
);
