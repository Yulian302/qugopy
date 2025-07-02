

class PriorityQueue:
    def __init__(self):
        self.data = []
        self.map = {}  # value : idx

    def _parent(self, i):
        return (i-1)//2

    def _left(self, i):
        return 2*i+1

    def _right(self, i):
        return 2*i+2

    def peek(self):
        if self.is_empty():
            raise IndexError("Heap is empty")
        return self.data[0]

    def is_empty(self):
        """
        >>> pq.push(3)
        >>> pq.is_empty()
        False
        >>> pq.pop()
        3
        >>> pq.is_empty()
        True
        """
        return len(self.data) == 0

    def push(self, value):
        """
        >>> pq.push(3)
        >>> pq.push(4)
        >>> pq.push(5)
        >>> pq.delete(4)
        >>> pq.pop()
        3
        >>> pq.pop()
        5
        """
        self.data.append(value)
        self._heapify_up(len(self.data)-1)

    def pop(self):
        if self.is_empty():
            raise IndexError("Heap is empty")

        root = self.data[0]
        last = self.data.pop()
        if self.data:
            self.data[0] = last
            self._heapify_down(0)
        return root

    def delete(self, value_to_remove):
        if self.is_empty():
            return

        for idx, val in enumerate(self.data):
            if val == value_to_remove:
                self.data[idx], self.data[-1] = self.data[-1], self.data[idx]
                self.data.pop()

                if idx >= len(self.data):
                    return

                parent = self._parent(idx)

                if parent >= 0 and self.data[parent] > self.data[idx]:
                    self._heapify_up(idx)
                else:
                    self._heapify_down(idx)
                return
        return

    def _heapify_up(self, index):
        while index > 0:
            parent = self._parent(index)
            if self.data[index] < self.data[parent]:
                self.data[index], self.data[parent] = self.data[parent], self.data[index]
                index = parent
            else:
                break

    def _heapify_down(self, index):
        size = len(self.data)

        smallest = index
        left = self._left(index)
        right = self._right(index)
        if left < size and self.data[left] < self.data[smallest]:
            smallest = left
        if right < size and self.data[right] < self.data[smallest]:
            smallest = right

        if smallest != index:
            self.data[index], self.data[smallest] = self.data[smallest], self.data[index]
            self._heapify_down(smallest)


if __name__ == '__main__':
    import doctest
    doctest.testmod(extraglobs={'pq': PriorityQueue()})
