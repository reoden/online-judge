package test

import (
	"context"
	"github.com/redis/go-redis/v9"
	"online-judge/models"
	"online-judge/utils"
	"testing"
	"time"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func TestRedisSet(t *testing.T) {
	rdb.Set(ctx, "name", "mmc", time.Second*10)
}

func TestRedisGet(t *testing.T) {
	result, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	utils.DPrintln(result)
}

func TestRedisGetByModel(t *testing.T) {
	result, err := models.RDB.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	utils.DPrintln(result)
}
