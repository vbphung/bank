create table sessions (
    id uuid primary key,
    email varchar not null,
    refresh_token varchar not null,
    user_agent varchar not null,
    client_ip varchar not null,
    blocked boolean not null default false,
    expired_at timestamptz not null,
    created_at timestamptz not null default (now())
);

alter table sessions
add constraint fk_sessions_email
foreign key (email) references accounts(email);
