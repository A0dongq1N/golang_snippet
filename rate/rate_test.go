package ratesnippet

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

//TODO 2025年09月14日00:46:08 还没懂

func TestRate(t *testing.T) {
	lim := rate.NewLimiter(100, 10) // 100 QPS，允许 10 的瞬时突发
	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	count := 0
	for {
		select {
		case <-tick.C:
			fmt.Println("qps:", count)
			count = 0
		// 没有匹配到case就会执行default
		default:
			// 如果令牌桶中的牌数量小于100就++，否则就执行else,因为很快就执行到了100，所以大部分时间都是sleep,所以定时输出100
			if lim.Allow() {
				count++
			} else {
				fmt.Printf("当前桶内令牌 = %.2f\n", lim.TokensAt(time.Now()))
				// 丢弃或小睡一会儿
				time.Sleep(1 * time.Millisecond)
			}
		}
	}

}
