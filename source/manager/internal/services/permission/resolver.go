package permission

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Resolver сущность для работы с сервисом ограничения доступа.
type Resolver struct {
	settings Settings
}

func MakeResolver(settings Settings) Resolver {
	return Resolver{settings: settings}
}

// Get метод для получения флага разрешенности выполнения какого-то действия актором с определенной ролью.
// В случае любой ошибки возвращаем true, т.к., например, падение сервиса ограничения доступа приведет к полной недоступности функционала бекенда.
func (resolver Resolver) Get(key ResolveRequest) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, resolver.settings.PermissionResolverURL, nil)
	if err != nil {
		return true, err
	}

	resolver.makeQuery(key, req)

	client := http.Client{Timeout: time.Duration(resolver.settings.RequestTimeout) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return true, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"request":         key,
			"response_status": res.StatusCode,
		}).Errorf("Unable to resolve permission - unsuccessful response status")

		return true, nil
	}

	var r Response
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		log.Errorf("Unable to decode response body: %v", err)

		return true, nil
	}

	if r.Status == okStatus {
		return r.Data.Effect == permitEffect, nil
	} else {
		log.WithFields(log.Fields{
			"response_data": r.Data,
		}).Error("Got error from permission resolver")

		return true, nil
	}
}

func (resolver Resolver) makeQuery(request ResolveRequest, httpReq *http.Request) {
	q := httpReq.URL.Query()

	q.Add("roleId", request.RoleId)
	q.Add("resourceId", request.ResourceId)
	q.Add("operation", request.Operation)
	if request.RolesVersionId != "" {
		q.Add("rolesVersionId", request.RolesVersionId)
	}

	httpReq.URL.RawQuery = q.Encode()
}
