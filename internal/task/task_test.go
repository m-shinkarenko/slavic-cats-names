package task

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTask(t *testing.T) {
	task, err := NewTask(20)
	require.NoError(t, err)

	taskBytes := task.ToBytes()

	gotTask, err := ParseTask(taskBytes)
	require.NoError(t, err)
	require.Equal(t, task, gotTask)

	nonce := gotTask.Solve()
	ok := task.Validate(nonce)
	require.True(t, ok)

}
