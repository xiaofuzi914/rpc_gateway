package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/xiaofuzi914/rpc-gateway/internal/config"
	"github.com/xiaofuzi914/rpc-gateway/internal/rpcproxy"
)

func main() {

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	mux := http.NewServeMux()

	// Print loaded configuration for debugging
	// log.Printf("Loaded config: %+v", cfg)

	// // 暂时先只用 ethereum 的第一个节点，后面再做多节点、负载均衡
	// ethCfg, ok := cfg.Chains["ethereum"]
	// if !ok || len(ethCfg.RPCEndpoints) == 0 {
	// 	log.Fatal("ethereum chain config or endpoints missing")
	// }

	// proxy, err := rpcproxy.NewProxy(ethCfg.RPCEndpoints[1])
	// if err != nil {
	// 	log.Fatal("Failed to create proxy:", err)
	// }

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Health check requested")
		// w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	proxies := make(map[string]*rpcproxy.Proxy)

	for name, chain := range cfg.Chains {
		p, err := rpcproxy.NewProxy(chain.RPCEndpoints)
		if err != nil {
			log.Fatalf("Failed to create proxy for chain %s: %v", name, err)
		}

		p.StartHealthCheck() // 启动健康检查

		proxies[name] = p
		log.Printf("Proxy created for chain %s with endpoints: %v", name, chain.RPCEndpoints)
	}

	mux.HandleFunc("/rpc/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("rpc run start")
		path := strings.TrimPrefix(r.URL.Path, "/rpc/")
		parts := strings.Split(path, "/")
		chain := parts[0]

		if chain == "" {
			http.Error(w, "missing chain named", http.StatusBadRequest)
			return
		}

		p, ok := proxies[chain]

		if !ok {
			http.Error(w, "unknown chain: "+chain, http.StatusBadRequest)
			return
		}
		p.ServeHTTP(w, r)
	})

	log.Println("listening on", cfg.Server.Addr)
	if err := http.ListenAndServe(cfg.Server.Addr, mux); err != nil {
		log.Fatal(err)
	}
}
