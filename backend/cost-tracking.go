package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"log"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func getPerUserLimits() (int64, int64) {
	ctx := context.Background()
	dailyLimit, _ := rdb.Get(ctx, "config:per_user_daily_limit").Int64()
	monthlyLimit, _ := rdb.Get(ctx, "config:per_user_monthly_limit").Int64()

	return dailyLimit, monthlyLimit
}

// InitRedis initializes the Redis client
func InitRedis() {
	redis_db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("invalid redis db id '%s'", os.Getenv("REDIS_DB"))
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redis_db,
	})

	ctx := context.Background()
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = rdb.Ping(pingCtx).Result()
	if err != nil {
		log.Fatalf("Redis is NOT available: %v", err)
	}
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

	dailyLim, monthlyLim := getPerUserLimits()

	if dailyLim > 0 && currentSpend >= dailyLim {
		return errors.New("daily budget exceeded")
	}
	if monthlyLim > 0 && currentSpend >= monthlyLim {
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
