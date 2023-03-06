package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type RoundRobinBalancer struct {
	targets []*url.URL
	current int
	lock    sync.Mutex
}

func (r *RoundRobinBalancer) NextTarget() *url.URL {
	r.lock.Lock()
	defer r.lock.Unlock()
	target := r.targets[r.current]
	r.current = (r.current + 1) % len(r.targets)
	return target
}

func main() {
	// Crea un slice de URL de destino que se balancearán
	targetUrls := []*url.URL{
		{Scheme: "http", Host: "localhost:8081"},
		{Scheme: "http", Host: "localhost:8082"},
		{Scheme: "http", Host: "localhost:8083"},
	}

	// Crea un balanceador de carga Round Robin
	balancer := &RoundRobinBalancer{
		targets: targetUrls,
		current: 0,
	}
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			targetUrl := balancer.NextTarget()
			req.URL.Scheme = targetUrl.Scheme
			req.URL.Host = targetUrl.Host
		},
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
		},
	}

	// Configura un servidor HTTP para escuchar y manejar solicitudes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
	fmt.Println("Listening on port 8888...")
	http.ListenAndServe(":8080", nil)
}
