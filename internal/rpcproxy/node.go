package rpcproxy

import (
	"sync/atomic"
)

type Node struct {
	URL       string
	healthy   atomic.Bool
	failCount int32
}

// 标记为健康
func (n *Node) SetHealthy() {
	n.healthy.Store(true)
	n.failCount = 0
}

// 标记为不健康
func (n *Node) SetUnhealthy() {
	n.healthy.Store(false)
}

// 是否健康
func (n *Node) IsHealthy() bool {
	return n.healthy.Load()
}
