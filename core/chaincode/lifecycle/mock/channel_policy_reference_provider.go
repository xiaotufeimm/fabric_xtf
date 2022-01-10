// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"sync"

	"github.com/hyperledger/fabric/v2/common/policies"
	"github.com/hyperledger/fabric/v2/core/chaincode/lifecycle"
)

type ChannelPolicyReferenceProvider struct {
	NewPolicyStub        func(string, string) (policies.Policy, error)
	newPolicyMutex       sync.RWMutex
	newPolicyArgsForCall []struct {
		arg1 string
		arg2 string
	}
	newPolicyReturns struct {
		result1 policies.Policy
		result2 error
	}
	newPolicyReturnsOnCall map[int]struct {
		result1 policies.Policy
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChannelPolicyReferenceProvider) NewPolicy(arg1 string, arg2 string) (policies.Policy, error) {
	fake.newPolicyMutex.Lock()
	ret, specificReturn := fake.newPolicyReturnsOnCall[len(fake.newPolicyArgsForCall)]
	fake.newPolicyArgsForCall = append(fake.newPolicyArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("NewPolicy", []interface{}{arg1, arg2})
	fake.newPolicyMutex.Unlock()
	if fake.NewPolicyStub != nil {
		return fake.NewPolicyStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.newPolicyReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChannelPolicyReferenceProvider) NewPolicyCallCount() int {
	fake.newPolicyMutex.RLock()
	defer fake.newPolicyMutex.RUnlock()
	return len(fake.newPolicyArgsForCall)
}

func (fake *ChannelPolicyReferenceProvider) NewPolicyCalls(stub func(string, string) (policies.Policy, error)) {
	fake.newPolicyMutex.Lock()
	defer fake.newPolicyMutex.Unlock()
	fake.NewPolicyStub = stub
}

func (fake *ChannelPolicyReferenceProvider) NewPolicyArgsForCall(i int) (string, string) {
	fake.newPolicyMutex.RLock()
	defer fake.newPolicyMutex.RUnlock()
	argsForCall := fake.newPolicyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ChannelPolicyReferenceProvider) NewPolicyReturns(result1 policies.Policy, result2 error) {
	fake.newPolicyMutex.Lock()
	defer fake.newPolicyMutex.Unlock()
	fake.NewPolicyStub = nil
	fake.newPolicyReturns = struct {
		result1 policies.Policy
		result2 error
	}{result1, result2}
}

func (fake *ChannelPolicyReferenceProvider) NewPolicyReturnsOnCall(i int, result1 policies.Policy, result2 error) {
	fake.newPolicyMutex.Lock()
	defer fake.newPolicyMutex.Unlock()
	fake.NewPolicyStub = nil
	if fake.newPolicyReturnsOnCall == nil {
		fake.newPolicyReturnsOnCall = make(map[int]struct {
			result1 policies.Policy
			result2 error
		})
	}
	fake.newPolicyReturnsOnCall[i] = struct {
		result1 policies.Policy
		result2 error
	}{result1, result2}
}

func (fake *ChannelPolicyReferenceProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newPolicyMutex.RLock()
	defer fake.newPolicyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChannelPolicyReferenceProvider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ lifecycle.ChannelPolicyReferenceProvider = new(ChannelPolicyReferenceProvider)
