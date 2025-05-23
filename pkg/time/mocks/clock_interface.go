// Code generated by mockery v2.53.3. DO NOT EDIT.

package time

import (
	mock "github.com/stretchr/testify/mock"

	time2 "time"
)

// MockClockInterface is an autogenerated mock type for the ClockInterface type
type MockClockInterface struct {
	mock.Mock
}

type MockClockInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockClockInterface) EXPECT() *MockClockInterface_Expecter {
	return &MockClockInterface_Expecter{mock: &_m.Mock}
}

// Now provides a mock function with no fields
func (_m *MockClockInterface) Now() time2.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Now")
	}

	var r0 time2.Time
	if rf, ok := ret.Get(0).(func() time2.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time2.Time)
	}

	return r0
}

// MockClockInterface_Now_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Now'
type MockClockInterface_Now_Call struct {
	*mock.Call
}

// Now is a helper method to define mock.On call
func (_e *MockClockInterface_Expecter) Now() *MockClockInterface_Now_Call {
	return &MockClockInterface_Now_Call{Call: _e.mock.On("Now")}
}

func (_c *MockClockInterface_Now_Call) Run(run func()) *MockClockInterface_Now_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockClockInterface_Now_Call) Return(_a0 time2.Time) *MockClockInterface_Now_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockClockInterface_Now_Call) RunAndReturn(run func() time2.Time) *MockClockInterface_Now_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockClockInterface creates a new instance of MockClockInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockClockInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockClockInterface {
	mock := &MockClockInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
