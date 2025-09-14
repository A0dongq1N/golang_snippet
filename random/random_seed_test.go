package randomsnippet

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// 演示 rand.Seed(time.Now().Unix()) 在并发环境下的问题
func TestSeedConcurrencyProblem(t *testing.T) {
	const numGoroutines = 10
	const numRands = 5

	t.Run("问题示例：相同种子导致重复", func(t *testing.T) {
		results := make([][]int, numGoroutines)
		var wg sync.WaitGroup

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				// 问题：多个 goroutine 在同一时刻调用，可能得到相同的种子
				rand.Seed(time.Now().Unix()) // Unix() 精度只到秒级

				// 生成随机数序列
				nums := make([]int, numRands)
				for j := 0; j < numRands; j++ {
					nums[j] = rand.Intn(1000)
				}
				results[goroutineID] = nums
			}(i)
		}
		wg.Wait()

		// 检查重复序列
		sequenceMap := make(map[string][]int)
		for i, seq := range results {
			key := fmt.Sprintf("%v", seq)
			sequenceMap[key] = append(sequenceMap[key], i)
		}

		t.Logf("生成的随机数序列：")
		for i, seq := range results {
			t.Logf("Goroutine %d: %v", i, seq)
		}

		duplicates := 0
		for sequence, goroutineIDs := range sequenceMap {
			if len(goroutineIDs) > 1 {
				duplicates++
				t.Logf("重复序列 %s 出现在 Goroutines: %v", sequence, goroutineIDs)
			}
		}

		if duplicates > 0 {
			t.Logf("发现 %d 个重复的随机数序列！", duplicates)
		} else {
			t.Logf("没有发现重复序列（可能是运行时间差异导致）")
		}
	})

	t.Run("解决方案1：使用纳秒级种子", func(t *testing.T) {
		results := make([][]int, numGoroutines)
		var wg sync.WaitGroup

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				// 解决方案：使用纳秒级时间戳 + goroutine ID
				seed := time.Now().UnixNano() + int64(goroutineID)
				rand.Seed(seed)

				nums := make([]int, numRands)
				for j := 0; j < numRands; j++ {
					nums[j] = rand.Intn(1000)
				}
				results[goroutineID] = nums
			}(i)
		}
		wg.Wait()

		// 检查是否还有重复
		sequenceMap := make(map[string][]int)
		for i, seq := range results {
			key := fmt.Sprintf("%v", seq)
			sequenceMap[key] = append(sequenceMap[key], i)
		}

		duplicates := 0
		for _, goroutineIDs := range sequenceMap {
			if len(goroutineIDs) > 1 {
				duplicates++
			}
		}
		t.Logf("使用纳秒级种子后，重复序列数: %d", duplicates)
	})

	t.Run("解决方案2：独立随机数生成器", func(t *testing.T) {
		results := make([][]int, numGoroutines)
		var wg sync.WaitGroup

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				// 最佳解决方案：每个 goroutine 独立的随机数生成器
				source := rand.NewSource(time.Now().UnixNano() + int64(goroutineID))
				rng := rand.New(source)

				nums := make([]int, numRands)
				for j := 0; j < numRands; j++ {
					nums[j] = rng.Intn(1000)
				}
				results[goroutineID] = nums
			}(i)
		}
		wg.Wait()

		// 检查重复
		sequenceMap := make(map[string][]int)
		for i, seq := range results {
			key := fmt.Sprintf("%v", seq)
			sequenceMap[key] = append(sequenceMap[key], i)
		}

		duplicates := 0
		for _, goroutineIDs := range sequenceMap {
			if len(goroutineIDs) > 1 {
				duplicates++
			}
		}
		t.Logf("使用独立生成器后，重复序列数: %d", duplicates)
	})
}

