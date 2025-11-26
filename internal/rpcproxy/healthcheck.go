package rpcproxy

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type rpcReq struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	Id      int    `json:"id"`
}

func (p *Proxy) StartHealthCheck() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			for _, node := range p.nodes {
				checkNode(p.client, node)
			}
		}

	}()
}

func checkNode(client *http.Client, n *Node) {
	body, _ := json.Marshal(rpcReq{
		Jsonrpc: "2.0",
		Method:  "ping",
		Params:  []any{},
		Id:      1,
	})

	resp, err := client.Post(n.URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		n.failCount++

		if n.failCount >= 3 && n.IsHealthy() {
			log.Printf("Health check failed for %s: %v", n.URL, err)
			n.SetUnhealthy()
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		n.failCount++
		if n.failCount >= 3 && n.IsHealthy() {
			log.Println("[health] DOWN (bad status):", n.URL)
			n.SetUnhealthy()
		}
		return
	}

	if !n.IsHealthy() {
		log.Println("[health] UP:", n.URL)
	}
	n.SetUnhealthy()
}
