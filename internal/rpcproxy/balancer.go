package rpcproxy

import (
	"sync/atomic"
)

// EndpointSelector 定义“选一个节点”的接口
type EndpointSelector interface {
	Next() *Node
}

type RoundRobinSelector struct {
	nodes   []*Node
	counter uint32
}

func NewRoundRobinSelector(nodes []*Node) *RoundRobinSelector {
	return &RoundRobinSelector{
		nodes: nodes,
	}
}

func (r *RoundRobinSelector) Next() *Node {
	n := uint32(len(r.nodes))

	if n == 0 {
		return nil
	}

	for i := 0; i < int(n); i++ {
		idx := atomic.AddUint32(&r.counter, 1)
		node := r.nodes[idx%n]
		if node.IsHealthy() {
			return node
		}

	}

	return nil

}

// RoundRobin 是最简单、最稳定的 LB 实现
type RoundRobin struct {
	endpoints []string
	counter   uint32
}

// NewRoundRobin 构造
func NewRoundRobin(endpoints []string) *RoundRobin {
	return &RoundRobin{
		endpoints: endpoints,
		counter:   0,
	}
}

// Next 返回下一个可用 endpoint（并发安全）
func (rr *RoundRobin) Next() string {
	n := uint32(len(rr.endpoints))
	if n == 0 {
		return ""
	}

	idx := atomic.AddUint32(&rr.counter, 1)
	return rr.endpoints[idx%n]
}
