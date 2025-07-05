

class TrieNode:
    def __init__(self, value=None):
        self.value = value
        self.children = {}
        self.is_end_of_word = False


class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word: str):
        if not word:
            return
        curr = self.root
        for ch in word:
            if ch not in curr.children:
                curr.children[ch] = TrieNode(ch)
            curr = curr.children[ch]
        curr.is_end_of_word = True

    def delete(self, word: str):

        def _delete(node: TrieNode, word: str, index: int):
            if index == len(word):
                if not node.is_end_of_word:
                    return False
                node.is_end_of_word = False
                return len(node.children) == 0

            ch = word[index]
            if ch not in node.children:
                return False

            should_delete_child = _delete(node.children[ch], word, index+1)
            if should_delete_child:
                del node.children[ch]
                return len(node.children) == 0
            return False

        return _delete(self.root, word, 0)

    def getAllWords(self):
        if not self.root.children:
            return []
        words = []

        def dfs(node: TrieNode, path: str):
            if node.is_end_of_word:
                words.append(path)

            for ch, child in node.children.items():
                dfs(child, path+ch)

        dfs(self.root, "")
        return words

    def getLongestCommonPrefix(self) -> str:
        prefix = []
        curr = self.root
        while len(curr.children) == 1 and not curr.is_end_of_word:
            ch = next(iter(curr.children))
            prefix.append(ch)
            curr = curr.children[ch]
        return "".join(prefix)

    def fuzzySearch(self, pattern: str):
        if not pattern:
            return []

        words = []

        def dfs(node: TrieNode, path: str, index: int):
            if index == len(pattern):
                if node.is_end_of_word:
                    words.append(path)
                return

            ch = pattern[index]
            if ch == "?":
                for ch, child in node.children.items():
                    dfs(child, path+ch, index+1)
            elif ch == "*":
                dfs(node, path, index + 1)
                for child_ch, child in node.children.items():
                    dfs(child, path + child_ch, index)
            else:
                if ch in node.children:
                    dfs(node.children[ch], path + ch, index + 1)

        dfs(self.root, "", 0)
        return words

    def populate(self, words: list[str]):
        if not words:
            return
        for word in words:
            self.insert(word)

    def startsWith(self, prefix: str) -> bool:
        if not prefix:
            return True

        curr = self.root
        for ch in prefix:
            if ch not in curr.children:
                return False
            curr = curr.children[ch]
        return True

    def searchPrefix(self, prefix: str) -> list[str]:
        if not prefix:
            return []
        curr = self.root
        for ch in prefix:
            if ch not in curr.children:
                return []
            curr = curr.children[ch]

        words = []

        def dfs(node: TrieNode, path: str):
            if not node:
                return

            if node.is_end_of_word:
                words.append(path)

            for ch, child in node.children.items():
                dfs(child, path+ch)

        dfs(curr, prefix)
        return words
