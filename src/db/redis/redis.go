package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"short-links/db"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisDb struct {
	store *redis.Client
}

func (m *RedisDb) Get(shortLink string) (*db.Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	v, err := m.store.Get(ctx, buildKey(shortLink)).Result()
	if err != nil {
		return nil, fmt.Errorf("short link does not exist: %s", err)
	}
	l := &db.Link{}
	if err := json.Unmarshal([]byte(v), l); err != nil {
		return nil, fmt.Errorf("unable to unmarshal value from redis: %s", err)
	}
	return l, nil
}

func (m *RedisDb) Create(l *db.Link) (*db.Link, error) {
	l.Id = buildKey(l.ShortLink)
	v, err := json.Marshal(l)
	if err != nil {
		return l, fmt.Errorf("unable to serialize link object: %s", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = m.store.Set(ctx, l.Id.(string), v, 0).Err()
	return l, err
}

func (m *RedisDb) Update(l *db.Link) error {
	if l.Id.(string) != buildKey(l.ShortLink) {
		return fmt.Errorf(
			"the link %q has an invalid id: got=%q want=%q",
			l.ShortLink, l.Id, buildKey(l.ShortLink),
		)
	}
	_, err := m.Create(l)
	return err
}

func NewRedisDb(u *url.URL) (db.Db, error) {
	opt, err := redis.ParseURL(u.String())
	if err != nil {
		return nil, fmt.Errorf("unable to parse redis url %q: %s", u.String(), err)
	}
	return &RedisDb{
		store: redis.NewClient(opt),
	}, nil
}

func buildKey(shortLink string) string {
	return "LINK_" + shortLink
}

func init() {
	db.RegisterDbBuilder("redis", NewRedisDb)
	db.RegisterDbBuilder("unix", NewRedisDb)
}
