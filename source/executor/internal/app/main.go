package app

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
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

var ctx = context.Background()

const taskExecutorPidKey = "task-executor-pid"

func init() {
	myLogger.Configure()
}

func Main() {
	configsPath := os.Args[1]
	settings, err := configs.ParseConfigs(configsPath)
	if err != nil {
		log.Fatalf("Unable to parse configs by path %s, got error %v\n", configsPath, err)
	}

	r := redisClient(settings.RedisSettings)
	mustStopExistExecutorInstance(r)
	mustStoreCurrentPid(r)

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

func redisClient(redisSettings configs.RedisSettings) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisSettings.Address,
		Password: redisSettings.Password,
		DB:       redisSettings.DB,
	})
}

func mustStopExistExecutorInstance(redisClient *redis.Client) {
	val, err := redisClient.Get(ctx, taskExecutorPidKey).Result()
	if err == redis.Nil {
		log.Info("No cached pid found")
		return
	} else if err != nil {
		log.Fatalf("Unable to get task executor pid from redis: %v", err)
	}

	pid, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Fatalf("Unable to parse cached pid: %v", err)
	}

	procExists := syscall.Kill(int(pid), 0) == nil
	if !procExists {
		log.Infof("Unable to find process with pid <%d>, skipping killing", pid)
		return
	}

	log.Infof("Found exist task executor instance with pid <%d>, killing...", pid)
	err = syscall.Kill(int(pid), 2)
	if err != nil {
		log.Fatalf("Unable to kill exists task executor instace with pid <%d>: %v", pid, err)
	}
}

func mustStoreCurrentPid(redisClient *redis.Client) {
	err := redisClient.Set(ctx, taskExecutorPidKey, os.Getpid(), 0).Err()
	if err != nil {
		log.Fatalf("Unable to store current pid to redis: %v", err)
	}
}
