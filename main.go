package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const upstreamURL = "https://cloudflare-eth.com"

func handleRequest(w http.ResponseWriter, r *http.Request) {

	// 1. 我们只处理 POST 请求 (因为 JSON-RPC 都是 POST)
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. 核心挑战：如何把用户发来的 r.Body 转发给 upstreamURL？
	// TODO: 在这里写转发逻辑
	// 提示：你需要创建一个新的 request，然后用 http.DefaultClient.Do() 发送它

	req, err := http.NewRequest(http.MethodPost, upstreamURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// 3. 复制用户请求的 headers 到新的请求中
	req.Header = r.Header.Clone()

	// 4. 发送请求到上游节点
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to reach upstream", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 5. 把上游节点的响应头和状态码复制回用户响应
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	// 6. 把上游节点的响应体复制回用户响应体
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Println("Failed to copy response body:", err)
	}

}

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("🚀 Gateway running on :8080 forwarding to", upstreamURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
