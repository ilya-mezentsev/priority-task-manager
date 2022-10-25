package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"priority-task-manager/manager/internal/services/permission"
)

type (
	Alert struct {
		Annotations struct {
			LatencySeconds float64 `json:"latency_seconds"`
		} `json:"annotations"`
	}

	AlertRequest struct {
		Alters []Alert `json:"alters"`
	}
)

type AlertController struct {
	service *permission.VersionManager
}

func makeAlertController(service *permission.VersionManager) AlertController {
	return AlertController{service: service}
}

func (ac AlertController) Handle(context *gin.Context) {
	var r AlertRequest
	if err := context.ShouldBindJSON(&r); err != nil {
		context.String(http.StatusBadRequest, "Invalid JSON")
		return
	}

	for _, alert := range r.Alters {
		ac.service.UpdateVersionId(alert.Annotations.LatencySeconds)
	}

	context.String(http.StatusOK, "OK")
}
