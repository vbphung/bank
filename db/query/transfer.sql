-- name: CreateTransfer :one
insert into transfers (from_id, to_id, amount)
values ($1, $2, $3)
returning *;

-- name: ReadTransfer :one
select * from transfers
where id = $1
limit 1;

-- name: ReadAllTransfers :many
select * from transfers
where from_id = $1 or to_id = $2
limit $3
offset $4;
