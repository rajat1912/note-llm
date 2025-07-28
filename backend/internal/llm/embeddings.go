package llm

import (
	"context"
	"fmt"
	"sync"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/spf13/viper"
)

var (
	client  openai.Client
	once    sync.Once
	initErr error
)

func InitOpenAIClient() {
	once.Do(func() {
		apiKey := viper.GetString("OPENAI_API")
		if apiKey == "" {
			initErr = fmt.Errorf("Missing OPENAI_API_KEY in config")
			return
		}
		client = openai.NewClient(option.WithAPIKey(apiKey))
	})
}

func GetEmbeddings(texts []string) ([][]float32, error) {
	InitOpenAIClient()
	if initErr != nil {
		return nil, initErr
	}

	// Call embeddings endpoint
	res, err := client.Embeddings.New(context.Background(), openai.EmbeddingNewParams{
		Model: openai.EmbeddingModelTextEmbedding3Small, // Fixed model constant
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: texts,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}

	var result [][]float32
	for _, item := range res.Data {
		// Convert []float64 to []float32
		embedding := make([]float32, len(item.Embedding))
		for i, val := range item.Embedding {
			embedding[i] = float32(val)
		}
		result = append(result, embedding)
	}

	return result, nil
}
