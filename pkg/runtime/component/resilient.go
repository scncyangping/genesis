// @Author: YangPing
// @Create: 2023/10/23
// @Description: 弹性组件配置

package component

import (
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

const (
	backoffTime = 5 * time.Second
)

type resilientComponent struct {
	log       *slog.Logger
	component Component
	cname     string
}

func NewResilientComponent(log *slog.Logger, component Component, name ...string) Component {
	c := &resilientComponent{
		log:       log,
		component: component,
	}
	if len(name) > 0 {
		c.cname = name[0]
	} else {
		c.cname = reflect.TypeOf(component).String()
	}
	return c
}

func (r *resilientComponent) Start(stop <-chan struct{}) error {
	r.log.Info(fmt.Sprintf("starting resilient component: %s", r.cname))
	for generationID := uint64(1); ; generationID++ {
		errCh := make(chan error, 1)
		go func(errCh chan<- error) {
			defer close(errCh)
			// recover from a panic
			defer func() {
				if e := recover(); e != nil {
					if err, ok := e.(error); ok {
						errCh <- err
					} else {
						errCh <- errors.Errorf("component: %s  error: %v", r.cname, e)
					}
				}
			}()

			errCh <- r.component.Start(stop)
		}(errCh)
		select {
		case <-stop:
			r.log.Info(fmt.Sprintf("resilient done component: %s", r.cname))
			return nil
		case err := <-errCh:
			if err != nil {
				r.log.With("generationID", generationID).Error(fmt.Sprintf(" component: %s terminated with an error: %v", r.cname, err))
			}
		}
		<-time.After(backoffTime)
	}
}

func (r *resilientComponent) NeedLeaderElection() bool {
	return r.component.NeedLeaderElection()
}
