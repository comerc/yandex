package article

func RemoveDuplicates1(arr []int32) []int32 {
	lastValue := arr[0]
	result := []int32{lastValue}
	for i := 1; i < len(arr); i++ {
		if arr[i] == lastValue {
			continue
		} else {
			result = append(result, arr[i])
			lastValue = arr[i]
		}
	}
	return result
}

// При решении этой задачи также не нужно использовать дополнительную память.

func RemoveDuplicates2(arr []int32) []int32 {
	uniqueIndex := 1
	for i := 1; i < len(arr); i++ {
		j := 0
		for j < uniqueIndex {
			if arr[i] == arr[j] {
				break
			}
			j++
		}
		if j == uniqueIndex {
			arr[uniqueIndex] = arr[i]
			uniqueIndex++
		}
	}
	return arr[:uniqueIndex]
}
