package syncsnippet

import (
	"encoding/json"
	"sync"
	"testing"
)

type Student struct {
	Name string
	Age  int
}

var buf, _ = json.Marshal(Student{
	Name: "qpz",
	Age:  30,
})
var stuPool = sync.Pool{New: func() interface{} { return new(Student) }}

func BenchmarkUnmarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stu := &Student{}
		err := json.Unmarshal(buf, stu)
		if err != nil {
			return
		}
	}
}

func BenchmarkUnmarshalWithPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stu := stuPool.Get().(*Student)
		_ = json.Unmarshal(buf, stu)
		stuPool.Put(stu) // 别忘了归还，否则 Pool 失去意义
	}
}

/*
注意点1：benchmark函数中必须要写for n := 0; n < b.N; n++
原因：
b.N 是 Go benchmark 框架根据目标精度自动调整的循环次数（可能几千到几百万次）。
只有把你真正要测的那段代码放进这个循环里，框架才能统计“每次操作的平均耗时/内存”。
如果你把代码写在循环外，框架只能反复调用整个函数来凑时长，测到的是“函数调用开销”，而不是你关心的那段逻辑。
所以，所有 benchmark 函数都必须包含这个固定循环，把被测代码放在循环体内

注意点2：运行结果
[root@jeremyqin-1ja635hgwa /data/workspace/github/golang_snippet/sync_snippet]# go test -bench=. -benchmem sync_pool_test.go
goos: linux
goarch: amd64
cpu: AMD EPYC 9754 128-Core Processor
BenchmarkUnmarshal-64                    2092827               624.0 ns/op           248 B/op          6 allocs/op
BenchmarkUnmarshalWithPool-64            2026566               581.6 ns/op           224 B/op          5 allocs/op
PASS
ok      command-line-arguments  3.674s

---------------------
BenchmarkUnmarshal-64 //-64代表GOMAXPROCS
2092827     //执行次数
624.0 ns/op //每次操作的耗时
248 B/op    //每次操作所占用的内存
6 allocs/op //每次操作的堆分配次数


注意点3：
-bench 和 -benchmem 都是 go test 在跑基准测试（benchmark）时的可选标志：
-bench=.
含义：告诉 go test 运行哪些基准测试函数。
语法：-bench=正则表达式，. 表示“匹配所有以 Benchmark 开头的函数”。
示例：
-bench=. → 运行全部 benchmark；
-bench=BenchmarkUnmarshal → 只跑 BenchmarkUnmarshal；
-bench=Fib|Sort → 跑名字里带 Fib 或 Sort 的 benchmark。
-benchmem
含义：在 benchmark 结果里额外打印内存分配统计。
输出列：
B/op —— 每次操作平均分配了多少字节；
allocs/op —— 每次操作平均发生了多少次堆分配。
如果不加 -benchmem，这两列不会出现。


注意点4：为什么2个benchmark的结果非常接近？是不是代表sync.pool没什么用？
不是 `sync.Pool` 没用，而是**这个场景里它带来的收益本来就很小**，所以看起来差距不大。原因主要有三点：

### 1. 对象太轻量
`Student` 只有两个字段，**一次分配只有 24 字节**（`string` 指向的底层字节数组是复用的，不计入本次分配）。
创建/销毁这种小对象的开销极低，Pool 能省的只是“一次 24 B 的内存分配 + 一次 GC 扫描”，收益自然不明显。

### 2. 额外开销抵消了收益
使用 Pool 时，每次都要执行：

```go
stu := stuPool.Get().(*Student) // 原子操作 + 类型断言
stuPool.Put(stu)                // 原子操作 + 写屏障
```

这些**同步开销**在对象很轻量的场景下，几乎把“省一次小对象分配”的收益吃掉了。

### 3. 测试数据量太小
benchmark 只反序列化一条固定 JSON，**没有并发压力**，也没有频繁创建/销毁对象。
Pool 真正发挥威力的地方通常是：

- 高并发（大量 goroutine 同时需要临时对象）；
- 对象本身较大或构造代价高（如 `bytes.Buffer`、`gzip.Writer`）；
- 生命周期极短（用完后立即丢弃）。

## 如何验证 Pool 的价值？

把测试改得“更重”一点就能看到差距：

1. 把 `Student` 放大（加很多字段或嵌套结构）；
2. 并发跑 benchmark（`go test -bench=. -cpu=1,8,32`）；
3. 把 JSON 变大（几十 KB 以上）；
4. 在业务代码里用 `pprof` 看 GC 压力。

通常在这些场景下，Pool 能把 **B/op 和 allocs/op 降一个数量级**，延迟也会更稳定。

### 结论

- **当前场景**：对象太小、无并发、无 GC 压力 → Pool 收益≈0。
- **真实高并发/大对象/频繁创建场景**：Pool 通常能把内存分配和 GC 时间压到原来的 1/10 甚至更低。

*/
