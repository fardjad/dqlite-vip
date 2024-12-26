package api

import (
	"net/http"

	"fardjad.com/dqlite-vip/cluster"
)

type Handlers struct {
	clusterNode cluster.ClusterNode
}

func NewHandlers(clusterNode cluster.ClusterNode) Handlers {
	return Handlers{clusterNode: clusterNode}
}

func (s *Handlers) GetHealth(w http.ResponseWriter, r *http.Request) {
}

func (s *Handlers) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.GetHealth)

	return mux
}
