package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var mongoClient *mongo.Client
var requestsCollection *mongo.Collection

type RequestRecord struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Endpoint       string        `bson:"endpoint" json:"endpoint"`
	Model          string        `bson:"model" json:"model"`
	Provider       string        `bson:"provider" json:"provider"`
	Prompt         string        `bson:"prompt" json:"prompt"`
	Messages       []ChatMessage `bson:"messages,omitempty" json:"messages,omitempty"`
	Response       string        `bson:"response" json:"response"`
	Confidence     *float32      `bson:"confidence" json:"confidence"`
	Toxicity       float32       `bson:"toxicity" json:"toxicity"`
	NLI            *NLIReport    `bson:"nli,omitempty" json:"nli,omitempty"`
	LatencyMs      int64         `bson:"latency_ms" json:"latency_ms"`
	CostMicrocents int64         `bson:"cost_microcents" json:"cost_microcents"`
	Cached         bool          `bson:"cached,omitempty" json:"cached,omitempty"`
	Timestamp      time.Time     `bson:"timestamp" json:"timestamp"`
}

// InitMongoDB initializes the MongoDB client and sets the requests collection.
func InitMongoDB() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://127.0.0.1:27017/controlplane"
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB is NOT available: %v", err)
	}

	mongoClient = client
	requestsCollection = mongoClient.Database("controlplane").Collection("requests")
	log.Println("Connected to MongoDB")
}

// LogRequest inserts a request record into MongoDB. Intended to be called from a goroutine.
func LogRequest(record RequestRecord) {
	if requestsCollection == nil {
		return
	}
	record.Timestamp = time.Now().UTC()
	_, err := requestsCollection.InsertOne(context.Background(), record)
	if err != nil {
		log.Printf("Failed to log request to MongoDB: %v", err)
	}
}

// FetchRequests retrieves request records sorted by timestamp descending with pagination.
func FetchRequests(limit, offset int64) ([]RequestRecord, int64, error) {
	ctx := context.Background()

	total, err := requestsCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := requestsCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var records []RequestRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, 0, err
	}

	// Return empty slice instead of nil for consistent JSON serialization
	if records == nil {
		records = []RequestRecord{}
	}

	return records, total, nil
}

// FetchRequestByID retrieves a single request record by its hex ObjectID.
func FetchRequestByID(id string) (*RequestRecord, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var record RequestRecord
	err = requestsCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: objID}}).Decode(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

// ComputeSemanticCacheSavingsFromDB calculates total savings by iterating over
// all cached requests in MongoDB. For each cached request where provider != "nvidia",
// $1 is added to the savings.
func ComputeSemanticCacheSavingsFromDB(ctx context.Context) (float64, error) {
	if requestsCollection == nil {
		return 0, errors.New("mongodb requests collection not initialized")
	}

	filter := bson.D{{Key: "cached", Value: true}}
	cursor, err := requestsCollection.Find(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to find cached requests in mongodb: %w", err)
	}
	defer cursor.Close(ctx)

	var cachedRequests []RequestRecord
	if err := cursor.All(ctx, &cachedRequests); err != nil {
		return 0, fmt.Errorf("failed to decode cached requests: %w", err)
	}

	var totalSavings float64
	for _, req := range cachedRequests {
		if strings.ToLower(strings.TrimSpace(req.Provider)) != "nvidia" {
			totalSavings += 1.0
		}
	}

	return totalSavings, nil
}

// ComputeAverageCost calculates the average cost across all logged requests in MongoDB.
// Returns average cost in dollars and in microcents.
func ComputeAverageCost(ctx context.Context) (float64, float64, error) {
	if requestsCollection == nil {
		return 0, 0, errors.New("mongodb requests collection not initialized")
	}

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalCost", Value: bson.D{{Key: "$sum", Value: "$cost_microcents"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := requestsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to aggregate average cost from mongodb: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0, 0, fmt.Errorf("failed to decode average cost aggregation results: %w", err)
	}

	if len(results) == 0 {
		return 0.0, 0.0, nil
	}

	var totalCost float64
	var count float64

	switch v := results[0]["totalCost"].(type) {
	case int32:
		totalCost = float64(v)
	case int64:
		totalCost = float64(v)
	case float64:
		totalCost = v
	case float32:
		totalCost = float64(v)
	case int:
		totalCost = float64(v)
	}

	switch v := results[0]["count"].(type) {
	case int32:
		count = float64(v)
	case int64:
		count = float64(v)
	case float64:
		count = v
	case float32:
		count = float64(v)
	case int:
		count = float64(v)
	}

	if count == 0 {
		return 0.0, 0.0, nil
	}

	avgMicrocents := totalCost / count
	avgDollars := avgMicrocents / 1000000.0

	return avgDollars, avgMicrocents, nil
}

