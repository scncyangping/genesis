package watchdog

import (
	"time"

	"github.com/robfig/cron"
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
	NewTicker func() *time.Ticker
	OnTick    func() error
	OnError   func(error)
	OnStop    func()
}

func (w *SimpleWatchdog) Start(stop <-chan struct{}) {
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
