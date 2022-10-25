package configs

import "priority-task-manager/shared/pkg/services/settings"

type (
	Settings struct {
		Web settings.WebSettings `json:"web_metrics"`
		DB  settings.DBSettings  `json:"db"`
	}
)
