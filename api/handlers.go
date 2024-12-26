package api

import (
	"encoding/json"
	"net/http"

	"fardjad.com/dqlite-vip/cluster"
)

type Handlers struct {
	clusterNode cluster.ClusterNode
}

func NewHandlers(clusterNode cluster.ClusterNode) Handlers {
	return Handlers{clusterNode: clusterNode}
}

func (s *Handlers) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /status", s.GetStatus)
	mux.HandleFunc("PUT /vip", s.SetVIP)

	return mux
}

func (s *Handlers) writeJSON(w http.ResponseWriter, statusCode int, headers map[string]string, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	for key, value := range headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
