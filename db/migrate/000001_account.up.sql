create table accounts (
    id bigserial primary key,
    balance bigint not null default 0,
    created_at timestamptz not null default (now())
);
