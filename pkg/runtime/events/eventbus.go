// @Author: YangPing
// @Create: 2023/10/23
// @Description: 事件总线配置

package events

import (
	"sync"

	"github.com/pkg/errors"
)

func NewEventBus() *EventBus {
	return &EventBus{}
}

type EventBus struct {
	mtx         sync.RWMutex
	subscribers []chan Event
}

func (b *EventBus) New() Listener {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	events := make(chan Event, 10)
	b.subscribers = append(b.subscribers, events)
	return &reader{
		events: events,
	}
}

func (b *EventBus) Send(event Event) {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	// 所有订阅者 都发送
	for _, s := range b.subscribers {
		s <- event
	}
}

type reader struct {
	events chan Event
}

func (k *reader) Receive(stop <-chan struct{}) (Event, error) {
	select {
	case event, ok := <-k.events:
		if !ok {
			return nil, errors.New("end of events channel")
		}
		return event, nil
	case <-stop:
		return nil, ListenerStoppedErr
	}
}
