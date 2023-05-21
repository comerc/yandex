package article

func FindLongestSequence1(arr []int) int {
	maxLen := 0
	curLen := 0
	for _, v := range arr {
		if v == 1 {
			curLen++
			if curLen > maxLen {
				maxLen = curLen
			}
		} else {
			curLen = 0
		}
	}
	return maxLen
}

func FindLongestSequence2(arr []int) int {
	maxLen, curLen := 0, 0
	for _, val := range arr {
		curLen = (curLen + 1) * val
		if curLen > maxLen {
			maxLen = curLen
		}
	}
	return maxLen
}
