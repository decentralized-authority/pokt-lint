package maths

import "sort"

// Mean returns the average value of a set of numbers
func Mean(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}

	var tot float64
	for _, n := range nums {
		tot += n
	}
	return tot / float64(len(nums))
}

// Median sorts and returns the middle value, positionally in an odd set of numbers,
// or the Mean of the middle 2 numbers in an even set of numbers
func Median(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}

	sort.Slice(nums, func(i int, j int) bool { return nums[i] < nums[j] })
	if len(nums)%2 == 0 {
		half := len(nums) / 2
		return (nums[half-1] + nums[half]) / 2
	}

	pos := (len(nums) - 1) / 2
	return nums[pos]
}

// Min returns the smallest number in a set
func Min(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}

	min := nums[0]
	for i, n := range nums {
		if n < min {
			min = nums[i]
		}
	}
	return min
}

// Max returns the largest number in a set
func Max(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}

	max := nums[0]
	for i, n := range nums {
		if n > max {
			max = nums[i]
		}
	}
	return max
}
