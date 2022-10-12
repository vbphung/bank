-- name: CreateAccount :one
insert into accounts (
    balance
) values (
    $1
) returning *;

-- name: GetAccount :one
select * from accounts
where id = $1
limit 1;


-- name: UpdateAccount :exec
update accounts
set balance = $2
where id = $1;

-- name: DeleteAccount :exec
delete from accounts
where id = $1;
