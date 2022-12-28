package worker

import (
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

	// todo понять, как передавать время "выполнения" в данных задачи (task.Data)
	time.Sleep(time.Second)

	task.Completed = time.Now()
	task.Status = types.SuccessfullyDoneStatus
	e.updateTaskInStorage(task)
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
