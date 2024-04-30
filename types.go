package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Request struct {
	StartDate string `json:"startDate" bson:"startDate"`
	EndDate   string `json:"endDate" bson:"endDate"`
	MinCount  int    `json:"minCount" bson:"minCount"`
	MaxCount  int    `json:"maxCount" bson:"maxCount"`
}

type Response struct {
	Count   int         `json:"count" bson:"count"`
	Msg     string      `json:"msg" bson:"msg"`
	Records []RecordDTO `json:"records" bson:"records"`
}

type RecordDTO struct {
	Key        string    `json:"key" bson:"key"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	TotalCount int       `json:"totalCount" bson:"totalCount"`
}

type Record struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Key       string             `json:"key" bson:"key"`
	Counts    []int              `json:"counts" bson:"counts"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

func ResponseFromRecord(record Record) *RecordDTO {
	return &RecordDTO{
		Key:        record.Key,
		CreatedAt:  record.CreatedAt,
		TotalCount: record.SumCounts(),
	}
}

func (r *Record) SumCounts() int {
	var sum int
	for _, count := range r.Counts {
		sum += count
	}
	return sum
}
