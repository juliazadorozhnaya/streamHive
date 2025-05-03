package state

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	coll *mongo.Collection
}

type Job struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	JobID     string             `bson:"job_id"`
	StreamURL string             `bson:"stream_url"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"created_at"`
	Objects   []string           `bson:"objects,omitempty"`
	Duration  float64            `bson:"duration,omitempty"`
}

func NewMongo(uri, dbName, collection string) *MongoStore {
	ctx := context.Background()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	//  TTL-индекс по created_at
	_, _ = client.Database(dbName).Collection(collection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"created_at": 1},
		Options: options.Index().SetExpireAfterSeconds(259200), // 72 часа
	})

	return &MongoStore{
		coll: client.Database(dbName).Collection(collection),
	}
}

func (m *MongoStore) SaveJob(jobID, streamURL, status string) error {
	doc := Job{
		JobID:     jobID,
		StreamURL: streamURL,
		Status:    status,
		CreatedAt: time.Now(),
	}
	_, err := m.coll.InsertOne(context.Background(), doc)
	return err
}

func (m *MongoStore) UpdateJobStatus(jobID, status string) error {
	filter := bson.M{"job_id": jobID}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := m.coll.UpdateOne(context.Background(), filter, update)
	return err
}

func (m *MongoStore) GetJobStatus(jobID string) (string, error) {
	filter := bson.M{"job_id": jobID}
	var result Job
	err := m.coll.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Status, nil
}
