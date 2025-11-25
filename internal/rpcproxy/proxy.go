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
	selector EndpointSelector
	client   *http.Client
}

func NewProxy(endpoints []string) (*Proxy, error) {

	if len(endpoints) == 0 {
		return nil, ErrNoEndpoints
	}

	for _, e := range endpoints {
		if _, err := url.Parse(e); err != nil {
			return nil, err
		}
	}

	return &Proxy{
		selector: NewRoundRobin(endpoints),
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

	upstream := p.selector.Next()
	if upstream == "" {
		http.Error(w, "No upstream available", http.StatusBadGateway)
		return
	}

	// 构建上游请求
	req, err := http.NewRequest(http.MethodPost, upstream, r.Body)
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
	io.Copy(w, resp.Body)
}
