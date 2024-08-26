// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	metadata "google.golang.org/grpc/metadata"

	orderer "github.com/hyperledger/fabric-protos-go-apiv2/orderer"
)

// StepStream is an autogenerated mock type for the StepStream type
type StepStream struct {
	mock.Mock
}

// Context provides a mock function with given fields:
func (_m *StepStream) Context() context.Context {
	ret := _m.Called()

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// Recv provides a mock function with given fields:
func (_m *StepStream) Recv() (*orderer.StepRequest, error) {
	ret := _m.Called()

	var r0 *orderer.StepRequest
	if rf, ok := ret.Get(0).(func() *orderer.StepRequest); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*orderer.StepRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RecvMsg provides a mock function with given fields: m
func (_m *StepStream) RecvMsg(m interface{}) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Send provides a mock function with given fields: response
func (_m *StepStream) Send(response *orderer.StepResponse) error {
	ret := _m.Called(response)

	var r0 error
	if rf, ok := ret.Get(0).(func(*orderer.StepResponse) error); ok {
		r0 = rf(response)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendHeader provides a mock function with given fields: _a0
func (_m *StepStream) SendHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendMsg provides a mock function with given fields: m
func (_m *StepStream) SendMsg(m interface{}) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetHeader provides a mock function with given fields: _a0
func (_m *StepStream) SetHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetTrailer provides a mock function with given fields: _a0
func (_m *StepStream) SetTrailer(_a0 metadata.MD) {
	_m.Called(_a0)
}

type mockConstructorTestingTNewStepStream interface {
	mock.TestingT
	Cleanup(func())
}

// NewStepStream creates a new instance of StepStream. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStepStream(t mockConstructorTestingTNewStepStream) *StepStream {
	mock := &StepStream{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
