package configs

import "priority-task-manager/shared/pkg/services/settings"

type (
	Settings struct {
		Web settings.WebSettings `json:"web-metrics"`
		DB  settings.DBSettings  `json:"db"`
	}
)
