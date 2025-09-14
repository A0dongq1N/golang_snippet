package redis

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// getEnvWithDefault 获取环境变量，如果不存在则返回默认值
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 创建Redis客户端的辅助函数
func createRedisClient() *redis.Client {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Warning: Could not load config.env file: %v", err)
	}

	addr := getEnvWithDefault("REDIS_ADDR", "127.0.0.1:6379")
	password := getEnvWithDefault("REDIS_PASSWORD", "")

	return redis.NewClient(&redis.Options{
		Password:     password,
		Addr:         addr,
		DB:           0,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolTimeout:  6 * time.Second,
	})
}

/*
Redis Pipeline 的核心作用是**把多条命令打包一次性发给服务器，再一次性拿回结果**，从而**减少网络往返次数（RTT）**，提高吞吐量。

优点概括：
1. **省 RTT**：N 条命令只需 1 次往返，延迟大幅降低。
2. **高吞吐**：单位时间内可处理更多请求，CPU 利用率更高。
3. **无阻塞**：客户端把命令放进本地缓冲区即可继续干别的，最后统一 `Exec` 拿回结果。
4. **简单易用**：go-redis 里 `pipe := rdb.Pipeline()` → 攒命令 → `pipe.Exec(ctx)` 即可，代码改动极小。
5. **与事务无关**：Pipeline 只“打包”，不保证原子性；需要原子性请用 `TxPipeline` 或 `Watch`+`TxPipelined`。

一句话：**Pipeline 是“批处理”，不是“事务”；它让网络成为瓶颈前再快一点。**

*/

func TestRedisPipeline(t *testing.T) {
	rdb := createRedisClient()
	defer rdb.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 测试Redis连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		t.Skipf("Redis server not available: %v", err)
	}

	// 清理测试数据
	rdb.Del(ctx, "key1", "key2", "counter")

	log.Printf("=== 测试 Pipeline（批量执行，非事务）===")

	// 使用Pipeline进行批量操作
	pipe := rdb.Pipeline()
	setCmd1 := pipe.Set(ctx, "key1", "value1", 0)
	setCmd2 := pipe.Set(ctx, "key2", "value2", 0)
	incrCmd := pipe.Incr(ctx, "counter")

	// 执行Pipeline
	cmds, err := pipe.Exec(ctx)
	assert.NoError(t, err)
	log.Printf("Pipeline执行了 %d 个命令", len(cmds))

	// 验证结果
	assert.NoError(t, setCmd1.Err())
	assert.NoError(t, setCmd2.Err())
	assert.NoError(t, incrCmd.Err())

	// 清理测试数据
	rdb.Del(ctx, "key1", "key2", "counter")
}
