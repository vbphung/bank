-- name: CreateSession :one
insert into sessions (id, email, refresh_token, user_agent, client_ip, expired_at)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: ReadSession :one
select * from sessions
where id = $1
limit 1;
