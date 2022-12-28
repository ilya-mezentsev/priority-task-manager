package queue

import "priority-task-manager/shared/pkg/types"

type TaskReceiver interface {
	Add(task types.Task, priority int) error
}
