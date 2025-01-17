// Code generated by mockery v2.50.4. DO NOT EDIT.

package vip

import mock "github.com/stretchr/testify/mock"

// Configurator is an autogenerated mock type for the Configurator type
type Configurator struct {
	mock.Mock
}

type Configurator_Expecter struct {
	mock *mock.Mock
}

func (_m *Configurator) EXPECT() *Configurator_Expecter {
	return &Configurator_Expecter{mock: &_m.Mock}
}

// EnsureVIP provides a mock function with given fields: iface, address
func (_m *Configurator) EnsureVIP(iface string, address string) error {
	ret := _m.Called(iface, address)

	if len(ret) == 0 {
		panic("no return value specified for EnsureVIP")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(iface, address)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Configurator_EnsureVIP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EnsureVIP'
type Configurator_EnsureVIP_Call struct {
	*mock.Call
}

// EnsureVIP is a helper method to define mock.On call
//   - iface string
//   - address string
func (_e *Configurator_Expecter) EnsureVIP(iface interface{}, address interface{}) *Configurator_EnsureVIP_Call {
	return &Configurator_EnsureVIP_Call{Call: _e.mock.On("EnsureVIP", iface, address)}
}

func (_c *Configurator_EnsureVIP_Call) Run(run func(iface string, address string)) *Configurator_EnsureVIP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *Configurator_EnsureVIP_Call) Return(_a0 error) *Configurator_EnsureVIP_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Configurator_EnsureVIP_Call) RunAndReturn(run func(string, string) error) *Configurator_EnsureVIP_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveVIP provides a mock function with given fields: iface, address
func (_m *Configurator) RemoveVIP(iface string, address string) error {
	ret := _m.Called(iface, address)

	if len(ret) == 0 {
		panic("no return value specified for RemoveVIP")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(iface, address)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Configurator_RemoveVIP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveVIP'
type Configurator_RemoveVIP_Call struct {
	*mock.Call
}

// RemoveVIP is a helper method to define mock.On call
//   - iface string
//   - address string
func (_e *Configurator_Expecter) RemoveVIP(iface interface{}, address interface{}) *Configurator_RemoveVIP_Call {
	return &Configurator_RemoveVIP_Call{Call: _e.mock.On("RemoveVIP", iface, address)}
}

func (_c *Configurator_RemoveVIP_Call) Run(run func(iface string, address string)) *Configurator_RemoveVIP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *Configurator_RemoveVIP_Call) Return(_a0 error) *Configurator_RemoveVIP_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Configurator_RemoveVIP_Call) RunAndReturn(run func(string, string) error) *Configurator_RemoveVIP_Call {
	_c.Call.Return(run)
	return _c
}

// NewConfigurator creates a new instance of Configurator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConfigurator(t interface {
	mock.TestingT
	Cleanup(func())
}) *Configurator {
	mock := &Configurator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
