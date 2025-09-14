package timesnippet

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	// 周期性触发
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	//只触发一次
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-ticker.C:
			t.Log("定时触发分支")
		default:
			t.Log("非定时分支")
			time.Sleep(time.Second)
		}

	}
}
