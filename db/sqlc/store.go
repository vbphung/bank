package db

import (
	"context"
	"database/sql"
	"fmt"
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

		fromAcc, err := q.ReadAccountForUpdate(ctx, args.FromAcc)
		if err != nil {
			return err
		}

		if res.FromAcc, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      fromAcc.ID,
			Balance: fromAcc.Balance - args.Amount,
		}); err != nil {
			return err
		}

		toAcc, err := q.ReadAccountForUpdate(ctx, args.ToAcc)
		if err != nil {
			return err
		}

		if res.ToAcc, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      toAcc.ID,
			Balance: toAcc.Balance + args.Amount,
		}); err != nil {
			return err
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
