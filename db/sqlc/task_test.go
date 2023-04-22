package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {
	user := CreateRandomUser(t)

	// remember to change the UserID type below to int32
	arg := CreateTaskParams{
		UserID:          int64(user.UserID),
		TaskName:        "testTask",
		TaskDescription: "some random description",
		Status:          1,
		Priority:        sql.NullInt32{},
		Command:         "sudo something",
	}

	task, err := testQueries.CreateTask(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, task)

	require.Equal(t, arg.UserID, task.UserID)
	require.Equal(t, arg.TaskName, task.TaskName)

	require.NotZero(t, task.TaskID)
	require.NotZero(t, task.CreatedAt)
}
