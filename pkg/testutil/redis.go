package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupRedisContainer(t testing.TB) testcontainers.Container {
	ctx := context.Background()
	redisContainer, err := redis.Run(ctx,
		"docker.io/redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(5*time.Second)),
	)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, redisContainer.Terminate(ctx))
	})

	return redisContainer
}

func CreateAsynqClient(t testing.TB) (*asynq.Client, *asynq.Inspector) {
	container := SetupRedisContainer(t)
	host, _ := container.Host(context.Background())
	port, _ := container.MappedPort(context.Background(), "6379")

	redisAddr := fmt.Sprintf("%s:%d", host, port.Int())
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: redisAddr})

	t.Cleanup(func() {
		client.Close()
	})

	return client, inspector
}
