// @Author: YangPing
// @Create: 2023/10/23
// @Description: 组件定义

package leader

import "genesis/pkg/runtime/component"

type noopLeaderElector struct {
	alwaysLeader bool
	callbacks    []component.LeaderCallbacks
}

func NewAlwaysLeaderElector() component.LeaderElector {
	return &noopLeaderElector{
		alwaysLeader: true,
	}
}

func NewNeverLeaderElector() component.LeaderElector {
	return &noopLeaderElector{
		alwaysLeader: false,
	}
}

func (n *noopLeaderElector) AddCallbacks(callbacks component.LeaderCallbacks) {
	n.callbacks = append(n.callbacks, callbacks)
}

func (n *noopLeaderElector) IsLeader() bool {
	return n.alwaysLeader
}

func (n *noopLeaderElector) Start(stop <-chan struct{}) {
	if n.alwaysLeader {
		for _, callback := range n.callbacks {
			callback.OnStartedLeading()
		}
	}
}

var _ component.LeaderElector = &noopLeaderElector{}
