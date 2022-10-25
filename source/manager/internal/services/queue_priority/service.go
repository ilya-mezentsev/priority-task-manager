package queue_priority

import (
	log "github.com/sirupsen/logrus"
	"priority-task-manager/manager/internal/services/permission"
	"priority-task-manager/shared/pkg/repositories"
	"sort"
)

type Service struct {
	queuePriorities    []Settings
	versionManager     *permission.VersionManager
	permissionResolver repositories.Reader[bool, permission.ResolveRequest]
}

func MakeService(
	queuePriorities []Settings,
	versionManager *permission.VersionManager,
	permissionResolver repositories.Reader[bool, permission.ResolveRequest],
) Service {
	return Service{
		queuePriorities:    sortQueuePriorities(queuePriorities),
		versionManager:     versionManager,
		permissionResolver: permissionResolver,
	}
}

func sortQueuePriorities(queuePriorities []Settings) []Settings {
	sortedQueuePriorities := make([]Settings, len(queuePriorities))
	copy(sortedQueuePriorities, queuePriorities)

	sort.Slice(sortedQueuePriorities, func(i, j int) bool {
		// по убыванию
		return sortedQueuePriorities[i].Value > sortedQueuePriorities[j].Value
	})

	return sortedQueuePriorities
}

func (s Service) DetermineMaxPriority(role string) (int, error) {
	for _, priority := range s.queuePriorities {
		isPriorityAvailable, err := s.permissionResolver.Get(permission.ResolveRequest{
			RoleId:         role,
			ResourceId:     priority.Key,
			Operation:      operation,
			RolesVersionId: s.versionManager.VersionId(),
		})

		if err != nil {
			return 0, err
		}

		if isPriorityAvailable {
			return priority.Value, nil
		}
	}

	log.Warnf("Unable to find available priority for role %s. Has %d priorities", role, len(s.queuePriorities))
	return s.queuePriorities[0].Value, nil
}
