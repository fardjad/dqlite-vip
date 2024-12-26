package api

import (
	"net/http"

	"fardjad.com/dqlite-vip/cluster"
)

type GetStatusResponseBody struct {
	ID             uint64                       `json:"id"`
	LeaderID       uint64                       `json:"leader_id"`
	ClusterMembers []*cluster.ClusterMemberInfo `json:"cluster_members"`
}

func (s *Handlers) GetStatus(w http.ResponseWriter, r *http.Request) {
	clusterMembers, err := s.clusterNode.ClusterMembers(r.Context())
	if err != nil {
		s.writeErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	id := s.clusterNode.ID()
	leaderID, err := s.clusterNode.LeaderID(r.Context())
	if err != nil {
		s.writeErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	response := GetStatusResponseBody{
		ID:             id,
		LeaderID:       leaderID,
		ClusterMembers: clusterMembers,
	}

	s.writeJSON(w, http.StatusOK, nil, response)
}
