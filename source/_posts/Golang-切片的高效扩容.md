title: Golang 切片的扩容机制
author: cheel
abbrlink: f7b21975
tags: []
categories: []
date: 2025-07-29 11:30:00
---
介绍一下golang中切片的复制，扩容等机制。
<!-- more-->
## slice底层结构

```golang
type slice struct {
    ptr *T       // 指向底层数组的指针
    len int      // 当前元素数
    cap int      // 容量（数组大小）
}
```


### 扩容


v1.8 之前新容量的计算规则如下： 
- 需要的容量比2倍容量大：使用需要的容量 
  - 一般发生于append一个较大的slice时，例如 append(s, s1...)   

- 如果原容量小于1024，按照2倍扩容  
- 如果原容量大于等于1024，按照1.25倍扩容，对应源码newcap += newcap / 4  
- 如果newcap溢出了int最大值，不扩容

最新的扩容规则为：

- 需要的容量比2倍原容量大：使用需要的容量
- 如果原容量小于256， 2倍扩容
- 如果原容量大于等于256，newcap = 1.25倍的newcap + 192 计算新的容量直到大于所需容量
- 溢出则报panic
- 计算出新的cap之后，还需要进行内存对齐。


#### 源码 （v1.24.1）

```golang
/**
oldPtr	原切片底层数组的指针
newLen	扩容后的长度（=旧长度 + 新添加元素数量）
oldCap	原切片的容量
num	要添加的元素数量
et	元素类型 _type（用于确定大小、是否包含指针等）
*/
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
	...
	newcap := nextslicecap(newLen, oldCap)
	switch {
	case et.Size_ == 1:
		lenmem = uintptr(oldLen)
		newlenmem = uintptr(newLen)
		capmem = roundupsize(uintptr(newcap), noscan)
		overflow = uintptr(newcap) > maxAlloc
		newcap = int(capmem)
	case et.Size_ == goarch.PtrSize:
		...
	case isPowerOfTwo(et.Size_):
		...
	default:
		...
	}
	...
	return slice{p, newLen, newcap}
}

func nextslicecap(newLen, oldCap int) int {
	newcap := oldCap
	doublecap := newcap + newcap
	if newLen > doublecap {
		return newLen // 如果两倍的原cap 仍然不足以放下新的元素，则直接使用 旧长度 + 新添加元素数量
	}

	const threshold = 256
	if oldCap < threshold {
		return doublecap // 如果原cap 小于256 则扩容倍数为2
	}
	for {
        // 等价于 newcap = 1.25倍的newcap + 192
		newcap += (newcap + 3*threshold) >> 2

		if uint(newcap) >= uint(newLen) {
			break
		}
	}

	if newcap <= 0 {
		return newLen
	}
	return newcap
}

// 内存对齐
func roundupsize(size uintptr, noscan bool) (reqSize uintptr) {
	reqSize = size
	if reqSize <= maxSmallSize-mallocHeaderSize {
		if !noscan && reqSize > minSizeForMallocHeader { 
			reqSize += mallocHeaderSize
		}
		if reqSize <= smallSizeMax-8 {
			return uintptr(class_to_size[size_to_class8[divRoundUp(reqSize, smallSizeDiv)]]) - (reqSize - size)
		}
		return uintptr(class_to_size[size_to_class128[divRoundUp(reqSize-smallSizeMax, largeSizeDiv)]]) - (reqSize - size)
	}
	// Large object. Align reqSize up to the next page. Check for overflow.
	reqSize += pageSize - 1
	if reqSize < size {
		return size
	}
	return reqSize &^ (pageSize - 1)
}
```
### 注意事项

1. 进行大容量的合并时，避免使用golang自动扩容（append方法），而是使用预留cap和copy方法


### 内存对齐

#### 提高访问速率

如 
```golang 
type dog struct {
	age uint8
   balance uint64
}
```
定义一个新的dog对象后, age 占8位，balance 占64位，在64位的操作系统中,通过```fmt.Println(unsafe.Sizeof(Args{}))``` 输出结果位16, 但我们可以看到结构体其实只有(8+64)/8=9 位，其与7位即内存对齐的效果。

在没有内存对齐之前，64位操作系统一次可以读取8个字节,在读取balance时需要进行两次内存读取。
内存对齐之后，只需要进行一次内存读取。