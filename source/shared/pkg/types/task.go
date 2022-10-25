package types

import "time"

type Task struct {
	UUID               string    `db:"uuid"`
	AccountHash        string    `db:"account_hash"`
	Status             string    `db:"status"`
	AddedToQueue       time.Time `db:"added_to_queue"`
	ExtractedFromQueue time.Time `db:"extracted_from_queue"`
	Completed          time.Time `db:"completed"`
}

// Прямо сейчас сложно придумать, куда это положить, т.к. использоваться будет чуть ли не везде
const (
	InitialStatus          = "pending"
	InProgressStatus       = "in_progress"
	FailedStatus           = "failed"
	SuccessfullyDoneStatus = "completed"
)
