alter table entries
drop constraint fk_entries_account_id;

drop index idx_entries_account_id;

drop table if exists entries;
