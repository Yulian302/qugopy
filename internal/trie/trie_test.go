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
	words := [][]string{
		{"send",
			"set",
			"run",
			"start",
			"restart",
			"script"},
	}
	trie := GenerateRuneTrie(words, []int{0})
	count := 0
	for _, word := range trie.GetAllWords(0) {
		for i := 0; i < len(words); i++ {
			if word == words[0][i] {
				count++
			}
		}
	}
	assert.Equal(t, len(words), count)
}

func TestTrieRuneFuzzySearch_WithGroups(t *testing.T) {
	words := [][]string{
		{"send", "sent", "set"},          // Group 0
		{"run", "runner", "restart"},     // Group 1
		{"script", "service", "session"}, // Group 2
	}

	groups := []int{
		0,
		1,
		2,
	}

	trie := GenerateRuneTrie(words, groups)

	tests := []struct {
		Pattern string
		Group   int
		Result  []string
	}{
		{
			Pattern: "s?n?",
			Group:   0,
			Result:  []string{"send", "sent"},
		},
		{
			Pattern: "r*n",
			Group:   1,
			Result:  []string{"run"},
		},
		{
			Pattern: "re*",
			Group:   1,
			Result:  []string{"restart"},
		},
		{
			Pattern: "s*",
			Group:   2,
			Result:  []string{"script", "service", "session"},
		},
		{
			Pattern: "set",
			Group:   0,
			Result:  []string{"set"},
		},
		{
			Pattern: "*",
			Group:   1,
			Result:  []string{"run", "runner", "restart"},
		},
	}

	for _, tt := range tests {
		foundWords := trie.FuzzySearch(tt.Pattern, tt.Group)
		for _, expected := range tt.Result {
			assert.Contains(t, foundWords, expected, "pattern: %s group: %d", tt.Pattern, tt.Group)
		}
		assert.Len(t, foundWords, len(tt.Result), "unexpected count for pattern: %s group: %d", tt.Pattern, tt.Group)
	}
}

func TestTrieRunePrefixSearchWithGroups(t *testing.T) {
	groups := []struct {
		words []string
		group int
	}{
		{
			words: []string{"send", "set", "run", "start"}, group: 1,
		},
		{
			words: []string{"email", "file", "name", "password", "script"}, group: 2,
		},
		{
			words: []string{"hello", "world", "free", "server"}, group: 0, // no group
		},
	}
	trie := NewRuneTrie()
	for _, groupedWords := range groups {
		trie.Populate(groupedWords.words, groupedWords.group)
	}

	tests := []struct {
		prefix   string
		group    int
		expected []string
	}{
		{
			prefix: "s", group: 1,
			expected: []string{"send", "set", "start"},
		},
		{
			prefix: "f", group: 2,
			expected: []string{"file"},
		},
		{
			prefix: "s", group: 2,
			expected: []string{"script"},
		},
		{
			prefix: "w", group: 0,
			expected: []string{"world"},
		},
		{
			prefix: "s", group: 0,
			expected: []string{"server"},
		},
		{
			prefix: "x", group: 1,
			expected: []string{}, // no match
		},
	}

	for _, tt := range tests {
		results := trie.SearchPrefix(tt.prefix, true, tt.group)
		assert.ElementsMatch(t, tt.expected, results, "Prefix: %s, Group: %d", tt.prefix, tt.group)
	}

}
