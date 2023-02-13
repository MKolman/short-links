package redis

import (
	"fmt"
	"net/url"
	"short-links/db/internal"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func TestRedisDb(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(fmt.Errorf("cannot start redis server: %w", err))
	}
	defer mr.Close()
	u := &url.URL{Scheme: "redis", Host: mr.Addr()}
	rdb, err := NewRedisDb(u)
	if err != nil {
		panic(fmt.Errorf("cannot start redis client: %w", err))
	}
	internal.TestDbAll(t, rdb)
}
