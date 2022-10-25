package permission

type (
	ResolveRequest struct {
		RoleId         string
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
