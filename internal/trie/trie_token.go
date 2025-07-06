package trie

import "strings"

type TrieTokenNode struct {
	children    map[string]*TrieTokenNode
	isEndOfWord bool
}

type TrieToken struct {
	Root *TrieTokenNode
}

func NewTokenTrie() *TrieToken {
	return &TrieToken{
		Root: &TrieTokenNode{children: make(map[string]*TrieTokenNode)},
	}
}

func (t *TrieToken) Populate(tokensGroups [][]string) {
	if len(tokensGroups) == 0 {
		return
	}
	for _, tokens := range tokensGroups {
		t.Insert(tokens)
	}
}

func (t *TrieToken) Insert(tokens []string) {
	if len(tokens) == 0 {
		return
	}
	curr := t.Root
	for _, tok := range tokens {
		if curr.children[tok] == nil {
			curr.children[tok] = &TrieTokenNode{children: make(map[string]*TrieTokenNode)}
		}
		curr = curr.children[tok]
	}
	curr.isEndOfWord = true

}

func (t *TrieToken) Delete(tokens []string) bool {
	if t.Root == nil || len(tokens) == 0 {
		return false
	}
	return t.deleteHelper(t.Root, tokens, 0)
}

func (t *TrieToken) deleteHelper(node *TrieTokenNode, tokens []string, index int) bool {
	if index == len(tokens) {
		if !node.isEndOfWord {
			return false
		}
		node.isEndOfWord = false
		return len(node.children) == 0
	}

	token := tokens[index]
	child, exists := node.children[token]
	if !exists {
		return false
	}

	shouldDeleteChild := t.deleteHelper(child, tokens, index+1)
	if shouldDeleteChild {
		delete(node.children, token)
		return len(node.children) == 0
	}
	return false
}

func (t *TrieToken) GetAllWords() []string {
	tokens := make([]string, 0)
	if len(t.Root.children) == 0 {
		return []string{}
	}
	t.getAllWordsDfs(t.Root, []string{}, &tokens)
	return tokens
}

func (t *TrieToken) getAllWordsDfs(node *TrieTokenNode, path []string, tokens *[]string) {
	if node.isEndOfWord {
		*tokens = append(*tokens, strings.Join(path, " "))
	}

	for token, child := range node.children {
		t.getAllWordsDfs(child, append(path, token), tokens)
	}
}

func (t *TrieToken) SearchPrefix(tokens []string, includePrefix bool) []string {
	if len(tokens) == 0 {
		return []string{}
	}
	curr := t.Root
	for _, token := range tokens {
		if curr.children[token] == nil {
			return []string{}
		}
		curr = curr.children[token]
	}
	words := make([]string, 0)
	if includePrefix {
		t.searchPrefixDfs(curr, tokens, &words)
	} else {
		t.searchPrefixDfs(curr, []string{}, &words)
	}
	return words
}
func (t *TrieToken) searchPrefixDfs(node *TrieTokenNode, path []string, words *[]string) {
	if node == nil {
		return
	}

	if node.isEndOfWord {
		*words = append(*words, strings.Join(path, " "))
	}

	for token, child := range node.children {
		t.searchPrefixDfs(child, append(path, token), words)
	}
}

func (t *TrieToken) StartsWith(tokens []string) bool {
	if len(tokens) == 0 {
		return true
	}
	curr := t.Root
	for _, token := range tokens {
		if curr.children[token] == nil {
			return false
		}
		curr = curr.children[token]
	}
	return true
}

func (t *TrieToken) FuzzySearch(pattern []string) []string {
	if len(pattern) == 0 {
		return []string{}
	}
	tokens := make([]string, 0)
	t.fuzzySearchDfs(t.Root, pattern, []string{}, 0, &tokens)
	return tokens
}

func (t *TrieToken) fuzzySearchDfs(node *TrieTokenNode, pattern, path []string, index int, words *[]string) {
	if index == len(pattern) {
		if node.isEndOfWord {
			*words = append(*words, strings.Join(path, " "))
		}
		return
	}

	token := pattern[index]
	switch token {
	case "?":
		for childToken, child := range node.children {
			t.fuzzySearchDfs(child, pattern, append(path, childToken), index+1, words)
		}
	case "*":
		t.fuzzySearchDfs(node, pattern, path, index+1, words)
		for childToken, child := range node.children {
			t.fuzzySearchDfs(child, pattern, append(path, childToken), index, words)
		}
	default:
		if node.children[token] != nil {
			t.fuzzySearchDfs(node.children[token], pattern, append(path, token), index+1, words)
		}

	}
}

func GenerateTokenTrie(tokenGroups [][]string) *TrieToken {
	if len(tokenGroups) == 0 {
		return nil
	}
	trie := NewTokenTrie()
	trie.Populate(tokenGroups)
	return trie
}