type DailyModelCount struct {
	Date  string `bson:"date" json:"date"`
	Count int64  `bson:"count" json:"count"`
}

type ModelAnalyticsItem struct {
	Model               string            `bson:"model" json:"model"`
	Provider            string            `bson:"provider" json:"provider"`
	RequestCount        int64             `bson:"request_count" json:"request_count"`
	Percentage          float64           `bson:"percentage" json:"percentage"`
	AvgConfidence       *float64          `bson:"avg_confidence" json:"avg_confidence"`
	ConfidenceCount     int64             `bson:"confidence_count" json:"confidence_count"`
	AvgHallucination    *float64          `bson:"avg_hallucination" json:"avg_hallucination"`
	NLICount            int64             `bson:"nli_count" json:"nli_count"`
	AvgToxicity         float64           `bson:"avg_toxicity" json:"avg_toxicity"`
	ToxicityCount       int64             `bson:"toxicity_count" json:"toxicity_count"`
	TotalCostMicrocents int64             `bson:"total_cost_microcents" json:"total_cost_microcents"`
	TotalCostDollars    float64           `bson:"total_cost_dollars" json:"total_cost_dollars"`
	DailyCounts         []DailyModelCount `bson:"daily_counts" json:"daily_counts"`
}

type AnalyticsResponse struct {
	TotalRequests   int64                `json:"total_requests"`
	WeeklyRequests  int64                `json:"weekly_requests"`
	WeeklyStartDate string               `json:"weekly_start_date"`
	WeeklyEndDate   string               `json:"weekly_end_date"`
	Models          []ModelAnalyticsItem `json:"models"`
	DailyTotals     []DailyModelCount    `json:"daily_totals"`
}

