// Code generated by mockery v2.50.1. DO NOT EDIT.

package mocks

import (
	cluster "fardjad.com/dqlite-vip/cluster"
	mock "github.com/stretchr/testify/mock"
)

// ClusterNodeFactory is an autogenerated mock type for the ClusterNodeFactory type
type ClusterNodeFactory struct {
	mock.Mock
}

type ClusterNodeFactory_Expecter struct {
	mock *mock.Mock
}

func (_m *ClusterNodeFactory) EXPECT() *ClusterNodeFactory_Expecter {
	return &ClusterNodeFactory_Expecter{mock: &_m.Mock}
}

// NewClusterNode provides a mock function with given fields: dataDir, bindCluster, bindHttp, join
func (_m *ClusterNodeFactory) NewClusterNode(dataDir string, bindCluster string, bindHttp string, join []string) (cluster.ClusterNode, error) {
	ret := _m.Called(dataDir, bindCluster, bindHttp, join)

	if len(ret) == 0 {
		panic("no return value specified for NewClusterNode")
	}

	var r0 cluster.ClusterNode
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, string, []string) (cluster.ClusterNode, error)); ok {
		return rf(dataDir, bindCluster, bindHttp, join)
	}
	if rf, ok := ret.Get(0).(func(string, string, string, []string) cluster.ClusterNode); ok {
		r0 = rf(dataDir, bindCluster, bindHttp, join)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cluster.ClusterNode)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, string, []string) error); ok {
		r1 = rf(dataDir, bindCluster, bindHttp, join)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClusterNodeFactory_NewClusterNode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewClusterNode'
type ClusterNodeFactory_NewClusterNode_Call struct {
	*mock.Call
}

// NewClusterNode is a helper method to define mock.On call
//   - dataDir string
//   - bindCluster string
//   - bindHttp string
//   - join []string
func (_e *ClusterNodeFactory_Expecter) NewClusterNode(dataDir interface{}, bindCluster interface{}, bindHttp interface{}, join interface{}) *ClusterNodeFactory_NewClusterNode_Call {
	return &ClusterNodeFactory_NewClusterNode_Call{Call: _e.mock.On("NewClusterNode", dataDir, bindCluster, bindHttp, join)}
}

func (_c *ClusterNodeFactory_NewClusterNode_Call) Run(run func(dataDir string, bindCluster string, bindHttp string, join []string)) *ClusterNodeFactory_NewClusterNode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string), args[3].([]string))
	})
	return _c
}

func (_c *ClusterNodeFactory_NewClusterNode_Call) Return(_a0 cluster.ClusterNode, _a1 error) *ClusterNodeFactory_NewClusterNode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ClusterNodeFactory_NewClusterNode_Call) RunAndReturn(run func(string, string, string, []string) (cluster.ClusterNode, error)) *ClusterNodeFactory_NewClusterNode_Call {
	_c.Call.Return(run)
	return _c
}

// NewClusterNodeFactory creates a new instance of ClusterNodeFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClusterNodeFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *ClusterNodeFactory {
	mock := &ClusterNodeFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
