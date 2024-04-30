package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	uri    = ""
	dbName = ""
)

func TestGetRecords(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}
	mongo := newMongoStore(client, dbName)
	memory := newInMemoryDb()
	handler := newHandler(mongo, memory)

	reqBody := Request{
		StartDate: "2015-01-01",
		EndDate:   "2020-01-02",
		MinCount:  0,
		MaxCount:  6000,
	}

	requestBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/records", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	handler.HandleGetRecords(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetRecordsNotFound(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}
	mongo := newMongoStore(client, dbName)
	memory := newInMemoryDb()
	handler := newHandler(mongo, memory)

	reqBody := Request{
		StartDate: "2010-01-01",
		EndDate:   "2010-01-02",
		MinCount:  0,
		MaxCount:  0,
	}

	requestBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/records", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	handler.HandleGetRecords(rr, req)
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestPostPayload(t *testing.T) {
	InMemoryDb := newInMemoryDb()
	handler := newHandler(nil, InMemoryDb)

	type Req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	reqBody := Req{
		Key:   "key",
		Value: "value",
	}

	requestBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/payload", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	handler.HandlePostPayload(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestPostPayloadFail(t *testing.T) {
	InMemoryDb := newInMemoryDb()
	handler := newHandler(nil, InMemoryDb)

	reqBody := "asd"
	requestBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/payload", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	handler.HandlePostPayload(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

func TestGetPayloadById(t *testing.T) {
	InMemoryDb := newInMemoryDb()
	InMemoryDb.Insert("key", "value")
	handler := newHandler(nil, InMemoryDb)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/payload/?key=key", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.HandleGetPayloadById(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetPayloadByIdFail(t *testing.T) {
	InMemoryDb := newInMemoryDb()
	handler := newHandler(nil, InMemoryDb)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/payload/?key=key", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.HandleGetPayloadById(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
