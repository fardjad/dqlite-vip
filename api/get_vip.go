package api

import (
	"net/http"
)

type GetVIPResponseBody struct {
	VIP string `json:"vip"`
}

func (s *Handlers) GetVIP(w http.ResponseWriter, r *http.Request) {
	vip, err := s.clusterNode.GetString(r.Context(), "vip")
	if err != nil {
		s.writeErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	s.writeJSON(w, http.StatusOK, nil, GetVIPResponseBody{VIP: vip})
}
