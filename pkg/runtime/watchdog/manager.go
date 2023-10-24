// @Author: YangPing
// @Create: 2023/10/23
// @Description: watchdog

package watchdog

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Manager struct {
	watchDogs []Watchdog
	logger    *zap.SugaredLogger
}

func NewWatchDogManager(log *zap.SugaredLogger) *Manager {
	return &Manager{
		logger: log,
	}
}

func (w *Manager) WithWatchDog(wcs Watchdog) *Manager {
	w.watchDogs = append(w.watchDogs, wcs)
	return w
}

func (w *Manager) Start(stop <-chan struct{}) error {
	w.logger.Infof("watchdog start")
	errCh := make(chan error, 1)
	defer close(errCh)

	for _, wd := range w.watchDogs {
		go func(w Watchdog, errCh chan<- error) {
			// recover from a panic
			defer func() {
				if e := recover(); e != nil {
					if err, ok := e.(error); ok {
						errCh <- err
					} else {
						errCh <- errors.Errorf("%v", e)
					}
				}
			}()
			w.Start(stop)
		}(wd, errCh)
	}

	select {
	case <-stop:
		w.logger.Info("watchdog done")
		return nil
	case err := <-errCh:
		if err != nil {
			w.logger.Error(errors.Wrap(err, "watchdog terminated with an error"))
		}
	}
	return nil
}

func (w *Manager) NeedLeaderElection() bool {
	return false
}
