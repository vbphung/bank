create table entries (
    id bigserial primary key,
    account_id bigserial not null,
    amount bigint not null,
    created_at timestamptz not null default (now())
);

alter table entries
add constraint fk_entries_account_id
foreign key (account_id) references accounts(id);

create index idx_entries_account_id
on entries(account_id);
