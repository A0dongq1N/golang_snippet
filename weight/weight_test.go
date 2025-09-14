package weightsnippet

import (
	"fmt"
	"testing"

	"github.com/smallnest/weighted"
)

// 结论（两者都是加权轮询，但策略不同）
// weighted.SW（Smooth Weighted）: 平滑加权轮询（Nginx 同款算法）
// 机制: 每次选择前为每个节点累加“当前权重”，选出最大的，再减去总权重。
// 特点: 在短时间窗口内分配更均匀、波动更小（更“平滑”）。
// 代价: 计算略复杂（但仍是 O(n)），更适合对“连续性/平滑度”敏感的场景。

// weighted.RRW（Round Robin Weighted, LVS 算法）:
// 机制: 维护 gcd、maxW、cw 与下标 i，按权重循环扫描选择。
// 特点: 实现更轻量、吞吐更高，但短窗口可能出现“同一节点连续命中”（不如 SW 平滑）。
// 适用: 追求简单/性能、对短期平滑度不敏感，或权重较大时（gcd 降低步进次数）。

// 相同点
// 都是 O(n) 选择；都非并发安全，需在并发调用 Next() 时加锁保护。

// 如何选择
// 需要“更均匀的瞬时分布/更平滑的体验”→ 选 SW
// 需要“实现更简单/计算更轻/吞吐更高”→ 选 RRW

func TestExampleSW_Next(t *testing.T) {
	w := &weighted.SW{}
	w.Add("a", 5)
	w.Add("b", 2)
	w.Add("c", 3)

	for i := 0; i < 10; i++ {
		fmt.Printf("%s ", w.Next())
	}
}

func TestExampleRRW_Next(t *testing.T) {
	w := weighted.RRW{}
	w.Add("a", 5)
	w.Add("b", 2)
	w.Add("c", 3)

	for i := 0; i < 10; i++ {
		fmt.Printf("%s ", w.Next())
	}
}
