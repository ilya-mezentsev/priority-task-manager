package worker

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/types"
	"time"
)

type Executor struct {
	taskRepository repositories.Repository[types.Task, string]
}

func MakeExecutor(taskRepository repositories.Repository[types.Task, string]) Executor {
	return Executor{taskRepository: taskRepository}
}

// Exec делаем вид, что выполняем какую-то работу
func (e Executor) Exec(task types.Task) {
	task.ExtractedFromQueue = time.Now()
	task.Status = types.InProgressStatus
	e.updateTaskInStorage(task)

	e.imitateWork(task)

	task.Completed = time.Now()
	task.Status = types.SuccessfullyDoneStatus
	e.updateTaskInStorage(task)
}

func (e Executor) imitateWork(task types.Task) {
	someBuffer := make([]string, 0)
	done := time.After(time.Second * time.Duration(task.Data["exec_time"].(float64)))
	for {
		select {
		case <-done:
			return
		default:
			for i := 0; i < 1000; i++ {
				someBuffer = append(someBuffer, uuid.New().String())
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (e Executor) updateTaskInStorage(task types.Task) {
	err := e.taskRepository.Update(task)
	if err != nil {
		log.WithFields(log.Fields{
			"task": task,
		}).Errorf("Unable to update task: %v", err)

		// todo а дальше что?
	}
}
