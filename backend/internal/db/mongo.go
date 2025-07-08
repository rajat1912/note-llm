package db

import (
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func GetMongoDatabase() *mongo.Database {
	uri := viper.GetString("MONGODB_URI")

	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	db := client.Database("note-llm")

	return db
}
