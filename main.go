package main

import (
	"fmt"
	"time"

	"github.com/comerc/yandex/article"
)

func main() {
	var start time.Time
	var elapsed time.Duration
	var count int

	J := "abc"
	S := "aabbccd"

	start = time.Now()
	count = article.GetIntersectionCount1(J, S)
	elapsed = time.Since(start)
	fmt.Printf("article.GetIntersectionCount1 O(n+m): %d, time: %s\n", count, elapsed)

	start = time.Now()
	count = article.GetIntersectionCount2(J, S)
	elapsed = time.Since(start)
	fmt.Printf("article.GetIntersectionCount2 O(n*m): %d, time: %s\n", count, elapsed)

	start = time.Now()
	count = article.GetIntersectionCount3(J, S)
	elapsed = time.Since(start)
	fmt.Printf("article.GetIntersectionCount3 O(n^2): %d, time: %s\n", count, elapsed)

	start = time.Now()
	count = article.GetIntersectionCount4(J, S)
	elapsed = time.Since(start)
	fmt.Printf("article.GetIntersectionCount4 O(m*log(n)): %d, time: %s\n", count, elapsed)

	arr := []int{1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 1}

	start = time.Now()
	count = article.FindLongestSequence1(arr)
	elapsed = time.Since(start)
	fmt.Printf("article.FindLongestSequence1: %d, time: %s\n", count, elapsed)

	start = time.Now()
	count = article.FindLongestSequence2(arr)
	elapsed = time.Since(start)
	fmt.Printf("article.FindLongestSequence2: %d, time: %s\n", count, elapsed)

	arr32 := []int32{1, 2, 2, 3, 4, 4, 4, 5, 5, 6}
	fmt.Printf("article.RemoveDuplicates1: %v\n", article.RemoveDuplicates1(arr32))
	fmt.Printf("article.RemoveDuplicates2: %v\n", article.RemoveDuplicates2(arr32))

	fmt.Printf("article.GenerateParenthesis: %v\n", article.GenerateParenthesis(3))

	fmt.Printf("article.CompareStrings1: %v\n", article.CompareStrings1("hello", "ollhe"))
	fmt.Printf("article.CompareStrings2: %v\n", article.CompareStrings2("hello", "ollhe"))

	arrays := [][]int{{3, 5, 7}, {0, 6}, {0, 6, 28}}
	fmt.Printf("article.MergeArrays1: %v\n", article.MergeArrays1(arrays)) // [0 0 3 5 6 6 7 28]
	fmt.Printf("article.MergeArrays2: %v\n", article.MergeArrays2(arrays)) // [0 0 3 5 6 6 7 28]
	fmt.Printf("article.MergeArrays3: %v\n", article.MergeArrays3(arrays)) // [0 0 3 5 6 6 7 28]
}
