package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/barkar96/worker/libs/actor"
	"github.com/barkar96/worker/libs/logging"
)

var _ actor.Actor = (*Redis)(nil)
var _ Client = (*Redis)(nil)

type Client interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Close() error

	Do(ctx context.Context, args ...interface{}) *redis.Cmd

	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd

	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
}

type Redis struct {
	Client
}

func New(addrs []string, password string, timeout time.Duration) (*Redis, error) {
	var client redis.UniversalClient
	if len(addrs) == 1 {
		logging.Debug(context.Background(), "connecting to Redis")
		client = redis.NewClient(&redis.Options{
			Addr:         addrs[0],
			Password:     password,
			DialTimeout:  timeout,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		})
	} else {
		logging.Debug(context.Background(), "connecting to Redis")
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        addrs,
			Password:     password,
			DialTimeout:  timeout,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		})
	}
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &Redis{Client: client}, nil
}

func (r *Redis) Name() string {
	return "redis"
}

func (r *Redis) Start(ctx context.Context) error {
	t := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			r.checkConnection(ctx)
		}
	}
}

func (r *Redis) checkConnection(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := r.Ping(ctx).Err(); err != nil {
		logging.WithError(ctx, err, "failed to ping Redis")
	}
}

func (r *Redis) Stop(ctx context.Context) {
	r.Close()
}
