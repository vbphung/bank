alter table transfers
drop constraint fk_tranfers_from_id;

alter table transfers
drop constraint fk_transfers_to_id;

drop index idx_transfers_from_id;

drop index idx_transfers_to_id;

drop index idx_transfers;

drop table if exists transfers;
