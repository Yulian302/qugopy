package trie

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrieTokenCreate(t *testing.T) {
	tokens := [][]string{
		{"send", "email"},
		{"send", "file"},
		{"set", "name"},
		{"set", "password"},
		{"run", "script"},
	}
	trie := GenerateTokenTrie(tokens)
	count := 0
	for _, word := range trie.GetAllWords() {

		for i := 0; i < len(tokens); i++ {
			if word == strings.Join(tokens[i], " ") {
				count++
			}
		}
	}
	assert.Equal(t, len(tokens), count)
}

func TestTrieTokenFuzzySearch(t *testing.T) {
	tokens := [][]string{
		{"send", "email"},
		{"send", "file"},
		{"set", "name"},
		{"set", "password"},
		{"run", "script"},
		{"run", "executable"},
		{"start", "command", "workers"},
		{"start", "command", "redis"},
	}
	trie := GenerateTokenTrie(tokens)
	tests := []struct {
		Pattern []string
		Result  [][]string
	}{
		{
			Pattern: []string{"send", "*"},
			Result:  [][]string{{"send", "email"}, {"send", "file"}},
		},
		{
			Pattern: []string{"set", "?"},
			Result:  [][]string{{"set", "name"}, {"set", "password"}},
		},
		{
			Pattern: []string{"run", "*"},
			Result:  [][]string{{"run", "script"}, {"run", "executable"}},
		},
		{
			Pattern: []string{"start", "?", "workers"},
			Result:  [][]string{{"start", "command", "workers"}},
		},
	}
	for _, tt := range tests {
		foundTokens := trie.FuzzySearch(tt.Pattern)
		for _, resTokens := range tt.Result {
			resJoined := strings.Join(resTokens, " ")
			assert.Contains(t, foundTokens, resJoined)
		}
	}
}

func TestTrieRuneCreate(t *testing.T) {
	words := []string{
		"send",
		"set",
		"run",
		"start",
		"restart",
		"script",
	}
	trie := GenerateRuneTrie(words)
	count := 0
	for _, word := range trie.GetAllWords() {
		for i := 0; i < len(words); i++ {
			if word == words[i] {
				count++
			}
		}
	}
	assert.Equal(t, len(words), count)
}

func TestTrieRuneFuzzySearch(t *testing.T) {
	words := []string{
		"send",
		"sent",
		"set",
		"run",
		"runner",
		"restart",
		"script",
	}
	trie := GenerateRuneTrie(words)
	tests := []struct {
		Pattern string
		Result  []string
	}{
		{
			Pattern: "s?n?",
			Result:  []string{"send", "sent"},
		},
		{
			Pattern: "r*n",
			Result:  []string{"run"},
		},
		{
			Pattern: "re*",
			Result:  []string{"restart"},
		},
		{
			Pattern: "*t",
			Result:  []string{"set", "script"},
		},
	}
	for _, tt := range tests {
		foundWords := trie.FuzzySearch(tt.Pattern)
		for _, expected := range tt.Result {
			assert.Contains(t, foundWords, expected)
		}
	}
}
