package database

import (
	"context"
	"flag"

	"github.com/redis/go-redis/v9"
)

var (
	addr = flag.String("redis_addr", "", "")
)

type Database struct {
	rc *redis.Client
}

func New() *Database {
	return &Database{
		rc: redis.NewClient(&redis.Options{Addr: *addr}),
	}
}

func (d *Database) Get(ctx context.Context, key string) (int, error) {
	val, err := d.rc.Get(ctx, key).Int()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (d *Database) Incr(ctx context.Context, key string) error {
	return d.rc.Incr(ctx, key).Err()
}
