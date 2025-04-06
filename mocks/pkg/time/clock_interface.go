// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// ClockInterface is an autogenerated mock type for the ClockInterface type
type ClockInterface struct {
	mock.Mock
}

type ClockInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *ClockInterface) EXPECT() *ClockInterface_Expecter {
	return &ClockInterface_Expecter{mock: &_m.Mock}
}

// Now provides a mock function with no fields
func (_m *ClockInterface) Now() time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Now")
	}

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// ClockInterface_Now_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Now'
type ClockInterface_Now_Call struct {
	*mock.Call
}

// Now is a helper method to define mock.On call
func (_e *ClockInterface_Expecter) Now() *ClockInterface_Now_Call {
	return &ClockInterface_Now_Call{Call: _e.mock.On("Now")}
}

func (_c *ClockInterface_Now_Call) Run(run func()) *ClockInterface_Now_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ClockInterface_Now_Call) Return(_a0 time.Time) *ClockInterface_Now_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ClockInterface_Now_Call) RunAndReturn(run func() time.Time) *ClockInterface_Now_Call {
	_c.Call.Return(run)
	return _c
}

// NewClockInterface creates a new instance of ClockInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClockInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ClockInterface {
	mock := &ClockInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
