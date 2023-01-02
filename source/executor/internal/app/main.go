package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"priority-task-manager/executor/internal/configs"
	"priority-task-manager/executor/internal/services"
	"priority-task-manager/shared/pkg/services/db/connection"
	myLogger "priority-task-manager/shared/pkg/services/log"
	"strconv"
	"syscall"
)

func init() {
	myLogger.Configure()
}

func Main() {
	configsPath := os.Args[1]
	settings, err := configs.ParseConfigs(configsPath)
	if err != nil {
		log.Fatalf("Unable to parse configs by path %s, got error %v\n", configsPath, err)
	}

	processWorkersCount(&settings)

	db := connection.MustGetConnection(settings.DB)
	ss := services.MakeServices(settings, db)

	log.Info("Starting consume tasks from queue")
	go ss.TaskConsumer().StartConsume()
	waitForSignal(ss)
}

func processWorkersCount(settings *configs.Settings) {
	workersCount := os.Getenv("WORKERS_COUNT")
	if workersCount != "" {
		parsedWorkersCount, err := strconv.Atoi(workersCount)
		if err != nil {
			log.Infof("Got invalid workers count env <%s>, unable to parse: %v", workersCount, err)
		} else {
			settings.WorkersPool.MaxWorkersCount = parsedWorkersCount
		}
	}
}

func waitForSignal(ss services.Services) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs

		fmt.Println()
		log.Infof("Got signal: %v. Starting graceful shutdown", sig)

		// todo по-хорошему, нужно тут еще по тайм-ауту завершаться (см. time.AfterFunc)
		ss.TaskConsumer().StopConsume()
		ss.WorkersPool().WaitForAllDone()

		log.Infof("All tasks is done. Service ready to shutdown")

		done <- true
	}()

	<-done
	log.Info("Exiting")
}
