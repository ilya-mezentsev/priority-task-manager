package web

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"priority-task-manager/manager/internal/services/permission"
	"strconv"
)

type (
	Alert struct {
		Annotations struct {
			WaitingTime string `json:"waiting_time"`
		} `json:"annotations"`
	}

	AlertRequest struct {
		Alerts []Alert `json:"alerts"`
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
		log.Errorf("Unable to parse alert: %v", err)
		context.String(http.StatusBadRequest, "Invalid JSON")
		return
	}

	for _, alert := range r.Alerts {
		waitingTimer, err := strconv.ParseFloat(alert.Annotations.WaitingTime, 10)
		if err != nil {
			log.Errorf("Unable to parse waiting time (%s) in alert: %v", alert.Annotations.WaitingTime, err)
			context.String(http.StatusBadRequest, "Invalid waiting time in alert")
			continue
		}

		ac.service.UpdateVersionId(waitingTimer)
	}

	context.String(http.StatusOK, "OK")
}
