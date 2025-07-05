package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrieCreate(t *testing.T) {
	words := []string{"hello", "world"}
	trie := GenerateTrie(words)
	count := 0
	for _, word := range trie.GetAllWords() {
		if word == "hello" || word == "world" {
			count++
		}
	}
	assert.Equal(t, count, len(words))
}

func TestTrieFuzzySearch(t *testing.T) {
	words := []string{"to", "tea", "ted", "ten"}
	trie := GenerateTrie(words)
	tests := []struct {
		Pattern string
		Result  []string
	}{
		{
			Pattern: "t?a",
			Result:  []string{"tea"},
		},
		{
			Pattern: "t*",
			Result:  []string{"to", "tea", "ted", "ten"},
		},
		{
			Pattern: "t?*",
			Result:  []string{"to", "tea", "ted", "ten"},
		},
		{
			Pattern: "t*a",
			Result:  []string{"tea"},
		},
	}
	for _, tt := range tests {
		foundWords := trie.FuzzySearch(tt.Pattern)
		for _, res := range tt.Result {
			assert.Contains(t, foundWords, res)
		}
	}
}
