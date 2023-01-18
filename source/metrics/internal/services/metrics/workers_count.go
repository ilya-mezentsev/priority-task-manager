package metrics

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"priority-task-manager/shared/pkg/types"
	"strconv"
)

var (
	ctx          = context.Background()
	workersCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "workers_count",
	})
)

type WorkersCountService struct {
	redisClient *redis.Client
}

func MakeWorkersCountService(redisClient *redis.Client) WorkersCountService {
	return WorkersCountService{redisClient: redisClient}
}

func (wcs WorkersCountService) UpdateWorkersCount(types.Role) error {
	val, err := wcs.redisClient.Get(ctx, "max-workers-count").Result()
	if err != nil {
		return err
	}

	count, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return err
	}

	workersCount.Set(float64(count))

	return nil
}
