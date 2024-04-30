package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	coll = "records"
)

type MongoStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

type InMemoryDb struct {
	payload *map[string]string
}

func newInMemoryDb() *InMemoryDb {
	return &InMemoryDb{
		payload: &map[string]string{},
	}
}

func (s *InMemoryDb) Insert(key string, value string) (*map[string]string, error) {
	(*s.payload)[key] = value
	return &map[string]string{
		"key":   key,
		"value": value,
	}, nil
}

func (s *InMemoryDb) Get(key string) (*map[string]string, error) {
	value := (*s.payload)[key]

	if value == "" {
		return nil, fmt.Errorf("key %s not found", key)
	}

	return &map[string]string{
		"key":   key,
		"value": value,
	}, nil
}

func newMongoStore(client *mongo.Client, dbName string) *MongoStore {
	return &MongoStore{
		client: client,
		coll:   client.Database(dbName).Collection(coll),
	}
}

func (s *MongoStore) InsertRecord(record Record) (*Record, error) {
	_, err := s.coll.InsertOne(context.TODO(), record)

	return &record, err
}

func (s *MongoStore) GetRecords(start time.Time, end time.Time, min int, max int) (*[]Record, error) {
	filter := bson.M{
		"createdAt": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}
	cur, err := s.coll.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var records []Record
	if err := cur.All(context.Background(), &records); err != nil {
		return nil, err
	}

	var response []Record

	for _, record := range records {
		if record.SumCounts() >= min && record.SumCounts() <= max {
			response = append(response, record)
		}
	}
	return &response, nil
}
