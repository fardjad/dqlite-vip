package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	clusterMocks "fardjad.com/dqlite-vip/mocks/cluster"
)

type SetVIPTestSuite struct {
	suite.Suite

	clusterNode *clusterMocks.ClusterNode
	mux         *http.ServeMux
}

func (s *SetVIPTestSuite) SetupTest() {
	s.clusterNode = clusterMocks.NewClusterNode(s.T())
	handlers := NewHandlers(s.clusterNode)
	s.mux = handlers.Mux()
}

func (s *SetVIPTestSuite) TestSetVIP() {
	s.clusterNode.EXPECT().SetString(mock.Anything, "vip", "192.168.1.100").Return(nil)

	requestBody, _ := json.Marshal(&SetVIPRequestBody{
		VIP: "192.168.1.100",
	})
	request, _ := http.NewRequest(http.MethodPut, "/vip", bytes.NewReader(requestBody))
	response := httptest.NewRecorder()
	s.mux.ServeHTTP(response, request)

	s.Equal(http.StatusAccepted, response.Code)
	s.Equal("application/json", response.Header().Get("Content-Type"))

	var responseBody SetVIPResponseBody
	json.NewDecoder(response.Body).Decode(&responseBody)
	s.Equal(&SetVIPResponseBody{
		Message: "OK",
	}, &responseBody)
}

func TestSetVIPTestSuite(t *testing.T) {
	suite.Run(t, new(SetVIPTestSuite))
}
