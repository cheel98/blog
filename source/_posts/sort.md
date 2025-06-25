---
title: 常见的排序算法golang实现
description: 整理一下常用的排序算法
tags:
  - 算法
  - 排序
  - 数组
categories:
  - 算法
date: '2025年5月28日09:07:33'
abbrlink: a5f1835d
---
<!-- more -->

## 快读排序

特点: 1. 原地排序 2. 不稳定

步骤: 

step1: 找到一个基准数，将大于它的数放在右边, 小于它的数放在左边

step2: 对左右两个数组执行step1, 指导数组元素小于等于1

```golang
func QuickSort(nums []int) {
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
    QuickSort(nums[:l])
    QuickSort(nums[l+1:])
}
```

## 归并排序

特点: 1. 需要辅助数组 2. 稳定排序

步骤：

1. 分解：
   1. 找到数组中心位置，从中心位置分解为左右两个数组
   2. 对左右两个数组递归分解，直到分解为n个长度为1的数组

2. 合并:
   1. 从长度为1的数组开始依次合并，合并时比较左右元素的大小使合并后的数组保持有序

```golang
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
```