package rpcproxy

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var ErrNoEndpoints = errors.New("no rpc endpoints configured")

type Proxy struct {
	nodes    []*Node
	selector EndpointSelector
	client   *http.Client
}

func NewProxy(endpoints []string) (*Proxy, error) {

	if len(endpoints) == 0 {
		return nil, ErrNoEndpoints
	}

	nodes := make([]*Node, 0, len(endpoints))

	for _, e := range endpoints {
		//检查 URL 合法性
		if _, err := url.Parse(e); err != nil {
			return nil, err
		}
		nodes = append(nodes, &Node{
			URL: e,
		})
	}

	selector := NewRoundRobinSelector(nodes)

	return &Proxy{
		nodes:    nodes,
		selector: selector,
		client:   &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// 	u, err := url.Parse(upstream)
// 	if err != nil {
// 		return nil, err
// 	}
// 	client := &http.Client{Timeout: 10 * time.Second}

// 	return &Proxy{
// 		upstream: u,
// 		client:   client,
// 	}, nil
// }

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	node := p.selector.Next()
	if node == nil {
		http.Error(w, "o healthy upstream available", http.StatusServiceUnavailable)
		return
	}

	// 构建上游请求
	req, err := http.NewRequest(r.Method, node.URL, r.Body)
	if err != nil {
		log.Println("build upstream request error:", err)
		http.Error(w, "bad gateway", http.StatusBadGateway)
		return
	}

	req.Header = r.Header.Clone()

	// req, err := http.NewRequest(http.MethodPost, p.upstream.String(), r.Body)
	// if err != nil {
	// 	http.Error(w, "Failed to create request", http.StatusInternalServerError)
	// 	return
	// }
	// req.Header = r.Header

	resp, err := p.client.Do(req)
	if err != nil {
		log.Println("upstream request error:", err)
		http.Error(w, "Failed to reach upstream", http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	for k, vs := range resp.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println("copy response body error:", err)
	}
}
