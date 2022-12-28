package repositories

import (
	"github.com/jmoiron/sqlx"
	"priority-task-manager/shared/pkg/types"
)

//goland:noinspection ALL
const (
	getTaskQuery = `
	select
		uuid,
		account_hash,
		status,
		type,
		added_to_queue,
		extracted_from_queue,
		completed
	from task_stat
	where uuid = $1
	`

	addTaskQuery = `
	insert into
	task_stat(uuid, account_hash, status, type, added_to_queue, extracted_from_queue, completed)
	values(:uuid, :account_hash, :status, :type, :added_to_queue, :extracted_from_queue, :completed)
	`

	// todo как-то хитро обновлять время extracted_from_queue; нужно этого НЕ делать, если completed уже заполнено
	updateTaskQuery = `
	update task_stat
	set
		status = :status,
		added_to_queue = :added_to_queue,
		extracted_from_queue = :extracted_from_queue,
		completed = :completed
	where uuid = :uuid
	`
)

type Task struct {
	db *sqlx.DB
}

func MakeTaskRepository(db *sqlx.DB) Task {
	return Task{db: db}
}

func (t Task) Get(uuid string) (types.Task, error) {
	var task types.Task
	err := t.db.Get(&task, getTaskQuery, uuid)

	return task, err
}

func (t Task) Add(entity types.Task) error {
	_, err := t.db.NamedExec(addTaskQuery, t.taskToMap(entity))
	return err
}

func (t Task) taskToMap(entity types.Task) map[string]any {
	m := map[string]any{
		"uuid":                 entity.UUID,
		"account_hash":         entity.AccountHash,
		"status":               entity.Status,
		"type":                 entity.Type,
		"added_to_queue":       entity.AddedToQueue,
		"extracted_from_queue": entity.ExtractedFromQueue,
		"completed":            entity.Completed,
	}

	if entity.ExtractedFromQueue.IsZero() {
		m["extracted_from_queue"] = nil
	}

	if entity.Completed.IsZero() {
		m["completed"] = nil
	}

	return m
}

func (t Task) Update(entity types.Task) error {
	_, err := t.db.NamedExec(updateTaskQuery, t.taskToMap(entity))
	return err
}

func (t Task) Delete(_ types.ID) error {
	panic("not implemented")
}
