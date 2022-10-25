package task

import "priority-task-manager/shared/pkg/types"

type (
	Payload struct {
		Type string         `json:"type"`
		Data map[string]any `json:"data"`
	}

	Request struct {
		Account types.Account
		Task    Payload
	}
)
