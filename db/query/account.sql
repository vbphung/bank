-- name: CreateAccount :one
insert into accounts (full_name, balance)
values ($1, $2)
returning *;

-- name: ReadAccount :one
select * from accounts
where id = $1
limit 1;

-- name: ReadAccountForUpdate :one
select * from accounts
where id = $1
limit 1
for no key update;

-- name: UpdateAccount :one
update accounts
set balance = $2
where id = $1
returning *;

-- name: UpdateBalance :one
update accounts
set balance = balance + sqlc.arg(amount)
where id = sqlc.arg(id)
returning *;

-- name: DeleteAccount :one
delete from accounts
where id = $1
returning *;
