package trie

type TrieRuneNode struct {
	children    map[rune]*TrieRuneNode
	isEndOfWord bool
}

type TrieRune struct {
	Root *TrieRuneNode
}

func NewRuneTrie() *TrieRune {
	return &TrieRune{
		Root: &TrieRuneNode{children: make(map[rune]*TrieRuneNode)},
	}
}

func (t *TrieRune) Populate(words []string) {
	if len(words) == 0 {
		return
	}
	for _, word := range words {
		t.Insert(word)
	}
}

func (t *TrieRune) Insert(word string) {
	if len(word) == 0 {
		return
	}
	curr := t.Root
	for _, char := range word {
		if curr.children[char] == nil {
			curr.children[char] = &TrieRuneNode{children: make(map[rune]*TrieRuneNode)}
		}
		curr = curr.children[char]
	}
	curr.isEndOfWord = true

}

func (t *TrieRune) Delete(word string) bool {
	if t.Root == nil || len(word) == 0 {
		return false
	}
	return t.deleteHelper(t.Root, word, 0)
}

func (t *TrieRune) deleteHelper(node *TrieRuneNode, word string, index int) bool {
	if index == len(word) {
		if !node.isEndOfWord {
			return false
		}
		node.isEndOfWord = false
		return len(node.children) == 0
	}

	char := rune(word[index])
	child, exists := node.children[char]
	if !exists {
		return false
	}

	shouldDeleteChild := t.deleteHelper(child, word, index+1)
	if shouldDeleteChild {
		delete(node.children, char)
		return len(node.children) == 0
	}
	return false
}

func (t *TrieRune) GetAllWords() []string {
	words := make([]string, 0)
	if len(t.Root.children) == 0 {
		return []string{}
	}
	t.getAllWordsDfs(t.Root, "", &words)
	return words
}

func (t *TrieRune) getAllWordsDfs(node *TrieRuneNode, path string, tokens *[]string) {
	if node.isEndOfWord {
		*tokens = append(*tokens, path)
	}

	for char, child := range node.children {
		t.getAllWordsDfs(child, path+string(char), tokens)
	}
}

func (t *TrieRune) SearchPrefix(prefix string, includePrefix bool) []string {
	if len(prefix) == 0 {
		return []string{}
	}
	curr := t.Root
	for _, char := range prefix {
		if curr.children[char] == nil {
			return []string{}
		}
		curr = curr.children[char]
	}
	words := make([]string, 0)
	if includePrefix {
		t.searchPrefixDfs(curr, prefix, &words)
	} else {
		t.searchPrefixDfs(curr, "", &words)
	}
	return words
}
func (t *TrieRune) searchPrefixDfs(node *TrieRuneNode, path string, words *[]string) {
	if node == nil {
		return
	}

	if node.isEndOfWord {
		*words = append(*words, path)
	}

	for char, child := range node.children {
		t.searchPrefixDfs(child, path+string(char), words)
	}
}

func (t *TrieRune) StartsWith(prefix string) bool {
	if len(prefix) == 0 {
		return true
	}
	curr := t.Root
	for _, char := range prefix {
		if curr.children[char] == nil {
			return false
		}
		curr = curr.children[char]
	}
	return true
}

func (t *TrieRune) FuzzySearch(pattern string) []string {
	if len(pattern) == 0 {
		return []string{}
	}
	words := make([]string, 0)
	t.fuzzySearchDfs(t.Root, pattern, "", 0, &words)
	return words
}

func (t *TrieRune) fuzzySearchDfs(node *TrieRuneNode, pattern string, path string, index int, words *[]string) {
	if index == len(pattern) {
		if node.isEndOfWord {
			*words = append(*words, path)
		}
		return
	}

	char := rune(pattern[index])
	switch char {
	case '?':
		for childChar, child := range node.children {
			t.fuzzySearchDfs(child, pattern, path+string(childChar), index+1, words)
		}
	case '*':
		t.fuzzySearchDfs(node, pattern, path, index+1, words)
		for childChar, child := range node.children {
			t.fuzzySearchDfs(child, pattern, path+string(childChar), index, words)
		}
	default:
		if node.children[char] != nil {
			t.fuzzySearchDfs(node.children[char], pattern, path+string(char), index+1, words)
		}
	}
}

func GenerateRuneTrie(words []string) *TrieRune {
	if len(words) == 0 {
		return nil
	}
	trie := NewRuneTrie()
	trie.Populate(words)
	return trie
}
