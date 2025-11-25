package rpcproxy

import (
	"sync/atomic"
)

// EndpointSelector 定义“选一个节点”的接口
type EndpointSelector interface {
	Next() string
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
