package db

import (
	"context"
	"testing"

	"github.com/herbi-dino/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	store := CreateStore(testDb)

	fromAcc := createAcc(t)
	toAcc := createAcc(t)

	amounts := make(chan int64)

	errors := make(chan error)
	results := make(chan TransferResult)

	times := 10

	for i := 0; i < times; i++ {
		go func() {
			amount := utils.RandomBalance(100, 1000)

			res, err := store.Transfer(context.Background(), TransferParams{
				FromAcc: fromAcc.ID,
				ToAcc:   toAcc.ID,
				Amount:  amount,
			})

			amounts <- amount

			errors <- err
			results <- res
		}()
	}

	for i := 0; i < times; i++ {
		amount := <-amounts

		err := <-errors
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

		// test transfer
		trf := res.Transfer

		require.NotEmpty(t, trf)

		require.Equal(t, trf.FromID, fromAcc.ID)
		require.Equal(t, trf.ToID, toAcc.ID)
		require.Equal(t, trf.Amount, amount)

		_, err = store.ReadTransfer(context.Background(), trf.ID)
		require.NoError(t, err)

		// test from entry
		fromEtr := res.FromEntry

		require.NotEmpty(t, fromEtr)

		require.Equal(t, fromEtr.AccountID, fromAcc.ID)
		require.Equal(t, fromEtr.Amount, -amount)

		_, err = store.ReadEntry(context.Background(), fromEtr.ID)
		require.NoError(t, err)

		// test to entry
		toEtr := res.ToEntry

		require.NotEmpty(t, toEtr)

		require.Equal(t, toEtr.AccountID, toAcc.ID)
		require.Equal(t, toEtr.Amount, amount)

		_, err = store.ReadEntry(context.Background(), toEtr.ID)
		require.NoError(t, err)
	}
}
