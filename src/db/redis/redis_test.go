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
		panic(fmt.Errorf("Cannot start redis server: %s", err))
	}
	defer mr.Close()
	u := &url.URL{Scheme: "redis", Host: mr.Addr()}
	rdb, err := NewRedisDb(u)
	if err != nil {
		panic(fmt.Errorf("Cannot start redis client: %s", err))
	}
	internal.TestDbAll(t, rdb)
}
