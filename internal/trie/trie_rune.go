package trie

// group - is a group of a word (default: 0)
// it is required to obtain only words that are members of a specific group
// 0 as an argument means no group (returns all)
type TrieRuneNode struct {
	children    map[rune]*TrieRuneNode
	isEndOfWord bool
	group       int
}

type TrieRune struct {
	Root *TrieRuneNode
}

func NewRuneTrie() *TrieRune {
	return &TrieRune{
		Root: &TrieRuneNode{children: make(map[rune]*TrieRuneNode)},
	}
}

func (t *TrieRune) Populate(words []string, group int) {
	if len(words) == 0 {
		return
	}
	for _, word := range words {
		t.Insert(word, group)
	}
}

func (t *TrieRune) Insert(word string, group int) {
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
	curr.group = group
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
		node.group = 0
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

func (t *TrieRune) GetAllWords(group int) []string {
	words := make([]string, 0)
	if len(t.Root.children) == 0 {
		return []string{}
	}
	t.getAllWordsDfs(t.Root, "", &words, group)
	return words
}

func (t *TrieRune) getAllWordsDfs(node *TrieRuneNode, path string, tokens *[]string, group int) {
	if node.isEndOfWord {
		if node.group == group || group == 0 && node.group == 0 {
			*tokens = append(*tokens, path)
		}
	}

	for char, child := range node.children {
		t.getAllWordsDfs(child, path+string(char), tokens, group)
	}
}

func (t *TrieRune) SearchPrefix(prefix string, includePrefix bool, group int) []string {
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
		t.searchPrefixDfs(curr, prefix, &words, group)
	} else {
		t.searchPrefixDfs(curr, "", &words, group)
	}
	return words
}
func (t *TrieRune) searchPrefixDfs(node *TrieRuneNode, path string, words *[]string, group int) {
	if node == nil {
		return
	}

	if node.isEndOfWord {
		if node.group == group || group == 0 && node.group == 0 {
			*words = append(*words, path)
		}
	}

	for char, child := range node.children {
		t.searchPrefixDfs(child, path+string(char), words, group)
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

func (t *TrieRune) FuzzySearch(pattern string, group int) []string {
	if len(pattern) == 0 {
		return []string{}
	}
	words := make([]string, 0)
	t.fuzzySearchDfs(t.Root, pattern, "", 0, &words, group)
	return words
}

func (t *TrieRune) fuzzySearchDfs(node *TrieRuneNode, pattern string, path string, index int, words *[]string, group int) {
	if index == len(pattern) {
		if node.isEndOfWord {
			if node.group == group || group == 0 && node.group == 0 {
				*words = append(*words, path)
			}
		}
		return
	}

	char := rune(pattern[index])
	switch char {
	case '?':
		for childChar, child := range node.children {
			t.fuzzySearchDfs(child, pattern, path+string(childChar), index+1, words, group)
		}
	case '*':
		t.fuzzySearchDfs(node, pattern, path, index+1, words, group)
		for childChar, child := range node.children {
			t.fuzzySearchDfs(child, pattern, path+string(childChar), index, words, group)
		}
	default:
		if node.children[char] != nil {
			t.fuzzySearchDfs(node.children[char], pattern, path+string(char), index+1, words, group)
		}
	}
}

func GenerateRuneTrie(words [][]string, groupNumbers []int) *TrieRune {
	if len(words) == 0 || len(words) != len(groupNumbers) {
		return nil
	}
	trie := NewRuneTrie()
	for idx, groupedWords := range words {
		trie.Populate(groupedWords, groupNumbers[idx])
	}
	return trie
}
