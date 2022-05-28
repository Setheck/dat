// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	pflag "github.com/spf13/pflag"
	mock "github.com/stretchr/testify/mock"
)

// CobraCommand is an autogenerated mock type for the CobraCommand type
type CobraCommand struct {
	mock.Mock
}

// Execute provides a mock function with given fields:
func (_m *CobraCommand) Execute() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Flags provides a mock function with given fields:
func (_m *CobraCommand) Flags() *pflag.FlagSet {
	ret := _m.Called()

	var r0 *pflag.FlagSet
	if rf, ok := ret.Get(0).(func() *pflag.FlagSet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pflag.FlagSet)
		}
	}

	return r0
}