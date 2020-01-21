// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"context"
	"sync"

	"github.com/hyperledger/fabric/common/deliver"
)

type Semaphore struct {
	AcquireStub        func(context.Context) error
	acquireMutex       sync.RWMutex
	acquireArgsForCall []struct {
		arg1 context.Context
	}
	acquireReturns struct {
		result1 error
	}
	acquireReturnsOnCall map[int]struct {
		result1 error
	}
	ReleaseStub        func()
	releaseMutex       sync.RWMutex
	releaseArgsForCall []struct {
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Semaphore) Acquire(arg1 context.Context) error {
	fake.acquireMutex.Lock()
	ret, specificReturn := fake.acquireReturnsOnCall[len(fake.acquireArgsForCall)]
	fake.acquireArgsForCall = append(fake.acquireArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	fake.recordInvocation("Acquire", []interface{}{arg1})
	fake.acquireMutex.Unlock()
	if fake.AcquireStub != nil {
		return fake.AcquireStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.acquireReturns
	return fakeReturns.result1
}

func (fake *Semaphore) AcquireCallCount() int {
	fake.acquireMutex.RLock()
	defer fake.acquireMutex.RUnlock()
	return len(fake.acquireArgsForCall)
}

func (fake *Semaphore) AcquireCalls(stub func(context.Context) error) {
	fake.acquireMutex.Lock()
	defer fake.acquireMutex.Unlock()
	fake.AcquireStub = stub
}

func (fake *Semaphore) AcquireArgsForCall(i int) context.Context {
	fake.acquireMutex.RLock()
	defer fake.acquireMutex.RUnlock()
	argsForCall := fake.acquireArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Semaphore) AcquireReturns(result1 error) {
	fake.acquireMutex.Lock()
	defer fake.acquireMutex.Unlock()
	fake.AcquireStub = nil
	fake.acquireReturns = struct {
		result1 error
	}{result1}
}

func (fake *Semaphore) AcquireReturnsOnCall(i int, result1 error) {
	fake.acquireMutex.Lock()
	defer fake.acquireMutex.Unlock()
	fake.AcquireStub = nil
	if fake.acquireReturnsOnCall == nil {
		fake.acquireReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.acquireReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Semaphore) Release() {
	fake.releaseMutex.Lock()
	fake.releaseArgsForCall = append(fake.releaseArgsForCall, struct {
	}{})
	fake.recordInvocation("Release", []interface{}{})
	fake.releaseMutex.Unlock()
	if fake.ReleaseStub != nil {
		fake.ReleaseStub()
	}
}

func (fake *Semaphore) ReleaseCallCount() int {
	fake.releaseMutex.RLock()
	defer fake.releaseMutex.RUnlock()
	return len(fake.releaseArgsForCall)
}

func (fake *Semaphore) ReleaseCalls(stub func()) {
	fake.releaseMutex.Lock()
	defer fake.releaseMutex.Unlock()
	fake.ReleaseStub = stub
}

func (fake *Semaphore) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.acquireMutex.RLock()
	defer fake.acquireMutex.RUnlock()
	fake.releaseMutex.RLock()
	defer fake.releaseMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Semaphore) recordInvocation(key string, args []interface{}) {
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

var _ deliver.Semaphore = new(Semaphore)
