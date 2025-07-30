---
title: 浅谈sync.Pool
author: cheel
abbrlink: 73dd0434
date: 2025-07-30 17:09:38
tags:
---
顾名思义,sync.Pool是一个线程安全的池，池子里放的是一些可以复用的变量，利用sync.Pool可以有效减少GC和内存分配次数。
<!--more-->
## 1 适用场景
对象创建/销毁多，但更新少

## 2 使用方法

### 2.1 声明对象池

只需要实现 New 函数即可。对象池中没有对象时，将会调用 New 函数创建。

``` golang 
var studentPool = sync.Pool{
	New: func() interface{} { 
   	 return new(Student)        
   },
}
```

### 2.2 Get & Put

```golang
stu := studentPool.Get().(\*Student)   
json.Unmarshal(buf, stu)   
studentPool.Put(stu) 
```
- `Get()` 用于从对象池中获取对象，因为返回值是 `interface{}`，因此需要类型转换。
- `Put()` 则是在对象使用完毕后，返回对象池。

### 2.3 测试

```golang
func BenchmarkUnmarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stu := &Student{}
		json.Unmarshal(buf, stu)
	}
}

func BenchmarkUnmarshalWithPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stu := studentPool.Get().(*Student)
		json.Unmarshal(buf, stu)
		studentPool.Put(stu)
	}
}
```

这个例子中只是将sync.Pool中的对象读取，并没有将对象更改。

## 3 bytes.Buffer

```golang
var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

var data = make([]byte, 10000)

func BenchmarkBufferWithPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Write(data)
		buf.Reset()
		bufferPool.Put(buf)
	}
}

func BenchmarkBuffer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var buf bytes.Buffer
		buf.Write(data)
	}
}
```

从pool中取出buffer并使用后，不要忘记reset操作

## fmt

在golang的fmt源码中使用了pool。
