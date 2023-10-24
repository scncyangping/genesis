// @Author: YangPing
// @Create: 2023/10/23
// @Description: 事件定义

package events

import (
	"github.com/pkg/errors"
)

type Event any

var ListenerStoppedErr = errors.New("listener closed")

type Listener interface {
	Receive(stop <-chan struct{}) (Event, error)
}

type Emitter interface {
	Send(Event)
}

type ListenerFactory interface {
	New() Listener
}
