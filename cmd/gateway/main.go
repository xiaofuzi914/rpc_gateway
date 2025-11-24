package main

import (
	"log"
	"net/http"

	"github.com/xiaofuzi914/rpc-gateway/internal/config"
	"github.com/xiaofuzi914/rpc-gateway/internal/rpcproxy"
)

func main() {

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// Print loaded configuration for debugging
	log.Printf("Loaded config: %+v", cfg)

	// 暂时先只用 ethereum 的第一个节点，后面再做多节点、负载均衡
	ethCfg, ok := cfg.Chains["ethereum"]
	if !ok || len(ethCfg.RPCEndpoints) == 0 {
		log.Fatal("ethereum chain config or endpoints missing")
	}

	proxy, err := rpcproxy.NewProxy(ethCfg.RPCEndpoints[1])
	if err != nil {
		log.Fatal("Failed to create proxy:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Health check requested")
		// w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.Handle("/rpc", proxy)

	log.Println("listening on", cfg.Server.Addr)
	if err := http.ListenAndServe(cfg.Server.Addr, mux); err != nil {
		log.Fatal(err)
	}
}
