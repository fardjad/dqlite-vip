// Code generated by mockery v2.50.1. DO NOT EDIT.

package vip

import (
	cluster "fardjad.com/dqlite-vip/cluster"
	mock "github.com/stretchr/testify/mock"

	utils "fardjad.com/dqlite-vip/utils"

	vip "fardjad.com/dqlite-vip/vip"
)

// ManagerFactoryFunc is an autogenerated mock type for the ManagerFactoryFunc type
type ManagerFactoryFunc struct {
	mock.Mock
}

type ManagerFactoryFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *ManagerFactoryFunc) EXPECT() *ManagerFactoryFunc_Expecter {
	return &ManagerFactoryFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *ManagerFactoryFunc) Execute(_a0 cluster.ClusterNode, _a1 utils.BetterTicker, _a2 vip.Configurator, _a3 string) vip.Manager {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 vip.Manager
	if rf, ok := ret.Get(0).(func(cluster.ClusterNode, utils.BetterTicker, vip.Configurator, string) vip.Manager); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(vip.Manager)
		}
	}

	return r0
}

// ManagerFactoryFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type ManagerFactoryFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 cluster.ClusterNode
//   - _a1 utils.BetterTicker
//   - _a2 vip.Configurator
//   - _a3 string
func (_e *ManagerFactoryFunc_Expecter) Execute(_a0 interface{}, _a1 interface{}, _a2 interface{}, _a3 interface{}) *ManagerFactoryFunc_Execute_Call {
	return &ManagerFactoryFunc_Execute_Call{Call: _e.mock.On("Execute", _a0, _a1, _a2, _a3)}
}

func (_c *ManagerFactoryFunc_Execute_Call) Run(run func(_a0 cluster.ClusterNode, _a1 utils.BetterTicker, _a2 vip.Configurator, _a3 string)) *ManagerFactoryFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(cluster.ClusterNode), args[1].(utils.BetterTicker), args[2].(vip.Configurator), args[3].(string))
	})
	return _c
}

func (_c *ManagerFactoryFunc_Execute_Call) Return(_a0 vip.Manager) *ManagerFactoryFunc_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManagerFactoryFunc_Execute_Call) RunAndReturn(run func(cluster.ClusterNode, utils.BetterTicker, vip.Configurator, string) vip.Manager) *ManagerFactoryFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewManagerFactoryFunc creates a new instance of ManagerFactoryFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewManagerFactoryFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *ManagerFactoryFunc {
	mock := &ManagerFactoryFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}