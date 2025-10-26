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

// Initialises the Embedding model, redundant for now
func InitialiseEmbeddingModel() Embedding {
	return Embedding{
		embeddingModel: "embed-english-v3.0",
		apiKey:         os.Getenv("COHERE_API_KEY"),
		provider:       "cohere",
	}
}

// Reads each file from the  Absolute filepaths slice and returns the slice of the contents of each.
func ReadFiles(files []string) []string {
	var contentSlice []string
	for _, file := range files {
		contentByte, _ := os.ReadFile(string(file))
		content := string(contentByte)
		contentSlice = append(contentSlice, content)
	}

	return contentSlice
}

// takes in a slice of absolute filepaths and reads using readFiles, calls the Cohere API
// and embeds all the content of the files it also adds the filepaths and embeddings to the Store
func CallCohere(files []string) {
	contentSlice := ReadFiles(files)
	// contentSlice := files
	co := client.NewClient(client.WithToken(os.Getenv("COHERE_API_KEY")))
	model := "embed-english-v3.0"
	inputType := cohere.EmbedInputTypeSearchDocument // SDK-provided enum/pointer

	resp, err := co.Embed(
		context.TODO(),
		&cohere.EmbedRequest{
			Texts:     contentSlice,
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
	all := embedding["embeddings"]

	outer := all.([]interface{}) // top level = []interface{}
	var allEmbeddings [][]float32

	for _, row := range outer {
		inner := row.([]interface{}) // inner level = []interface{}
		vec := make([]float32, len(inner))
		for j, val := range inner {
			vec[j] = float32(val.(float64)) // convert float64 → float32
		}
		allEmbeddings = append(allEmbeddings, vec)
	}
	for i, embedding := range allEmbeddings {
		AddPair(files[i], embedding)
	}

}

func EmbedUserQuery(query string) []float32 {
	co := client.NewClient(client.WithToken(os.Getenv("COHERE_API_KEY")))
	model := "embed-english-v3.0"
	inputType := cohere.EmbedInputTypeSearchDocument // SDK enum

	resp, err := co.Embed(
		context.TODO(),
		&cohere.EmbedRequest{
			Texts:     []string{query},
			Model:     &model,
			InputType: &inputType,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Marshal + unmarshal to map for flexible parsing
	r, _ := json.Marshal(resp)
	var v map[string]interface{}
	json.Unmarshal(r, &v)

	// Extract the embeddings field
	rawEmbeddings, ok := v["embeddings"].([]interface{})
	if !ok || len(rawEmbeddings) == 0 {
		log.Fatal("no embeddings found in response")
	}

	// First embedding (since we passed one query)
	raw := rawEmbeddings[0].([]interface{})
	queryEmbed := make([]float32, len(raw))
	for i, val := range raw {
		queryEmbed[i] = float32(val.(float64))
	}

	return queryEmbed
}
