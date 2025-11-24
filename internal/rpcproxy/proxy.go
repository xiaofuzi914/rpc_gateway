package rpcproxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Proxy struct {
	upstream *url.URL
	client   *http.Client
}

func NewProxy(upstream string) (*Proxy, error) {

	u, err := url.Parse(upstream)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 10 * time.Second}

	return &Proxy{
		upstream: u,
		client:   client,
	}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := http.NewRequest(http.MethodPost, p.upstream.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header

	resp, err := p.client.Do(req)
	if err != nil {
		log.Println("upstream request error:", err)
		http.Error(w, "Failed to reach upstream", http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	for k, vs := range resp.Header {
		for _, vv := range vs {
			w.Header().Add(k, vv)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
