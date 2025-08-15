package redis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/katelinlis/NotificationService/internal/domain/model"
	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	*redis.Client
}

func (r *RedisService) GetInstanceID(ReceiverID int) (string, error) {
	return r.Get(context.Background(), "ws:"+strconv.Itoa(ReceiverID)).Result()

}

func (r *RedisService) SetInstanceID(clientID int, instanceID string) {
	r.Set(context.Background(), "ws:"+strconv.Itoa(clientID), instanceID, time.Minute).Err()
}

func (r *RedisService) DelInstanceID(clientID int) {
	r.Del(context.Background(), "ws:"+strconv.Itoa(clientID)).Err()
}

func (w *RedisService) SendPublish(instanceID string, msg model.MessageCreatedEvent) error {
	data, _ := json.Marshal(msg)
	return w.Publish(context.Background(), instanceID, data).Err()
}

func InitRedis(addr string) RedisService {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return RedisService{Client: client}
}
