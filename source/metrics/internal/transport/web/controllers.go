package web

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"priority-task-manager/shared/pkg/services/settings"
)

func MakeControllers(webSettings settings.WebSettings) {
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", webSettings.Domain, webSettings.Port), nil)
	if err != nil {
		log.Fatalf("Unable to start servig: %v", err)
	}
}
