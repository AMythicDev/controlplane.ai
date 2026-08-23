package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// DailyBudgetMicroCents represents a $5.00 daily limit (5,000,000 microcents)
const DailyBudgetMicroCents int64 = 5000000

var rdb *redis.Client

// InitRedis initializes the Redis client
func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// GetUserDailyKey generates the Redis key for a user's daily spend
func GetUserDailyKey(userID string) string {
	dateStr := time.Now().UTC().Format("2006-01-02")
	return fmt.Sprintf("spend:usr:%s:%s", userID, dateStr)
}

// CheckBudget checks if the user has exceeded their daily budget
func CheckBudget(ctx context.Context, userID string) error {
	key := GetUserDailyKey(userID)

	currentSpend, err := rdb.Get(ctx, key).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Key does not exist, spend is 0, which is fine
			return nil
		}
		// Redis error
		return fmt.Errorf("redis error checking budget: %w", err)
	}

	if currentSpend > DailyBudgetMicroCents {
		return errors.New("daily budget exceeded")
	}

	return nil
}

// RecordSpend increments the user's daily spend and sets a TTL
func RecordSpend(ctx context.Context, userID string, costMicroCents int64) error {
	key := GetUserDailyKey(userID)

	pipe := rdb.Pipeline()

	// Increment the spend
	pipe.IncrBy(ctx, key, costMicroCents)
	// Set TTL to 24 hours to clean up old keys
	pipe.Expire(ctx, key, 24*time.Hour)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to record spend: %w", err)
	}

	return nil
}

// GetCurrentSpend returns the current daily spend for a user
func GetCurrentSpend(ctx context.Context, userID string) (int64, error) {
	key := GetUserDailyKey(userID)

	currentSpend, err := rdb.Get(ctx, key).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Key does not exist, spend is 0
			return 0, nil
		}
		// Redis error
		return 0, fmt.Errorf("redis error getting spend: %w", err)
	}

	return currentSpend, nil
}
