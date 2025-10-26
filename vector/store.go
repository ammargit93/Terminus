package vector

import "fmt"

type Pair struct {
	key       string    // filenames
	embedding []float32 //  embeddings of files
}

var Store []Pair

// Adds a filename and embedding pair
// Adds it to the Store []Pair Slice
// file: embedding
func AddPair(filename string, embedding []float32) {
	Store = append(Store, Pair{key: filename, embedding: embedding})
}

// Iterates through the Store and deletes the key:embedding pair
// whose key matches filename
func RemovePair(filename string) {
	for i, pair := range Store {
		if pair.key == filename {
			Store = append(Store[:i], Store[i+1:]...)
		}
	}
}

func DisplayStore() {
	for _, pair := range Store {
		fmt.Printf("key: %s\nembedding: %v\n", pair.key, pair.embedding)
	}
}

func GetPair(key string) Pair {
	for _, pair := range Store {
		if pair.key == key {
			return pair
		}
	}
	return Pair{}
}
