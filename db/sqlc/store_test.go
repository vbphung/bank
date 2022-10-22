package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vbph/bank/utils"
)

func TestTransfer(t *testing.T) {
	store := CreateStore(testDb)

	initFromAcc := createAcc(t)
	initToAcc := createAcc(t)

	amounts := make(chan int64)

	errors := make(chan error)
	results := make(chan TransferResult)

	numOfCases := 10

	for i := 0; i < numOfCases; i++ {
		go func() {
			amount := utils.RandomAmount(100, 1000)

			res, err := store.Transfer(context.Background(), TransferParams{
				FromAcc: initFromAcc.ID,
				ToAcc:   initToAcc.ID,
				Amount:  amount,
			})

			amounts <- amount

			errors <- err
			results <- res
		}()
	}

	totalAmount := int64(0)

	for i := 0; i < numOfCases; i++ {
		amount := <-amounts
		totalAmount += amount

		err := <-errors
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

		trf := res.Transfer
		require.NotEmpty(t, trf)
		require.Equal(t, initFromAcc.ID, trf.FromID)
		require.Equal(t, initToAcc.ID, trf.ToID)
		require.Equal(t, amount, trf.Amount)

		fromEtr := res.FromEntry
		require.NotEmpty(t, fromEtr)
		require.Equal(t, initFromAcc.ID, fromEtr.AccountID)
		require.Equal(t, -amount, fromEtr.Amount)

		toEtr := res.ToEntry
		require.NotEmpty(t, toEtr)
		require.Equal(t, initToAcc.ID, toEtr.AccountID)
		require.Equal(t, amount, toEtr.Amount)

		trfFromAcc := res.FromAcc
		require.NotEmpty(t, trfFromAcc)
		require.Equal(t, initFromAcc.ID, trfFromAcc.ID)

		trfToAcc := res.ToAcc
		require.NotEmpty(t, trfToAcc)
		require.Equal(t, initToAcc.ID, trfToAcc.ID)

		fromDiff := initFromAcc.Balance - trfFromAcc.Balance
		toDiff := initToAcc.Balance - trfToAcc.Balance
		require.Equal(t, int64(0), fromDiff+toDiff)
	}

	finalFromAcc, err := testQueries.ReadAccount(context.Background(), initFromAcc.Email)
	require.NoError(t, err)
	require.Equal(t, initFromAcc.Balance-totalAmount, finalFromAcc.Balance)

	finalToAcc, err := testQueries.ReadAccount(context.Background(), initToAcc.Email)
	require.NoError(t, err)
	require.Equal(t, initToAcc.Balance+totalAmount, finalToAcc.Balance)
}

func TestDeadlock(t *testing.T) {
	store := CreateStore(testDb)

	acc1 := createAcc(t)
	acc2 := createAcc(t)

	errors := make(chan error)

	numOfCases := 10
	amount := utils.RandomAmount(100, 1000)

	for i := 0; i < numOfCases; i++ {
		fromId := acc1.ID
		toId := acc2.ID

		if i%2 == 1 {
			fromId = acc2.ID
			toId = acc1.ID
		}

		go func() {
			_, err := store.Transfer(context.Background(), TransferParams{
				FromAcc: fromId,
				ToAcc:   toId,
				Amount:  amount,
			})

			errors <- err
		}()
	}

	for i := 0; i < numOfCases; i++ {
		err := <-errors
		require.NoError(t, err)
	}

	finalAcc1, err := testQueries.ReadAccount(context.Background(), acc1.Email)
	require.NoError(t, err)
	require.Equal(t, acc1.Balance, finalAcc1.Balance)

	finalAcc2, err := testQueries.ReadAccount(context.Background(), acc2.Email)
	require.NoError(t, err)
	require.Equal(t, acc2.Balance, finalAcc2.Balance)
}
