package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_NAME")
	port := os.Getenv("LISTEN_ADDR")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	var (
		store      = newMongoStore(client, dbName)
		inMemoryDb = newInMemoryDb()
		handler    = newHandler(store, inMemoryDb)
	)

	http.HandleFunc("/records", handler.HandleGetRecords)
	http.HandleFunc("/payload", handler.HandlePostPayload)
	http.HandleFunc("/payload/", handler.HandleGetPayloadById)
	http.ListenAndServe(port, nil)
}
