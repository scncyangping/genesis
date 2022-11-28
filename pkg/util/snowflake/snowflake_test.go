package snowflake

import (
	"fmt"
	"github.com/sony/sonyflake"
	"sync"
	"testing"
	"time"
)

func TestNextId(t *testing.T) {
	var wg sync.WaitGroup
	var check sync.Map
	t1 := time.Now()
	sonyflake := sonyflake.NewSonyflake(struct {
		StartTime      time.Time
		MachineID      func() (uint16, error)
		CheckMachineID func(uint16) bool
	}{StartTime: time.Now(), MachineID: func() (uint16, error) {
		return 5, nil
	}, CheckMachineID: func(u uint16) bool {
		return true
	}})
	for i := 0; i < 200000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, _ := sonyflake.NextID()
			val := id
			if _, ok := check.Load(val); ok {
				// id冲突检查
				t.Error(fmt.Errorf("error#unique: val:%v", val))
				return
			}
			check.Store(val, 0)
		}()
	}
	wg.Wait()
	elapsed := time.Since(t1).Seconds()
	println(int64(elapsed))
}

func TestNextId2(t *testing.T) {
	var check sync.Map
	var wg sync.WaitGroup

	t1 := time.Now()
	for i := 0; i < 200000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val := NextId()
			if _, ok := check.Load(val); ok {
				// id冲突检查
				t.Error(fmt.Errorf("error#unique: val:%v", val))
				return
			}
			check.Store(val, 0)
		}()
	}
	wg.Wait()

	elapsed := time.Since(t1).Seconds()
	println(int64(elapsed))
}
