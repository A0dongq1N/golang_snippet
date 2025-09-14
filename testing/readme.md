## Go Testing 知识库

面向本仓库的单元/集成测试实践总结，持续补充。

### 基本约定
- **文件命名**：测试文件必须为 `*_test.go`。
- **测试函数**：`func TestXxx(t *testing.T)`；基准 `func BenchmarkXxx(b *testing.B)`；示例 `func ExampleXxx()`。
- **包名选择**：与被测包同名或使用 `package_name_test` 进行黑盒测试。

### 组织测试与测试套件
- **子测试与表驱动**：在一个 `TestXxx` 中使用 `t.Run` 组织多个子测试。
  ```go
  func TestDouble(t *testing.T) {
      cases := []struct{ name string; in, want int }{
          {"base", 1, 2}, {"zero", 0, 0}, {"negative", -3, -6},
      }
      for _, tc := range cases {
          tc := tc              // 关键：为闭包拷贝当前迭代变量
          t.Run(tc.name, func(t *testing.T) {
              // t.Parallel() 可选：并行执行子测试
              if got := tc.in*2; got != tc.want {
                  t.Fatalf("got=%d want=%d", got, tc.want)
              }
          })
      }
  }
  ```
- **并行执行 `t.Parallel()`**：将当前（子）测试标记为可与其他并行测试同时运行。并发上限由 `-parallel` 控制（默认 10）。调用前的代码串行、调用后的代码并发执行。
- **闭包捕获陷阱**：循环变量会被闭包按“同一个变量”捕获，需用 `tc := tc` 在每轮迭代创建副本，避免并发/延迟执行时数据错乱。

### 前置与后置（生命周期）
- **层级**：
  - 包级：`TestMain(m *testing.M)` 负责一次性初始化/销毁（数据库、容器、外部服务等）。
  - 用例级：每个 `TestXxx` 内部准备与清理。
  - 子测试级：每个 `t.Run` 的专属准备与清理。
- **常见前置类型**：
  - 环境变量：`t.Setenv`（测试结束自动回滚）。
  - 临时文件/目录：`t.TempDir`（自动清理）。
  - 外部依赖：`httptest.NewServer`、内存/本地 DB、容器等。
  - 共享重资源：配合 `sync.Once` 实现仅初始化一次的资源构建。
- **后置优先级**：
  - 首选 `t.Cleanup(func(){ ... })`，作用域清晰，随测试结束自动执行。
  - `defer` 仅限当前函数作用域；注意在 `TestMain` 里 `os.Exit` 会跳过 `defer`。

### TestMain 与 m.Run/退出码
```go
func TestMain(m *testing.M) {
    // 前置
    code := m.Run()  // 运行本包全部测试/基准/示例（按参数决定）
    // 后置
    os.Exit(code)    // 返回真实退出码；否则失败也可能被误判为通过
}
```
- `m.Run()` 返回进程应退出的状态码：0 全部通过，非 0 有失败/中断。
- 一旦自定义了 `TestMain`，必须用 `os.Exit(code)` 结束；否则进程默认 0 退出。
- 需要清理的逻辑要在 `os.Exit(code)` 之前完成（`os.Exit` 不执行 `defer`）。

### 日志与输出
- `TestMain` 里不能用 `t.Log`。可用：
  - `log`：`log.SetFlags(log.LstdFlags|log.Lshortfile); log.Println("...")`
  - `log/slog`（Go 1.21+）：`slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))`
- 普通测试函数内优先 `t.Log/t.Logf`；想总是看到日志请用 `-v`：`go test -v`。

### 失败与跳过
- 立即失败并中断：`t.Fatal/t.Fatalf/t.FailNow`。
- 标记失败但继续：`t.Error/t.Errorf/t.Fail`。
- 跳过：`t.Skip/t.Skipf`（不算失败）；也可用构建标签或检查环境条件决定跳过。

### 只运行部分测试
- 只跑某个用例：`go test -run '^TestFoo$' ./pkg`
- 只跑子测试：`go test -run 'TestFoo/^case1$' ./pkg`
- 只跑基准：`go test -bench '^BenchmarkX$' -run '^$' ./pkg`

### 基准测试与示例
- 基准：
  ```go
  func BenchmarkX(b *testing.B) {
      for i := 0; i < b.N; i++ { _ = work() }
  }
  // 运行：go test -bench . -benchmem ./pkg
  ```
- 示例：在函数注释中包含 `// Output:` 期望输出，可参与测试验证。

### 并发与隔离
- 避免共享可变全局状态；需要时：
  - 用 `t.TempDir` 隔离文件系统、`t.Setenv` 隔离环境变量。
  - 使用锁/原子操作或为每个测试准备独立实例。
- 数据竞争检测：`go test -race ./...`。

### 目录与测试数据
- 使用 `testdata/` 存放静态样例文件，Go 工具链会忽略编译。
- 临时文件优先 `t.TempDir`，避免污染工作区。

### 第三方工具（可选）
- `testify/assert` 与 `testify/require`：断言风格，`require.*` 失败即中断，`assert.*` 失败但继续。
- `testify/suite`：提供 `SetupSuite/TeardownSuite/SetupTest/TeardownTest` 等生命周期，适合复杂集成测试。

### 常用命令速查
- 基本：`go test ./...`，详细输出：`go test -v ./...`
- 过滤：`-run`、`-bench`、`-benchmem`、`-parallel N`
- 可靠性：`-race`、`-count=1`（禁缓存）、`-timeout 5m`
- 覆盖率：`go test -coverprofile=cover.out ./... && go tool cover -html=cover.out`

### 注意事项清单
- `*_test.go` 与 `TestXxx` 命名规范确保 IDE/工具链识别。
- `TestMain` 中务必 `os.Exit(m.Run())`；清理逻辑放在其前。
- 优先使用 `t.Cleanup` 管理后置；改环境用 `t.Setenv`；文件用 `t.TempDir`。
- 并发子测试请在 `t.Run` 开头调用 `t.Parallel()`，并使用 `tc := tc` 防闭包陷阱。
- 依赖全局顺序或共享资源的测试不要并行。
- 使用 `-run`/正则只运行目标测试，避免被其他包失败干扰。

---

后续新增的测试相关内容（模式、最佳实践、工具链等）请继续补充到本文件，并在上面相应章节归类或新增章节。


