package trie

type TrieNode struct {
	children    map[rune]*TrieNode
	isEndOfWord bool
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{children: make(map[rune]*TrieNode)},
	}
}

func (t *Trie) Populate(words []string) {
	if len(words) == 0 {
		return
	}
	for _, word := range words {
		t.Insert(word)
	}
}

func (t *Trie) Insert(word string) {
	if len(word) == 0 {
		return
	}
	curr := t.root
	for _, char := range word {
		if curr.children[char] == nil {
			curr.children[char] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		curr = curr.children[char]
	}
	curr.isEndOfWord = true
}

func (t *Trie) Delete(word string) bool {
	if t.root == nil || len(word) == 0 {
		return false
	}
	return t.deleteHelper(t.root, word, 0)
}

func (t *Trie) deleteHelper(node *TrieNode, word string, index int) bool {
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

func (t *Trie) GetAllWords() []string {
	words := make([]string, 0)
	if len(t.root.children) == 0 {
		return []string{}
	}
	t.getAllWordsDfs(t.root, "", &words)
	return words
}

func (t *Trie) getAllWordsDfs(node *TrieNode, path string, words *[]string) {
	if node.isEndOfWord {
		*words = append(*words, path)
	}

	for char, child := range node.children {
		t.getAllWordsDfs(child, path+string(char), words)
	}
}

func (t *Trie) FuzzySearch(pattern string) []string {
	if len(pattern) == 0 {
		return []string{}
	}
	words := make([]string, 0)
	t.fuzzySearchDfs(t.root, pattern, "", 0, &words)
	return words
}
func (t *Trie) SearchPrefix(prefix string) []string {
	if len(prefix) == 0 {
		return []string{}
	}
	curr := t.root
	for _, char := range prefix {
		if curr.children[char] == nil {
			return []string{}
		}
		curr = curr.children[char]
	}
	words := make([]string, 0)
	t.searchPrefixDfs(curr, prefix, &words)
	return words
}
func (t *Trie) searchPrefixDfs(node *TrieNode, path string, words *[]string) {
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

func (t *Trie) StartsWith(prefix string) bool {
	if len(prefix) == 0 {
		return true
	}
	curr := t.root
	for _, char := range prefix {
		if curr.children[char] == nil {
			return false
		}
		curr = curr.children[char]
	}
	return true
}

func (t *Trie) fuzzySearchDfs(node *TrieNode, pattern string, path string, index int, words *[]string) {
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

func GenerateTrie(words []string) *Trie {
	trie := NewTrie()
	trie.Populate(words)
	return trie
}
