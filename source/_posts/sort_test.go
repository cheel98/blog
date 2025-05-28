package _posts

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func TestName(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	//Quicksort(nums, 0, len(nums)-1)
	Quicksort2(nums)
	fmt.Println(nums)
}
func Quicksort(nums []int, start, end int) []int {
	n := end - start + 1
	if n < 2 {
		return nums
	}
	// 选择基准点
	pivot := start + rand.IntN(n)
	pivotValue := nums[pivot]
	nums[pivot], nums[start] = nums[start], pivotValue
	// l为小于pivot的边界:有l个小于pivot的数, i为数组索引
	l := start
	for i := start + 1; i <= end; i++ {
		if nums[i] < pivotValue {
			// i的值小于pivot，左边界移动,并准备交换
			l++
			if l != i {
				// 如果i不等于l，说明i和l之间的值是大于pivot,
				// nums[i] 与 num[l+1] 交换
				nums[l], nums[i] = nums[i], nums[l]
			}
		}
	}
	// 处理pivot
	nums[start], nums[l] = nums[l], nums[start]
	Quicksort(nums, start, l-1)
	Quicksort(nums, l+1, end)
	return nums
}

func Quicksort2(nums []int) {
	n := len(nums)
	if n < 2 {
		return
	}
	pivot := 0
	nums[pivot], nums[0] = nums[0], nums[pivot]
	l, r := 0, n-1
	for l < r {
		for r > l && nums[r] >= nums[l] {
			r--
		}
		for l < r && nums[l] <= nums[0] {
			l++
		}
		nums[l], nums[r] = nums[r], nums[l]
	}
	nums[0], nums[l] = nums[l], nums[0]
	Quicksort2(nums[:l])
	Quicksort2(nums[l+1:])
}

func TestMergeSort(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(MergeSort(nums))
}

func MergeSort(nums []int) []int {
	n := len(nums)
	if n < 2 {
		return nums
	}
	pivot := n / 2
	left := MergeSort(nums[:pivot])
	right := MergeSort(nums[pivot:])
	ans := Merge(left, right)
	return ans
}

func Merge(left, right []int) []int {
	n1, n2 := len(left), len(right)
	ans := make([]int, n1+n2)
	i, l, r := 0, 0, 0

	for l < n1 && r < n2 {
		if left[l] <= right[r] {
			ans[i] = left[l]
			l++
		} else {
			ans[i] = right[r]
			r++
		}
		i++
	}
	for l < n1 {
		ans[i] = left[l]
		l++
		i++
	}
	for r < n2 {
		ans[i] = right[r]
		r++
		i++
	}
	return ans
}
