package vector

type Pair struct {
	key       string      // filenames
	embedding [][]float32 //  embeddings of files
}

var Store []Pair
