// @Author: YangPing
// @Create: 2023/10/23
// @Description: watchdog

package watchdog

import (
	"github.com/robfig/cron"
	"time"
)

type Watchdog interface {
	Start(stop <-chan struct{})
}

type CronWatchdog struct {
	Cron   string
	Task   func()
	OnStop func()
}

func (c *CronWatchdog) Start(stop <-chan struct{}) {
	cro := cron.New()
	cro.AddFunc(c.Cron, c.Task)

	cro.Start()

	select {
	case <-stop:
		cro.Stop()
		if c.OnStop != nil {
			c.OnStop()
		}
		return
	}
}

type SimpleWatchdog struct {
	NewTicker  func() *time.Ticker
	OnTick     func() error
	OnError    func(error)
	OnStop     func()
	ExecuteNow bool
}

func (w *SimpleWatchdog) Start(stop <-chan struct{}) {
	if w.ExecuteNow {
		if err := w.OnTick(); err != nil {
			w.OnError(err)
		}
	}

	ticker := w.NewTicker()
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := w.OnTick(); err != nil {
				w.OnError(err)
			}
		case <-stop:
			if w.OnStop != nil {
				w.OnStop()
			}
			return
		}
	}
}

func (w *SimpleWatchdog) WithExecuteNow(we bool) *SimpleWatchdog {
	w.ExecuteNow = we
	return w
}
