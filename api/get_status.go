package api

import (
	"net/http"

	"fardjad.com/dqlite-vip/cluster"
)

type GetStatusResponseBody struct {
	ID             uint64                       `json:"id"`
	LeaderID       uint64                       `json:"leader_id"`
	IsLeader       bool                         `json:"is_leader"`
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

	isLeader := s.clusterNode.IsLeader(r.Context())

	response := GetStatusResponseBody{
		ID:             id,
		LeaderID:       leaderID,
		IsLeader:       isLeader,
		ClusterMembers: clusterMembers,
	}

	s.writeJSON(w, http.StatusOK, nil, response)
}
