package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	clusterMocks "fardjad.com/dqlite-vip/mocks/cluster"
)

type GetVIPTestSuite struct {
	suite.Suite

	clusterNode *clusterMocks.ClusterNode
	mux         *http.ServeMux
}

func (s *GetVIPTestSuite) SetupTest() {
	s.clusterNode = clusterMocks.NewClusterNode(s.T())
	handlers := NewHandlers(s.clusterNode)
	s.mux = handlers.Mux()
}

func (s *GetVIPTestSuite) TestGetVIP() {
	s.clusterNode.EXPECT().GetString(mock.Anything, "vip").Return("192.168.1.100", nil)

	request, _ := http.NewRequest(http.MethodGet, "/vip", nil)
	response := httptest.NewRecorder()
	s.mux.ServeHTTP(response, request)

	s.Equal(http.StatusOK, response.Code)
	s.Equal("application/json", response.Header().Get("Content-Type"))

	var responseBody GetVIPResponseBody
	json.NewDecoder(response.Body).Decode(&responseBody)
	s.Equal(&GetVIPResponseBody{
		VIP: "192.168.1.100",
	}, &responseBody)
}

func TestGetVIPTestSuite(t *testing.T) {
	suite.Run(t, new(GetVIPTestSuite))
}
