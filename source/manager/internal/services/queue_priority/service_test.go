package queue_priority

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"priority-task-manager/manager/internal/services/permission"
	"testing"
)

var (
	versionManager = permission.MakeVersionManager(permission.Settings{
		DefaultVersionId: "default-version-id",
		LoadLevels: []permission.LoadLevel{
			{
				VersionID: "default-version-id",
			},
		},
	})
)

type MockPermissionResolver struct {
	calledWithPriorityKey string
}

func (m *MockPermissionResolver) Get(key permission.ResolveRequest) (bool, error) {
	m.calledWithPriorityKey = key.ResourceId

	if key.RoleId == "error" {
		return false, errors.New("some-error")
	} else if key.RoleId == "not-found" {
		return false, nil
	} else {
		return true, nil
	}
}

func TestService_DetermineMaxPriority_Success(t *testing.T) {
	service := regularService()

	priority, err := service.DetermineMaxPriority("bronze-client")
	assert.Nil(t, err)
	assert.Equal(t, 100, priority)
	assert.Equal(t, "priority-100", service.permissionResolver.(*MockPermissionResolver).calledWithPriorityKey)
}

func TestService_DetermineMaxPriority_Error(t *testing.T) {
	service := regularService()

	_, err := service.DetermineMaxPriority("error")
	assert.NotNil(t, err)
	assert.Equal(t, "priority-100", service.permissionResolver.(*MockPermissionResolver).calledWithPriorityKey)
}

func TestService_DetermineMaxPriority_NotFound(t *testing.T) {
	service := regularService()

	priority, err := service.DetermineMaxPriority("not-found")
	assert.Nil(t, err)
	assert.Equal(t, 100, priority)

	// Тут именно 90, т.к. это последний элемент в списке
	assert.Equal(t, "priority-90", service.permissionResolver.(*MockPermissionResolver).calledWithPriorityKey)
}

func regularService() Service {
	return MakeService(
		[]Settings{
			{
				Key:   "priority-90",
				Value: 90,
			},
			{
				Key:   "priority-100",
				Value: 100,
			},
		},
		versionManager,
		&MockPermissionResolver{},
	)
}
