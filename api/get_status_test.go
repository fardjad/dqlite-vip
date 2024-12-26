package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"fardjad.com/dqlite-vip/cluster"
	clusterMocks "fardjad.com/dqlite-vip/mocks/cluster"
)

type GetStatusTestSuite struct {
	suite.Suite

	clusterNode *clusterMocks.ClusterNode
	mux         *http.ServeMux
}

func (s *GetStatusTestSuite) SetupTest() {
	s.clusterNode = clusterMocks.NewClusterNode(s.T())
	handlers := NewHandlers(s.clusterNode)
	s.mux = handlers.Mux()
}

func (s *GetStatusTestSuite) TestGetStatus_Healthy() {
	clusterMembers := []*cluster.ClusterMemberInfo{
		{ID: 1, Address: "192.168.1.1", Role: "voter"},
		{ID: 2, Address: "192.168.1.2", Role: "voter"},
		{ID: 3, Address: "192.168.1.3", Role: "voter"},
	}
	s.clusterNode.EXPECT().ID().Return(uint64(1))
	s.clusterNode.EXPECT().LeaderID(mock.Anything).Return(uint64(1), nil)
	s.clusterNode.EXPECT().IsLeader(mock.Anything).Return(true)
	s.clusterNode.EXPECT().ClusterMembers(mock.Anything).Return(clusterMembers, nil)

	request, _ := http.NewRequest(http.MethodGet, "/status", nil)
	response := httptest.NewRecorder()
	s.mux.ServeHTTP(response, request)

	s.clusterNode.AssertExpectations(s.T())

	s.Equal(http.StatusOK, response.Code)
	s.Equal("application/json", response.Header().Get("Content-Type"))

	var responseBody GetStatusResponseBody
	json.NewDecoder(response.Body).Decode(&responseBody)
	s.Equal(&GetStatusResponseBody{
		ID:             1,
		LeaderID:       1,
		IsLeader:       true,
		ClusterMembers: clusterMembers,
	}, &responseBody)
}

func TestGetStatusTestSuite(t *testing.T) {
	suite.Run(t, new(GetStatusTestSuite))
}
