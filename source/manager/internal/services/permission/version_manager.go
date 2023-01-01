package permission

import (
	"sort"
	"sync"
)

// VersionManager сущность для управления версиями, с которыми делаются запросы к сервису ограничения доступа
type VersionManager struct {
	versionId        string
	sortedLoadLevels []LoadLevel

	settings Settings

	sync.Mutex
}

func MakeVersionManager(settings Settings) *VersionManager {
	return &VersionManager{
		versionId:        settings.DefaultVersionId,
		sortedLoadLevels: sortLoadLevels(settings.LoadLevels),

		settings: settings,
	}
}

func sortLoadLevels(levels []LoadLevel) []LoadLevel {
	sortedLevels := make([]LoadLevel, len(levels))
	copy(sortedLevels, levels)

	sort.Slice(sortedLevels, func(i, j int) bool {
		return sortedLevels[i].MaxLatency < sortedLevels[j].MaxLatency
	})

	return sortedLevels
}

func (vm *VersionManager) VersionId() string {
	vm.Lock()
	defer vm.Unlock()

	return vm.versionId
}

func (vm *VersionManager) UpdateVersionId(latency float64) {
	newVersionId := ""
	for _, level := range vm.sortedLoadLevels {
		if latency < level.MaxLatency {
			newVersionId = level.VersionID
			break
		}
	}

	if newVersionId == "" {
		newVersionId = vm.settings.CriticalVersionID
	}

	vm.setVersionId(newVersionId)
}

func (vm *VersionManager) setVersionId(newVersionId string) {
	vm.Lock()
	defer vm.Unlock()

	vm.versionId = newVersionId
}
