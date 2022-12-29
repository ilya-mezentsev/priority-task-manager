package permission

import "priority-task-manager/shared/pkg/types"

type (
	ResolveRequest struct {
		RoleId         types.Role
		ResourceId     string
		Operation      string
		RolesVersionId string
	}

	Response struct {
		Status string `json:"status"`
		Data   struct {
			Effect string `json:"effect"` // для успешных ответов

			// для ответов с ошибкой
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"data"`
	}
)
