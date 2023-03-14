package watchdog

import (
	"fmt"
	"genesis/pkg/util"
	"testing"
	"time"
)

func TestCronWatchdog(t *testing.T) {
	cc := CronWatchdog{
		Cron: "*/5 * * * * ?",
		Task: func() {
			fmt.Println(util.NowDateTimeFormat())
		},
		OnStop: func() {
			fmt.Printf("stop")
		},
	}
	stop := make(chan struct{})
	go cc.Start(stop)
	time.Sleep(20 * time.Second)

	stop <- struct{}{}
	fmt.Println("stop >>>>>>")
	time.Sleep(10 * time.Second)
	fmt.Println("end >>>>>>")
}
