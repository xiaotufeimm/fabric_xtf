// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"sync"

	"github.com/hyperledger/fabric-protos-go/common"
)

type BlockIterator struct {
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	NextStub        func() (*common.Block, common.Status)
	nextMutex       sync.RWMutex
	nextArgsForCall []struct {
	}
	nextReturns struct {
		result1 *common.Block
		result2 common.Status
	}
	nextReturnsOnCall map[int]struct {
		result1 *common.Block
		result2 common.Status
	}
	WaitForNextBlockStub        func() (bool, bool)
	waitForNextBlockMutex       sync.RWMutex
	waitForNextBlockArgsForCall []struct {
	}
	waitForNextBlockReturns struct {
		result1 bool
		result2 bool
	}
	waitForNextBlockReturnsOnCall map[int]struct {
		result1 bool
		result2 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *BlockIterator) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		fake.CloseStub()
	}
}

func (fake *BlockIterator) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *BlockIterator) CloseCalls(stub func()) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *BlockIterator) Next() (*common.Block, common.Status) {
	fake.nextMutex.Lock()
	ret, specificReturn := fake.nextReturnsOnCall[len(fake.nextArgsForCall)]
	fake.nextArgsForCall = append(fake.nextArgsForCall, struct {
	}{})
	fake.recordInvocation("Next", []interface{}{})
	fake.nextMutex.Unlock()
	if fake.NextStub != nil {
		return fake.NextStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.nextReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *BlockIterator) NextCallCount() int {
	fake.nextMutex.RLock()
	defer fake.nextMutex.RUnlock()
	return len(fake.nextArgsForCall)
}

func (fake *BlockIterator) NextCalls(stub func() (*common.Block, common.Status)) {
	fake.nextMutex.Lock()
	defer fake.nextMutex.Unlock()
	fake.NextStub = stub
}

func (fake *BlockIterator) NextReturns(result1 *common.Block, result2 common.Status) {
	fake.nextMutex.Lock()
	defer fake.nextMutex.Unlock()
	fake.NextStub = nil
	fake.nextReturns = struct {
		result1 *common.Block
		result2 common.Status
	}{result1, result2}
}

func (fake *BlockIterator) NextReturnsOnCall(i int, result1 *common.Block, result2 common.Status) {
	fake.nextMutex.Lock()
	defer fake.nextMutex.Unlock()
	fake.NextStub = nil
	if fake.nextReturnsOnCall == nil {
		fake.nextReturnsOnCall = make(map[int]struct {
			result1 *common.Block
			result2 common.Status
		})
	}
	fake.nextReturnsOnCall[i] = struct {
		result1 *common.Block
		result2 common.Status
	}{result1, result2}
}

func (fake *BlockIterator) WaitForNextBlock() (bool, bool) {
	fake.waitForNextBlockMutex.Lock()
	ret, specificReturn := fake.waitForNextBlockReturnsOnCall[len(fake.waitForNextBlockArgsForCall)]
	fake.waitForNextBlockArgsForCall = append(fake.waitForNextBlockArgsForCall, struct {
	}{})
	fake.recordInvocation("WaitForNextBlock", []interface{}{})
	fake.waitForNextBlockMutex.Unlock()
	if fake.WaitForNextBlockStub != nil {
		return fake.WaitForNextBlockStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.waitForNextBlockReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *BlockIterator) WaitForNextBlockCallCount() int {
	fake.waitForNextBlockMutex.RLock()
	defer fake.waitForNextBlockMutex.RUnlock()
	return len(fake.waitForNextBlockArgsForCall)
}

func (fake *BlockIterator) WaitForNextBlockCalls(stub func() (bool, bool)) {
	fake.waitForNextBlockMutex.Lock()
	defer fake.waitForNextBlockMutex.Unlock()
	fake.WaitForNextBlockStub = stub
}

func (fake *BlockIterator) WaitForNextBlockReturns(result1 bool, result2 bool) {
	fake.waitForNextBlockMutex.Lock()
	defer fake.waitForNextBlockMutex.Unlock()
	fake.WaitForNextBlockStub = nil
	fake.waitForNextBlockReturns = struct {
		result1 bool
		result2 bool
	}{result1, result2}
}

func (fake *BlockIterator) WaitForNextBlockReturnsOnCall(i int, result1 bool, result2 bool) {
	fake.waitForNextBlockMutex.Lock()
	defer fake.waitForNextBlockMutex.Unlock()
	fake.WaitForNextBlockStub = nil
	if fake.waitForNextBlockReturnsOnCall == nil {
		fake.waitForNextBlockReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 bool
		})
	}
	fake.waitForNextBlockReturnsOnCall[i] = struct {
		result1 bool
		result2 bool
	}{result1, result2}
}

func (fake *BlockIterator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.nextMutex.RLock()
	defer fake.nextMutex.RUnlock()
	fake.waitForNextBlockMutex.RLock()
	defer fake.waitForNextBlockMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *BlockIterator) recordInvocation(key string, args []interface{}) {
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
