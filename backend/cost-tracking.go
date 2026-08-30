package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

const SemanticCacheSavingsKey = "savings:semantic_cache"

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

// GetSemanticCacheSavings retrieves the total savings by the semantic cache from Redis.
// If the key is not available in Redis, it computes the savings from MongoDB,
// stores it in Redis, and returns the total.
func GetSemanticCacheSavings(ctx context.Context) (float64, error) {
	if rdb == nil {
		return ComputeSemanticCacheSavingsFromDB(ctx)
	}

	valStr, err := rdb.Get(ctx, SemanticCacheSavingsKey).Result()
	if err == nil {
		savings, parseErr := strconv.ParseFloat(valStr, 64)
		if parseErr == nil {
			return savings, nil
		}
		log.Printf("Warning: failed to parse redis savings '%s': %v", valStr, parseErr)
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("Redis error getting semantic cache savings: %v", err)
	}

	// Key not found in Redis or Redis error: compute from MongoDB
	savings, dbErr := ComputeSemanticCacheSavingsFromDB(ctx)
	if dbErr != nil {
		return 0, fmt.Errorf("failed to compute semantic cache savings from mongodb: %w", dbErr)
	}

	// Save computed savings to Redis
	if setErr := rdb.Set(ctx, SemanticCacheSavingsKey, savings, 0).Err(); setErr != nil {
		log.Printf("Warning: failed to save computed savings to redis: %v", setErr)
	}

	return savings, nil
}

// RecordCacheSavings increments the semantic cache savings in Redis when a request
// is served from cache. If provider != "nvidia", the cost saved is $1.
func RecordCacheSavings(ctx context.Context, provider string) error {
	if strings.ToLower(strings.TrimSpace(provider)) == "nvidia" {
		return nil
	}

	if rdb == nil {
		return errors.New("redis client not initialized")
	}

	exists, err := rdb.Exists(ctx, SemanticCacheSavingsKey).Result()
	if err != nil {
		log.Printf("Warning: failed to check redis key '%s': %v", SemanticCacheSavingsKey, err)
	}

	if exists == 0 {
		// Key does not exist in Redis, compute total from DB (which includes the logged request)
		savings, dbErr := ComputeSemanticCacheSavingsFromDB(ctx)
		if dbErr == nil {
			return rdb.Set(ctx, SemanticCacheSavingsKey, savings, 0).Err()
		}
		// Fallback to setting 1.0
		return rdb.Set(ctx, SemanticCacheSavingsKey, 1.0, 0).Err()
	}

	// Key exists, increment by 1.0
	_, err = rdb.IncrByFloat(ctx, SemanticCacheSavingsKey, 1.0).Result()
	return err
}
