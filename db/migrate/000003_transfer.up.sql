create table transfers (
    id bigserial primary key,
    from_id bigserial not null,
    to_id bigserial not null,
    amount bigint not null,
    created_at timestamptz not null default (now())
);

alter table transfers
add constraint fk_tranfers_from_id
foreign key (from_id) references accounts(id);

alter table transfers
add constraint fk_transfers_to_id
foreign key (to_id) references accounts(id);

create index idx_transfers_from_id
on transfers(from_id);

create index idx_transfers_to_id
on transfers(to_id);

create index idx_transfers
on transfers(from_id, to_id);
