package types

import "time"

type Task struct {
	UUID               string         `db:"uuid" json:"uuid"`
	AccountHash        string         `db:"account_hash" json:"account_hash"`
	Type               string         `db:"type" json:"type"`
	Data               map[string]any `db:"data" json:"data"`
	Status             string         `db:"status" json:"status"`
	AddedToQueue       time.Time      `db:"added_to_queue" json:"added_to_queue"`
	ExtractedFromQueue time.Time      `db:"extracted_from_queue" json:"extracted_from_queue"`
	Completed          time.Time      `db:"completed" json:"completed"`
}

// Прямо сейчас сложно придумать, куда это положить, т.к. использоваться будет чуть ли не везде
const (
	InitialStatus          = "pending"
	InProgressStatus       = "in_progress"
	SuccessfullyDoneStatus = "completed"
)
