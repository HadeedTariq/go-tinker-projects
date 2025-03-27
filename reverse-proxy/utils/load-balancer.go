package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	servers []*url.URL // List of backend servers
	mu      sync.Mutex
	current int
}

// Initialize load balancer with backend servers
func NewLoadBalancer(serverList []string) *LoadBalancer {
	var servers []*url.URL
	for _, addr := range serverList {
		serverURL, _ := url.Parse(addr)
		servers = append(servers, serverURL)
	}
	return &LoadBalancer{servers: servers}
}
func (lb *LoadBalancer) getNextServer() *url.URL {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server := lb.servers[lb.current]
	lb.current = (lb.current + 1) % len(lb.servers) // Round-robin cycling
	return server
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextServer()

	proxy := httputil.NewSingleHostReverseProxy(targetServer)
	r.URL.Host = targetServer.Host
	r.URL.Scheme = targetServer.Scheme
	r.Host = targetServer.Host

	fmt.Printf("Forwarding request to: %s\n", targetServer.String())
	proxy.ServeHTTP(w, r)
}

func main() {
	servers := []string{
		"http://localhost:5001",
		"http://localhost:5002",
		"http://localhost:5003",
	}

	lb := NewLoadBalancer(servers)

	fmt.Println("Load Balancer started at :8080")
	http.ListenAndServe(":8080", lb)
}
