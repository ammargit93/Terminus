package vector

import (
	"context"
	"encoding/json"
	"log"
	"os"

	cohere "github.com/cohere-ai/cohere-go/v2"
	client "github.com/cohere-ai/cohere-go/v2/client"
)

type Embedding struct {
	embeddingModel string
	apiKey         string
	provider       string
}

func InitialiseEmbeddingModel() Embedding {
	return Embedding{
		embeddingModel: "embed-english-v3.0",
		apiKey:         os.Getenv("COHERE_API_KEY"),
		provider:       "cohere",
	}
}

func readFiles(files string) []string {
	var contentSlice []string
	for _, file := range files {
		contentByte, _ := os.ReadFile(string(file))
		content := string(contentByte)
		contentSlice = append(contentSlice, content)
	}

	return contentSlice
}

func CallCohere(files []string) [][]float32 {
	co := client.NewClient(client.WithToken(os.Getenv("COHERE_API_KEY")))
	model := "embed-english-v3.0"
	inputType := cohere.EmbedInputTypeSearchDocument // SDK-provided enum/pointer

	resp, err := co.Embed(
		context.TODO(),
		&cohere.EmbedRequest{
			Texts:     files,
			Model:     &model,
			InputType: &inputType,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	r, _ := json.Marshal(resp)
	var v interface{}
	json.Unmarshal(r, &v)
	embedding := v.(map[string]interface{})
	return embedding["embeddings"].([][]float32)
}
