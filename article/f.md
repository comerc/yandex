```go
package article

import (
	"container/heap"
)

type Item struct {
	value    int // значение элемента
	arrayNum int // номер массива, из которого был взят элемент
	index    int // индекс элемента в массиве
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].value < pq[j].value
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func MergeArrays1(arrays [][]int) []int {
	k := len(arrays)
	pointers := make([]*int, k)
	for i := 0; i < k; i++ {
		pointers[i] = new(int)
	}
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	for i := 0; i < k; i++ {
		heap.Push(&pq, &Item{value: arrays[i][0], arrayNum: i, index: 0})
	}
	result := make([]int, 0)
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		result = append(result, item.value)
		if *pointers[item.arrayNum] < len(arrays[item.arrayNum])-1 {
			*pointers[item.arrayNum]++
			heap.Push(&pq, &Item{value: arrays[item.arrayNum][*pointers[item.arrayNum]], arrayNum: item.arrayNum, index: *pointers[item.arrayNum]})
		}
	}
	return result
}

func MergeArrays2(arrays [][]int) []int {
	k := len(arrays)
	pointers := make([]int, k)
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	for i := 0; i < k; i++ {
		heap.Push(&pq, &Item{value: arrays[i][0], arrayNum: i, index: 0})
	}
	result := make([]int, 0)
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		result = append(result, item.value)
		if pointers[item.arrayNum] < len(arrays[item.arrayNum])-1 {
			pointers[item.arrayNum]++
			heap.Push(&pq, &Item{value: arrays[item.arrayNum][pointers[item.arrayNum]], arrayNum: item.arrayNum, index: pointers[item.arrayNum]})
		}
	}
	return result
}

func MergeArrays3(arrays [][]int) []int {
	// Создание указателей на начало каждого массива
	pointers := make([]int, len(arrays))
	for i := range pointers {
		pointers[i] = -1
	}

	// Создание кучи и добавление первых элементов из каждого массива
	heap := make([]int, 0)
	for i, array := range arrays {
		if len(array) > 0 {
			pointers[i] = 0
			heap = append(heap, array[0])
		}
	}
	buildHeap(heap)

	// Слияние массивов
	result := make([]int, 0)
	for len(heap) > 0 {
		// Извлечение минимального элемента из кучи
		min := heap[0]
		result = append(result, min)
		heap[0] = heap[len(heap)-1]
		heap = heap[:len(heap)-1]
		heapify(heap, 0)

		// Добавление следующего элемента из соответствующего массива
		for i, array := range arrays {
			if pointers[i] >= 0 && pointers[i] < len(array) && array[pointers[i]] == min {
				pointers[i]++
				if pointers[i] < len(array) {
					heap = append(heap, array[pointers[i]])
					heapifyUp(heap, len(heap)-1)
				}
			}
		}
	}

	return result
}

func buildHeap(heap []int) {
	for i := len(heap) / 2; i >= 0; i-- {
		heapify(heap, i)
	}
}

func heapify(heap []int, i int) {
	left := 2*i + 1
	right := 2*i + 2
	smallest := i
	if left < len(heap) && heap[left] < heap[smallest] {
		smallest = left
	}
	if right < len(heap) && heap[right] < heap[smallest] {
		smallest = right
	}
	if smallest != i {
		heap[i], heap[smallest] = heap[smallest], heap[i]
		heapify(heap, smallest)
	}
}

func heapifyUp(heap []int, i int) {
	parent := (i - 1) / 2
	if parent >= 0 && heap[i] < heap[parent] {
		heap[i], heap[parent] = heap[parent], heap[i]
		heapifyUp(heap, parent)
	}
}
```