// 更极端的种子冲突示例
func TestSeedCollisionDemo(t *testing.T) {
	t.Run("故意制造种子冲突", func(t *testing.T) {
		const numWorkers = 5
		results := make([]int, numWorkers)
		var wg sync.WaitGroup

		// 让所有 goroutine 同时启动，使用相同的时间戳
		startTime := time.Now().Unix()

		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// 故意使用相同的种子
				rand.Seed(startTime)

				// 生成第一个随机数
				results[workerID] = rand.Intn(10000)
			}(i)
		}
		wg.Wait()

		t.Logf("使用相同种子 %d 的结果:", startTime)
		for i, result := range results {
			t.Logf("Worker %d: %d", i, result)
		}

		// 检查是否有相同结果
		valueMap := make(map[int][]int)
		for i, val := range results {
			valueMap[val] = append(valueMap[val], i)
		}

		for val, workers := range valueMap {
			if len(workers) > 1 {
				t.Logf("值 %d 在多个 workers 中重复: %v", val, workers)
			}
		}
	})
}

// 演示 math/rand 并发冲突问题
func TestRandomConcurrency(t *testing.T) {
	const numGoroutines = 100
	const numIterations = 1000

	// 方法1：使用全局随机数生成器（有并发问题）
	t.Run("unsafe_global_rand", func(t *testing.T) {
		results := make([]int, numGoroutines*numIterations)
		var wg sync.WaitGroup

		start := time.Now()
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				for j := 0; j < numIterations; j++ {
					// 全局 rand 包不是并发安全的
					// 多个 goroutine 同时调用可能导致数据竞争
					val := rand.Intn(1000)
					results[goroutineID*numIterations+j] = val
				}
			}(i)
		}
		wg.Wait()
		elapsed := time.Since(start)
		t.Logf("全局 rand（不安全）: %v", elapsed)

		// 检查是否有重复值（不是完美的并发检测，但能反映问题）
		uniqueVals := make(map[int]int)
		for _, val := range results {
			uniqueVals[val]++
		}
		t.Logf("生成了 %d 个不同的值", len(uniqueVals))
	})

	// 方法2：每个 goroutine 创建独立的随机数生成器（安全但可能重复）
	t.Run("separate_rand_sources", func(t *testing.T) {
		results := make([]int, numGoroutines*numIterations)
		var wg sync.WaitGroup

		start := time.Now()
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				// 每个 goroutine 有独立的随机数生成器
				rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(goroutineID)))
				for j := 0; j < numIterations; j++ {
					val := rng.Intn(1000)
					results[goroutineID*numIterations+j] = val
				}
			}(i)
		}
		wg.Wait()
		elapsed := time.Since(start)
		t.Logf("独立随机源: %v", elapsed)
	})

	// 方法3：使用 crypto/rand（线程安全但较慢）
	t.Run("crypto_rand", func(t *testing.T) {
		results := make([]int, numGoroutines*numIterations)
		var wg sync.WaitGroup

		start := time.Now()
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				for j := 0; j < numIterations; j++ {
					// crypto/rand 是线程安全的
					buf := make([]byte, 4)
					if _, err := crand.Read(buf); err != nil {
						t.Errorf("crypto/rand error: %v", err)
						return
					}
					val := int(binary.BigEndian.Uint32(buf)) % 1000
					results[goroutineID*numIterations+j] = val
				}
			}(i)
		}
		wg.Wait()
		elapsed := time.Since(start)
		t.Logf("crypto/rand（安全）: %v", elapsed)
	})

	// 方法4：使用互斥锁保护全局随机数生成器
	t.Run("mutex_protected_rand", func(t *testing.T) {
		results := make([]int, numGoroutines*numIterations)
		var wg sync.WaitGroup
		var mu sync.Mutex

		start := time.Now()
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				for j := 0; j < numIterations; j++ {
					mu.Lock()
					val := rand.Intn(1000)
					mu.Unlock()
					results[goroutineID*numIterations+j] = val
				}
			}(i)
		}
		wg.Wait()
		elapsed := time.Since(start)
		t.Logf("互斥锁保护的 rand: %v", elapsed)
	})
}

// 演示数据竞争检测
func TestRandomDataRace(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过数据竞争测试")
	}

	// 运行命令：go test -race ./random_snippet
	// 这将检测到数据竞争

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = rand.Intn(100) // 可能导致数据竞争
			}
		}()
	}
	wg.Wait()
}
