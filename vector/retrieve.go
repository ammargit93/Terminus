package vector

// calculates cosine similarity of the user query with all embeddings in the Store
// and returns the top 5.
import (
	"fmt"
	"sort"

	"gonum.org/v1/gonum/floats"
)

func calculateCosineSimilarity(query []float32, embedding []float32) float64 {
	if len(query) != len(embedding) || len(query) == 0 {
		return 0
	}

	// Convert []float32 → []float64 for gonum
	q := make([]float64, len(query))
	e := make([]float64, len(embedding))
	for i := range query {
		q[i] = float64(query[i])
		e[i] = float64(embedding[i])
	}

	dot := floats.Dot(q, e)
	normA := floats.Norm(q, 2)
	normB := floats.Norm(e, 2)

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (normA * normB)
}

func GetTopResults(query string) {
	type scorePair struct {
		filepath string
		score    float64
	}

	var scoresList []scorePair
	embeddedQ := EmbedUserQuery(query)

	for _, pair := range Store {
		score := calculateCosineSimilarity(embeddedQ, pair.embedding)
		scoresList = append(scoresList, scorePair{
			filepath: pair.key,
			score:    score,
		})
	}

	// Sort descending by score
	sort.Slice(scoresList, func(i, j int) bool {
		return scoresList[i].score > scoresList[j].score
	})

	// Print top results
	for _, s := range scoresList {
		fmt.Printf("%s : %.6f\n", s.filepath, s.score)
	}
}
