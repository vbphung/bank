package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/golang/mock/mockgen/model"
)

type Store struct {
	*Queries
	db *sql.DB
}

type TransferParams struct {
	FromAcc int64 `json:"from_id"`
	ToAcc   int64 `json:"to_id"`
	Amount  int64 `json:"amount"`
}

type TransferResult struct {
	Transfer  Transfer `json:"transfer"`
	FromAcc   Account  `json:"from_account"`
	ToAcc     Account  `json:"to_account"`
	FromEntry Entry    `json:"from_entry"`
	ToEntry   Entry    `json:"to_entry"`
}

func CreateStore(storeDb *sql.DB) *Store {
	return &Store{
		db:      storeDb,
		Queries: New(storeDb),
	}
}

func (store *Store) Transfer(ctx context.Context, args TransferParams) (TransferResult, error) {
	var res TransferResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		if res.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromID: args.FromAcc,
			ToID:   args.ToAcc,
			Amount: args.Amount,
		}); err != nil {
			return err
		}

		if res.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAcc,
			Amount:    -args.Amount,
		}); err != nil {
			return err
		}

		if res.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAcc,
			Amount:    args.Amount,
		}); err != nil {
			return err
		}

		if args.FromAcc > args.ToAcc {
			if res.FromAcc, res.ToAcc, err = transfer(
				ctx, q, args.FromAcc, args.ToAcc, args.Amount,
			); err != nil {
				return err
			}
		} else {
			if res.ToAcc, res.FromAcc, err = transfer(
				ctx, q, args.ToAcc, args.FromAcc, -args.Amount,
			); err != nil {
				return err
			}
		}

		return nil
	})

	return res, err
}

func (store *Store) execTx(ctx context.Context, txFunc func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(tx)

	if err = txFunc(queries); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("%v - %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

func transfer(
	ctx context.Context,
	q *Queries,
	id1, id2, amount int64,
) (
	acc1, acc2 Account,
	err error,
) {
	if acc1, err = q.UpdateBalance(ctx, UpdateBalanceParams{
		ID:     id1,
		Amount: -amount,
	}); err != nil {
		return
	}

	if acc2, err = q.UpdateBalance(ctx, UpdateBalanceParams{
		ID:     id2,
		Amount: amount,
	}); err != nil {
		return
	}

	return
}
