alter table sessions
drop constraint fk_sessions_email;

drop table if exists sessions;
