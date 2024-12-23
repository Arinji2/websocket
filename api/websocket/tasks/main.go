package websocketTasks

import (
	"context"
	"encoding/json"
	"fmt"
)

type WebsocketTask[T any] struct {
	TaskID   string `json:"task_id"`
	TaskData T      `json:"task_data"`
}

// TaskResponse interface to handle different task types
type TaskResponse interface {
	GetTaskID() string
}

// Make WebsocketTask implement TaskResponse
func (t WebsocketTask[T]) GetTaskID() string {
	return t.TaskID
}

func GetTaskData(ctx context.Context, taskID string, jsonData json.RawMessage) (TaskResponse, error) {
	fmt.Println(taskID)
	switch taskID {
	case "AUTHENTICATE_USER":
		return AuthenticateUser(ctx, taskID, jsonData)
	case "CREATE_ROOM":
		return CreateRoomTask(ctx, taskID, jsonData)
	case "JOIN_ROOM":
		return JoinRoomTask(ctx, taskID, jsonData)
	case "DELETE_ROOM":
		return DeleteRoomTask(ctx, taskID, jsonData)
	default:
		return nil, fmt.Errorf("unknown task ID: %s", taskID)
	}
}
