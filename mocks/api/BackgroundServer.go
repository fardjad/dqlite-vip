// Code generated by mockery v2.50.4. DO NOT EDIT.

package api

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// BackgroundServer is an autogenerated mock type for the BackgroundServer type
type BackgroundServer struct {
	mock.Mock
}

type BackgroundServer_Expecter struct {
	mock *mock.Mock
}

func (_m *BackgroundServer) EXPECT() *BackgroundServer_Expecter {
	return &BackgroundServer_Expecter{mock: &_m.Mock}
}

// ListenAndServeInBackground provides a mock function with no fields
func (_m *BackgroundServer) ListenAndServeInBackground() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ListenAndServeInBackground")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BackgroundServer_ListenAndServeInBackground_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListenAndServeInBackground'
type BackgroundServer_ListenAndServeInBackground_Call struct {
	*mock.Call
}

// ListenAndServeInBackground is a helper method to define mock.On call
func (_e *BackgroundServer_Expecter) ListenAndServeInBackground() *BackgroundServer_ListenAndServeInBackground_Call {
	return &BackgroundServer_ListenAndServeInBackground_Call{Call: _e.mock.On("ListenAndServeInBackground")}
}

func (_c *BackgroundServer_ListenAndServeInBackground_Call) Run(run func()) *BackgroundServer_ListenAndServeInBackground_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *BackgroundServer_ListenAndServeInBackground_Call) Return(_a0 error) *BackgroundServer_ListenAndServeInBackground_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BackgroundServer_ListenAndServeInBackground_Call) RunAndReturn(run func() error) *BackgroundServer_ListenAndServeInBackground_Call {
	_c.Call.Return(run)
	return _c
}

// Shutdown provides a mock function with given fields: ctx
func (_m *BackgroundServer) Shutdown(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Shutdown")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BackgroundServer_Shutdown_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Shutdown'
type BackgroundServer_Shutdown_Call struct {
	*mock.Call
}

// Shutdown is a helper method to define mock.On call
//   - ctx context.Context
func (_e *BackgroundServer_Expecter) Shutdown(ctx interface{}) *BackgroundServer_Shutdown_Call {
	return &BackgroundServer_Shutdown_Call{Call: _e.mock.On("Shutdown", ctx)}
}

func (_c *BackgroundServer_Shutdown_Call) Run(run func(ctx context.Context)) *BackgroundServer_Shutdown_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *BackgroundServer_Shutdown_Call) Return(_a0 error) *BackgroundServer_Shutdown_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BackgroundServer_Shutdown_Call) RunAndReturn(run func(context.Context) error) *BackgroundServer_Shutdown_Call {
	_c.Call.Return(run)
	return _c
}

// NewBackgroundServer creates a new instance of BackgroundServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBackgroundServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *BackgroundServer {
	mock := &BackgroundServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
