package qdrant

import (
	"github.com/qdrant/go-client/qdrant"
	"github.com/spf13/viper"
)

func GetQdrantClient() *qdrant.Client {
	apiKey := viper.GetString("QDRANT_API")
	host := viper.GetString("QDRANT_HOST")
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   host,
		Port:   6334,
		APIKey: apiKey,
		UseTLS: true,
	})
	if err != nil {
		panic(err)
	}
	return client
}
