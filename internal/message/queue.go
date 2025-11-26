package message

import (
	"context"
	"encoding/json"
	"time"

	"campaign-service/internal/models"

	"github.com/redis/go-redis/v9"
)

type Queue interface {
	Publish(msg models.CampaignMessage) error
}

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue() *RedisQueue {
	return &RedisQueue{
		client: redis.NewClient(&redis.Options{
			Addr: "localhost:6380", // ← safe port
		}),
	}
}

func (q *RedisQueue) Publish(msg models.CampaignMessage) error {
	data, _ := json.Marshal(msg)
	return q.client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "campaign_messages",
		Values: map[string]interface{}{"payload": string(data)},
	}).Err()
}

func (q *RedisQueue) Consume(ctx context.Context, handler func(models.CampaignMessage) error) {
	// Create consumer group if not exists
	q.client.XGroupCreateMkStream(ctx, "campaign_messages", "workers", "0")

	for {
		streams, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "workers",
			Consumer: "worker-1",
			Streams:  []string{"campaign_messages", ">"},
			Count:    10,
			Block:    5000,
		}).Result()

		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				var campaignMsg models.CampaignMessage
				if err := json.Unmarshal([]byte(msg.Values["payload"].(string)), &campaignMsg); err != nil {
					continue
				}

				if err := handler(campaignMsg); err != nil {
					continue
				}

				q.client.XAck(ctx, "campaign_messages", "workers", msg.ID)
			}
		}
	}
}
