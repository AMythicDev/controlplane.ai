package main

import (
	"context"
	"log"
	"os"
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
