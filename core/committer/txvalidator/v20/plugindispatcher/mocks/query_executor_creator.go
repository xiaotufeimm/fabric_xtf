// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	ledger "github.com/hyperledger/fabric/v2/core/ledger"
	mock "github.com/stretchr/testify/mock"
)

// QueryExecutorCreator is an autogenerated mock type for the QueryExecutorCreator type
type QueryExecutorCreator struct {
	mock.Mock
}

// NewQueryExecutor provides a mock function with given fields:
func (_m *QueryExecutorCreator) NewQueryExecutor() (ledger.QueryExecutor, error) {
	ret := _m.Called()

	var r0 ledger.QueryExecutor
	if rf, ok := ret.Get(0).(func() ledger.QueryExecutor); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ledger.QueryExecutor)
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
