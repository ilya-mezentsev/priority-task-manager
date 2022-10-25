package configs

import "priority-task-manager/shared/pkg/services/settings"

func ParseConfigs(path string) (Settings, error) {
	var s Settings
	err := settings.Parse(path, &s)

	return s, err
}
