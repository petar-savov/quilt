package balancer

import (
	"errors"
	"sync"
)

type Config struct {
	ListenAddr      string
	UpstreamServers []string
}

type Balancer interface {
	NextUpstream() (string, error)
}

type roundRobinBalancer struct {
	servers []string
	index   int
	mu      sync.Mutex
}

func New(config *Config) (Balancer, error) {
	if len(config.UpstreamServers) == 0 {
		return nil, errors.New("no upstream servers provided")
	}

	return &roundRobinBalancer{
		servers: config.UpstreamServers,
	}, nil
}

func (r *roundRobinBalancer) NextUpstream() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	upstream := r.servers[r.index]
	r.index = (r.index + 1) % len(r.servers)

	return upstream, nil
}