// GetModelAnalytics aggregates request metrics per model over the last 7 days and all-time.
func GetModelAnalytics(ctx context.Context) (*AnalyticsResponse, error) {
	if requestsCollection == nil {
		return nil, errors.New("mongodb requests collection not initialized")
	}

	cursor, err := requestsCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to query requests: %w", err)
	}
	defer cursor.Close(ctx)

	var records []RequestRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("failed to decode requests: %w", err)
	}

	now := time.Now().UTC()
	sevenDaysAgo := now.AddDate(0, 0, -7)

	dailyTotalsMap := make(map[string]int64)
	var dateKeys []string
	for i := 6; i >= 0; i-- {
		d := now.AddDate(0, 0, -i).Format("2006-01-02")
		dateKeys = append(dateKeys, d)
		dailyTotalsMap[d] = 0
	}

	type modelAccumulator struct {
		model            string
		provider         string
		totalRequests    int64
		weeklyRequests   int64
		confidenceSum    float64
		confidenceCount  int64
		hallucinationSum float64
		nliCount         int64
		toxicitySum      float64
		toxicityCount    int64
		costMicrocents   int64
		dailyMap         map[string]int64
	}

	modelMap := make(map[string]*modelAccumulator)
	var totalRequests int64
	var weeklyRequests int64

	for _, req := range records {
		totalRequests++
		modelKey := req.Model
		if modelKey == "" {
			modelKey = "unknown"
		}
		acc, exists := modelMap[modelKey]
		if !exists {
			acc = &modelAccumulator{
				model:    modelKey,
				provider: req.Provider,
				dailyMap: make(map[string]int64),
			}
			for _, d := range dateKeys {
				acc.dailyMap[d] = 0
			}
			modelMap[modelKey] = acc
		}
		if acc.provider == "" && req.Provider != "" {
			acc.provider = req.Provider
		}

		acc.totalRequests++
		acc.costMicrocents += req.CostMicrocents

		if req.Confidence != nil {
			acc.confidenceSum += float64(*req.Confidence)
			acc.confidenceCount++
		}

		if req.NLI != nil {
			acc.hallucinationSum += float64(req.NLI.ContradictionProb)
			acc.nliCount++
		}

		acc.toxicitySum += float64(req.Toxicity)
		acc.toxicityCount++

		isWeekly := req.Timestamp.After(sevenDaysAgo)
		if isWeekly {
			weeklyRequests++
			acc.weeklyRequests++
			dateStr := req.Timestamp.Format("2006-01-02")
			if _, ok := dailyTotalsMap[dateStr]; ok {
				dailyTotalsMap[dateStr]++
			}
			if _, ok := acc.dailyMap[dateStr]; ok {
				acc.dailyMap[dateStr]++
			}
		}
	}

	var modelItems []ModelAnalyticsItem
	for _, acc := range modelMap {
		var avgConf *float64
		if acc.confidenceCount > 0 {
			val := acc.confidenceSum / float64(acc.confidenceCount)
			avgConf = &val
		}

		var avgHallucination *float64
		if acc.nliCount > 0 {
			val := acc.hallucinationSum / float64(acc.nliCount)
			avgHallucination = &val
		}

		var avgTox float64
		if acc.toxicityCount > 0 {
			avgTox = acc.toxicitySum / float64(acc.toxicityCount)
		}

		var pct float64
		if weeklyRequests > 0 {
			pct = (float64(acc.weeklyRequests) / float64(weeklyRequests)) * 100.0
		} else if totalRequests > 0 {
			pct = (float64(acc.totalRequests) / float64(totalRequests)) * 100.0
		}

		var dailyCounts []DailyModelCount
		for _, d := range dateKeys {
			dailyCounts = append(dailyCounts, DailyModelCount{
				Date:  d,
				Count: acc.dailyMap[d],
			})
		}

		effectiveCount := acc.weeklyRequests
		if weeklyRequests == 0 {
			effectiveCount = acc.totalRequests
		}

		modelItems = append(modelItems, ModelAnalyticsItem{
			Model:               acc.model,
			Provider:            acc.provider,
			RequestCount:        effectiveCount,
			Percentage:          pct,
			AvgConfidence:       avgConf,
			ConfidenceCount:     acc.confidenceCount,
			AvgHallucination:    avgHallucination,
			NLICount:            acc.nliCount,
			AvgToxicity:         avgTox,
			ToxicityCount:       acc.toxicityCount,
			TotalCostMicrocents: acc.costMicrocents,
			TotalCostDollars:    float64(acc.costMicrocents) / 1000000.0,
			DailyCounts:         dailyCounts,
		})
	}

	// Sort modelItems by request count descending
	for i := 0; i < len(modelItems)-1; i++ {
		for j := i + 1; j < len(modelItems); j++ {
			if modelItems[i].RequestCount < modelItems[j].RequestCount {
				modelItems[i], modelItems[j] = modelItems[j], modelItems[i]
			}
		}
	}

	var dailyTotals []DailyModelCount
	for _, d := range dateKeys {
		dailyTotals = append(dailyTotals, DailyModelCount{
			Date:  d,
			Count: dailyTotalsMap[d],
		})
	}

	if modelItems == nil {
		modelItems = []ModelAnalyticsItem{}
	}
	if dailyTotals == nil {
		dailyTotals = []DailyModelCount{}
	}

	startDate := ""
	endDate := ""
	if len(dateKeys) > 0 {
		startDate = dateKeys[0]
		endDate = dateKeys[len(dateKeys)-1]
	}

	return &AnalyticsResponse{
		TotalRequests:   totalRequests,
		WeeklyRequests:  weeklyRequests,
		WeeklyStartDate: startDate,
		WeeklyEndDate:   endDate,
		Models:          modelItems,
		DailyTotals:     dailyTotals,
	}, nil
}
