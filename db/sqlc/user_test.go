package db

import (
	"context"
	"testing"

	"github.com/scarpart/distributed-task-scheduler/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.UserID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Username, user2.Username)
}

func CreateRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: util.RandomUsername(),
		Password: util.RandomPassword(),
		Email:    util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Username, user.Username)

	require.NotZero(t, user.UserID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestUpdateUser(t *testing.T) {
	user1 := CreateRandomUser(t)

	arg := UpdateUserParams{
		ID:       user1.UserID,
		Username: util.RandomUsername(),
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(err)
	require.NotEmpty(user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.NotEqual(t, user1.Username, user2.Username)
}
