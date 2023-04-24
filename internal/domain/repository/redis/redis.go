package redis_repository

import (
	"context"
	"time"

	"github.com/MrBorisT/go_url_shortener/internal/config"
	"github.com/redis/go-redis/v9"
)

type repo struct {
	client *redis.Client
	TTL    time.Duration
}

func NewInMemoryRepository() *repo {
	repo := &repo{
		client: redis.NewClient(&redis.Options{
			Addr:     config.ConfigData.Redis.Addr,
			Password: config.ConfigData.Redis.Password,
			DB:       0, // use default,
		}),
		TTL: time.Duration(config.ConfigData.Redis.TTL) * time.Minute,
	}
	return repo
}

func (r *repo) ExistsShortURL(fullURL string) *string {
	shortURL, err := r.client.Get(context.TODO(), fullURL).Result()
	if err != nil {
		if err != redis.Nil {
			panic(err)
		}
		return nil
	}
	return &shortURL
}

func (r *repo) SaveURL(shortURL, fullURL string) error {
	if err := r.client.Set(context.TODO(), shortURL, fullURL, r.TTL).Err(); err != nil {
		return err
	}
	if err := r.client.Set(context.TODO(), fullURL, shortURL, r.TTL).Err(); err != nil {
		return err
	}

	return nil
}

func (r *repo) GetURL(shortURL string) (string, error) {
	fullURL, err := r.client.Get(context.TODO(), shortURL).Result()
	if err != nil {
		return "", err
	}
	return fullURL, nil
}
