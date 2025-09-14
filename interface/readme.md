## Go 接口（Interface）知识库

本文档整理 Go 接口的核心概念、实现机制、类型断言等重要知识点。

### 基本概念

#### 接口定义
```go
type I interface {
    Get() int
    Set(int)
}
```

- 接口是方法签名的集合
- 接口定义了行为规范，不包含实现
- 类型通过实现接口的所有方法来"隐式"实现接口

#### 接口实现规则
```go
type S struct {
    Age int
}

// 值接收者方法
func (s S) Get() int {
    return s.Age
}

// 指针接收者方法  
func (s *S) Set(age int) {
    s.Age = age
}
```

**关键规则：**
- **指针类型 `*S`** 实现了接口 I（可以调用值接收者和指针接收者的方法）
- **值类型 `S`** 没有完全实现接口 I（无法调用指针接收者的方法）

### 方法接收者与接口实现

#### 接收者类型的影响

| 方法接收者 | 值类型调用 | 指针类型调用 | 接口实现能力 |
|----------|----------|------------|------------|
| `func (s S) Method()` | ✅ 直接调用 | ✅ 自动解引用 | 值和指针都可实现接口 |
| `func (s *S) Method()` | ✅ 自动取地址 | ✅ 直接调用 | 只有指针可实现接口 |

#### 自动转换机制

**Go 的智能转换：**
```go
s := S{1}
ptr := &s

// 指针调用值接收者方法 - 自动解引用
name := ptr.Get() // 等价于 (*ptr).Get()

// 值调用指针接收者方法 - 自动取地址  
s.Set(2)         // 等价于 (&s).Set(2)
```

**为什么要自动转换？**
1. **简化语法**：统一的方法调用体验
2. **减少重复**：避免为值和指针分别定义方法
3. **提升开发效率**：开发者无需过分关心值/指针区别

**转换限制：**
- 方法调用：支持自动转换
- 函数参数：不支持自动转换，必须显式转换
- 接口赋值：遵循接收者规则

### 类型断言与类型判断

#### 类型断言（Type Assertion）
```go
s := S{1}
var i I = &s

// 方式1：安全断言
if ptr, ok := i.(*S); ok {
    t.Log("i 的类型是 *S")
    // 使用 ptr...
}

// 方式2：直接断言（可能 panic）
ptr := i.(*S)  // 如果类型不匹配会 panic
```

#### 类型开关（Type Switch）
```go
switch v := i.(type) {
case *S:
    t.Log("类型是 *S，值为:", v)
case *R:
    t.Log("类型是 *R，值为:", v)
case nil:
    t.Log("接口为 nil")
default:
    t.Log("未知类型:", v)
}
```

### 接口的零值与 nil

```go
var i I           // i == nil，类型和值都为 nil
var s *S          // s == nil
var i2 I = s      // i2 != nil！类型为 *S，值为 nil

// 检查接口是否真正为空
if i == nil {
    // 接口本身为 nil
}

// 检查接口包装的值是否为 nil
if i != nil && reflect.ValueOf(i).IsNil() {
    // 接口不为 nil，但包装的值为 nil
}
```

### 空接口与类型断言

```go
// 空接口可以存储任何类型
var any interface{} = "hello"

// 类型断言获取具体类型
if str, ok := any.(string); ok {
    fmt.Println("字符串:", str)
}

// 类型开关处理多种类型
switch v := any.(type) {
case string:
    fmt.Println("字符串:", v)
case int:
    fmt.Println("整数:", v)
default:
    fmt.Printf("未知类型: %T\n", v)
}
```

### 接口组合

```go
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

// 接口组合
type ReadWriter interface {
    Reader  // 嵌入 Reader 接口
    Writer  // 嵌入 Writer 接口
}

// 等价于：
type ReadWriter interface {
    Read([]byte) (int, error)
    Write([]byte) (int, error)
}
```

### 最佳实践

#### 1. 接口设计原则
- **保持小而专注**：接口应该尽可能小，专注于特定行为
- **优先使用组合**：通过组合小接口构建复杂接口
- **面向行为设计**：接口名通常以 `-er` 结尾（如 `Reader`、`Writer`）

#### 2. 方法接收者选择
```go
// 何时使用指针接收者：
// 1. 需要修改接收者的状态
func (s *S) Set(age int) { s.Age = age }

// 2. 接收者是大型结构体（避免拷贝）
func (s *LargeStruct) Process() { }

// 3. 需要保证方法调用的一致性
type Counter struct { count int }
func (c *Counter) Increment() { c.count++ }  // 修改状态
func (c *Counter) Value() int { return c.count }  // 保持一致性，也用指针
```

#### 3. 接口使用建议
- **接收接口，返回具体类型**：函数参数使用接口，返回值使用具体类型
- **在使用方定义接口**：在需要使用的包中定义接口，而不是在实现的包中
- **避免过度抽象**：不要为了接口而接口，确保接口有实际价值

### 常见陷阱与注意事项

#### 1. nil 接口 vs nil 指针
```go
var i I = (*S)(nil)  // i != nil，但 i 包装的值为 nil
if i == nil {        // false！
    // 不会执行
}
```

#### 2. 值接收者与指针接收者混用
```go
type I interface { Method() }

type S struct{}
func (s S) Method() {}   // 值接收者

var s S
var i I = s   // ✅ 可以
var i I = &s  // ✅ 也可以

type T struct{}
func (t *T) Method() {}  // 指针接收者

var t T
var i I = t   // ❌ 编译错误！值类型无法实现需要指针接收者的接口
var i I = &t  // ✅ 正确
```

#### 3. 接口比较
```go
// 接口可以比较，但要小心
var i1, i2 I = &S{1}, &S{1}
fmt.Println(i1 == i2)  // false，不同的指针

var i3, i4 I = S{1}, S{1}
fmt.Println(i3 == i4)  // true，相同的值
```

### 性能考虑

- **接口调用有轻微开销**：比直接方法调用稍慢，但通常可忽略
- **避免频繁类型断言**：在性能敏感的代码中避免不必要的类型断言
- **接口切片的性能**：`[]interface{}` 不能直接转换为 `[]T`，需要逐个转换

---

后续接口相关的新知识点请继续补充到本文档相应章节。
