// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	context "context"
	io "io"

	mock "github.com/stretchr/testify/mock"

	pipe "github.com/antonmisa/1cctl/pkg/pipe"
)

// Piper is an autogenerated mock type for the Piper type
type Piper struct {
	mock.Mock
}

// Run provides a mock function with given fields: ctx, arg
func (_m *Piper) Run(ctx context.Context, arg ...string) (pipe.Commander, io.ReadCloser, error) {
	_va := make([]interface{}, len(arg))
	for _i := range arg {
		_va[_i] = arg[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 pipe.Commander
	var r1 io.ReadCloser
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, ...string) (pipe.Commander, io.ReadCloser, error)); ok {
		return rf(ctx, arg...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...string) pipe.Commander); ok {
		r0 = rf(ctx, arg...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pipe.Commander)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...string) io.ReadCloser); ok {
		r1 = rf(ctx, arg...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(io.ReadCloser)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, ...string) error); ok {
		r2 = rf(ctx, arg...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewPiper creates a new instance of Piper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPiper(t interface {
	mock.TestingT
	Cleanup(func())
}) *Piper {
	mock := &Piper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}