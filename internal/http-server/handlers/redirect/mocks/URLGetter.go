// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// URLGetter is an autogenerated mock type for the URLGetter type
type URLGetter struct {
	mock.Mock
}

type URLGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *URLGetter) EXPECT() *URLGetter_Expecter {
	return &URLGetter_Expecter{mock: &_m.Mock}
}

// GetURL provides a mock function with given fields: alias
func (_m *URLGetter) GetURL(alias string) (string, error) {
	ret := _m.Called(alias)

	if len(ret) == 0 {
		panic("no return value specified for GetURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(alias)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// URLGetter_GetURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetURL'
type URLGetter_GetURL_Call struct {
	*mock.Call
}

// GetURL is a helper method to define mock.On call
//   - alias string
func (_e *URLGetter_Expecter) GetURL(alias interface{}) *URLGetter_GetURL_Call {
	return &URLGetter_GetURL_Call{Call: _e.mock.On("GetURL", alias)}
}

func (_c *URLGetter_GetURL_Call) Run(run func(alias string)) *URLGetter_GetURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *URLGetter_GetURL_Call) Return(_a0 string, _a1 error) *URLGetter_GetURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *URLGetter_GetURL_Call) RunAndReturn(run func(string) (string, error)) *URLGetter_GetURL_Call {
	_c.Call.Return(run)
	return _c
}

// NewURLGetter creates a new instance of URLGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewURLGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *URLGetter {
	mock := &URLGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
