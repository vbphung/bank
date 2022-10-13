package db

import (
	"context"
	"testing"

	"github.com/herbi-dino/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createAcc(t)
}

func TestReadAccount(t *testing.T) {
	acc := createAcc(t)

	readAcc, err := testQueries.ReadAccount(context.Background(), acc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, readAcc)

	require.Equal(t, readAcc.FullName, acc.FullName)
	require.Equal(t, readAcc.Balance, acc.Balance)
	require.Equal(t, readAcc.CreatedAt, acc.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	acc := createAcc(t)

	args := UpdateAccountParams{
		ID:      acc.ID,
		Balance: utils.RandomAmount(1000, 2000),
	}

	updatedAcc, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)

	require.Equal(t, updatedAcc.ID, acc.ID)

	require.Equal(t, updatedAcc.Balance, args.Balance)
}

func TestDeleteAccount(t *testing.T) {
	acc := createAcc(t)

	deletedAcc, err := testQueries.DeleteAccount(context.Background(), acc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, deletedAcc)

	require.Equal(t, deletedAcc.ID, acc.ID)
}

func createAcc(t *testing.T) Account {
	args := CreateAccountParams{
		FullName: utils.RandomFullName(),
		Balance:  utils.RandomAmount(100, 1000),
	}

	acc, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	return acc
}
