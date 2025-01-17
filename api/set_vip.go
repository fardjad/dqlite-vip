package api

import (
	"encoding/json"
	"net/http"
)

type SetVIPResponseBody struct {
	Message string `json:"message"`
}

type SetVIPRequestBody struct {
	VIP string `json:"vip"`
}

func (s *Handlers) SetVIP(w http.ResponseWriter, r *http.Request) {
	var requestBody SetVIPRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		s.writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	err = s.clusterNode.SetString(r.Context(), "vip", requestBody.VIP)
	if err != nil {
		s.writeErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	s.writeJSON(w, http.StatusAccepted, nil, SetVIPResponseBody{Message: "OK"})
}
