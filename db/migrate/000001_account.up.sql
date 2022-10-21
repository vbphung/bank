create table accounts (
    id bigserial primary key,
    email varchar(255) not null,
    password varchar not null,
    balance bigint not null default 0,
    password_changed_at timestamptz not null default (now()),
    created_at timestamptz not null default (now())
);

alter table accounts
add constraint u_accounts_email
unique (email);

create index idx_accounts_email
on accounts(email);
