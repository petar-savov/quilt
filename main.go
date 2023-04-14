package main

import (
	"io"
	"log"
	"net/http"

	"github.com/petar-savov/quilt/pkg/balancer"
)

func main() {
	config := &balancer.Config{
		ListenAddr: ":8080",
		UpstreamServers: []string{
			"http://localhost:8001",
			"http://localhost:8002",
		},
	}

	quiltBalancer, err := balancer.New(config)
	if err != nil {
		log.Fatalf("Error initializing load balancer: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		upstream, err := quiltBalancer.NextUpstream()
		if err != nil {
			http.Error(w, "Failed to select an upstream server", http.StatusInternalServerError)
			return
		}

		client := http.Client{}

		r.URL.Scheme = upstream.Scheme
		r.URL.Host = upstream.Host

		resp, err := client.Do(r)
		if err != nil {
			http.Error(w, "Failed to forward request", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		for name, values := range resp.Header {
			w.Header()[name] = values
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	log.Printf("Starting load balancer on %s", config.ListenAddr)
	if err := http.ListenAndServe(config.ListenAddr, nil); err != nil {
		log.Fatalf("Error starting load balancer: %v", err)
	}

}
