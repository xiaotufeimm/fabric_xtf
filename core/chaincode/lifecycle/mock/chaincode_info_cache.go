// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"sync"

	"github.com/hyperledger/fabric/v2/core/chaincode/lifecycle"
)

type ChaincodeInfoCache struct {
	ChaincodeInfoStub        func(string, string) (*lifecycle.LocalChaincodeInfo, error)
	chaincodeInfoMutex       sync.RWMutex
	chaincodeInfoArgsForCall []struct {
		arg1 string
		arg2 string
	}
	chaincodeInfoReturns struct {
		result1 *lifecycle.LocalChaincodeInfo
		result2 error
	}
	chaincodeInfoReturnsOnCall map[int]struct {
		result1 *lifecycle.LocalChaincodeInfo
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChaincodeInfoCache) ChaincodeInfo(arg1 string, arg2 string) (*lifecycle.LocalChaincodeInfo, error) {
	fake.chaincodeInfoMutex.Lock()
	ret, specificReturn := fake.chaincodeInfoReturnsOnCall[len(fake.chaincodeInfoArgsForCall)]
	fake.chaincodeInfoArgsForCall = append(fake.chaincodeInfoArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ChaincodeInfo", []interface{}{arg1, arg2})
	fake.chaincodeInfoMutex.Unlock()
	if fake.ChaincodeInfoStub != nil {
		return fake.ChaincodeInfoStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.chaincodeInfoReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeInfoCache) ChaincodeInfoCallCount() int {
	fake.chaincodeInfoMutex.RLock()
	defer fake.chaincodeInfoMutex.RUnlock()
	return len(fake.chaincodeInfoArgsForCall)
}

func (fake *ChaincodeInfoCache) ChaincodeInfoCalls(stub func(string, string) (*lifecycle.LocalChaincodeInfo, error)) {
	fake.chaincodeInfoMutex.Lock()
	defer fake.chaincodeInfoMutex.Unlock()
	fake.ChaincodeInfoStub = stub
}

func (fake *ChaincodeInfoCache) ChaincodeInfoArgsForCall(i int) (string, string) {
	fake.chaincodeInfoMutex.RLock()
	defer fake.chaincodeInfoMutex.RUnlock()
	argsForCall := fake.chaincodeInfoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ChaincodeInfoCache) ChaincodeInfoReturns(result1 *lifecycle.LocalChaincodeInfo, result2 error) {
	fake.chaincodeInfoMutex.Lock()
	defer fake.chaincodeInfoMutex.Unlock()
	fake.ChaincodeInfoStub = nil
	fake.chaincodeInfoReturns = struct {
		result1 *lifecycle.LocalChaincodeInfo
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeInfoCache) ChaincodeInfoReturnsOnCall(i int, result1 *lifecycle.LocalChaincodeInfo, result2 error) {
	fake.chaincodeInfoMutex.Lock()
	defer fake.chaincodeInfoMutex.Unlock()
	fake.ChaincodeInfoStub = nil
	if fake.chaincodeInfoReturnsOnCall == nil {
		fake.chaincodeInfoReturnsOnCall = make(map[int]struct {
			result1 *lifecycle.LocalChaincodeInfo
			result2 error
		})
	}
	fake.chaincodeInfoReturnsOnCall[i] = struct {
		result1 *lifecycle.LocalChaincodeInfo
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeInfoCache) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.chaincodeInfoMutex.RLock()
	defer fake.chaincodeInfoMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChaincodeInfoCache) recordInvocation(key string, args []interface{}) {
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

var _ lifecycle.ChaincodeInfoCache = new(ChaincodeInfoCache)
