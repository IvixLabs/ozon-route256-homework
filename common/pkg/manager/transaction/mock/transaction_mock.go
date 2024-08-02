// Code generated by http://github.com/gojuno/minimock (v3.3.13). DO NOT EDIT.

package mock

//go:generate minimock -i route256/common/pkg/manager/transaction.Transaction -o transaction_mock.go -n TransactionMock -p mock

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// TransactionMock implements transaction.Transaction
type TransactionMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcCommit          func(ctx context.Context) (err error)
	inspectFuncCommit   func(ctx context.Context)
	afterCommitCounter  uint64
	beforeCommitCounter uint64
	CommitMock          mTransactionMockCommit

	funcRollback          func(ctx context.Context) (err error)
	inspectFuncRollback   func(ctx context.Context)
	afterRollbackCounter  uint64
	beforeRollbackCounter uint64
	RollbackMock          mTransactionMockRollback
}

// NewTransactionMock returns a mock for transaction.Transaction
func NewTransactionMock(t minimock.Tester) *TransactionMock {
	m := &TransactionMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CommitMock = mTransactionMockCommit{mock: m}
	m.CommitMock.callArgs = []*TransactionMockCommitParams{}

	m.RollbackMock = mTransactionMockRollback{mock: m}
	m.RollbackMock.callArgs = []*TransactionMockRollbackParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mTransactionMockCommit struct {
	optional           bool
	mock               *TransactionMock
	defaultExpectation *TransactionMockCommitExpectation
	expectations       []*TransactionMockCommitExpectation

	callArgs []*TransactionMockCommitParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// TransactionMockCommitExpectation specifies expectation struct of the Transaction.Commit
type TransactionMockCommitExpectation struct {
	mock      *TransactionMock
	params    *TransactionMockCommitParams
	paramPtrs *TransactionMockCommitParamPtrs
	results   *TransactionMockCommitResults
	Counter   uint64
}

// TransactionMockCommitParams contains parameters of the Transaction.Commit
type TransactionMockCommitParams struct {
	ctx context.Context
}

// TransactionMockCommitParamPtrs contains pointers to parameters of the Transaction.Commit
type TransactionMockCommitParamPtrs struct {
	ctx *context.Context
}

// TransactionMockCommitResults contains results of the Transaction.Commit
type TransactionMockCommitResults struct {
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmCommit *mTransactionMockCommit) Optional() *mTransactionMockCommit {
	mmCommit.optional = true
	return mmCommit
}

// Expect sets up expected params for Transaction.Commit
func (mmCommit *mTransactionMockCommit) Expect(ctx context.Context) *mTransactionMockCommit {
	if mmCommit.mock.funcCommit != nil {
		mmCommit.mock.t.Fatalf("TransactionMock.Commit mock is already set by Set")
	}

	if mmCommit.defaultExpectation == nil {
		mmCommit.defaultExpectation = &TransactionMockCommitExpectation{}
	}

	if mmCommit.defaultExpectation.paramPtrs != nil {
		mmCommit.mock.t.Fatalf("TransactionMock.Commit mock is already set by ExpectParams functions")
	}

	mmCommit.defaultExpectation.params = &TransactionMockCommitParams{ctx}
	for _, e := range mmCommit.expectations {
		if minimock.Equal(e.params, mmCommit.defaultExpectation.params) {
			mmCommit.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmCommit.defaultExpectation.params)
		}
	}

	return mmCommit
}

// ExpectCtxParam1 sets up expected param ctx for Transaction.Commit
func (mmCommit *mTransactionMockCommit) ExpectCtxParam1(ctx context.Context) *mTransactionMockCommit {
	if mmCommit.mock.funcCommit != nil {
		mmCommit.mock.t.Fatalf("TransactionMock.Commit mock is already set by Set")
	}

	if mmCommit.defaultExpectation == nil {
		mmCommit.defaultExpectation = &TransactionMockCommitExpectation{}
	}

	if mmCommit.defaultExpectation.params != nil {
		mmCommit.mock.t.Fatalf("TransactionMock.Commit mock is already set by Expect")
	}

	if mmCommit.defaultExpectation.paramPtrs == nil {
		mmCommit.defaultExpectation.paramPtrs = &TransactionMockCommitParamPtrs{}
	}
	mmCommit.defaultExpectation.paramPtrs.ctx = &ctx

	return mmCommit
}

// Inspect accepts an inspector function that has same arguments as the Transaction.Commit
func (mmCommit *mTransactionMockCommit) Inspect(f func(ctx context.Context)) *mTransactionMockCommit {
	if mmCommit.mock.inspectFuncCommit != nil {
		mmCommit.mock.t.Fatalf("Inspect function is already set for TransactionMock.Commit")
	}

	mmCommit.mock.inspectFuncCommit = f

	return mmCommit
}

// Return sets up results that will be returned by Transaction.Commit
func (mmCommit *mTransactionMockCommit) Return(err error) *TransactionMock {
	if mmCommit.mock.funcCommit != nil {
		mmCommit.mock.t.Fatalf("TransactionMock.Commit mock is already set by Set")
	}

	if mmCommit.defaultExpectation == nil {
		mmCommit.defaultExpectation = &TransactionMockCommitExpectation{mock: mmCommit.mock}
	}
	mmCommit.defaultExpectation.results = &TransactionMockCommitResults{err}
	return mmCommit.mock
}

// Set uses given function f to mock the Transaction.Commit method
func (mmCommit *mTransactionMockCommit) Set(f func(ctx context.Context) (err error)) *TransactionMock {
	if mmCommit.defaultExpectation != nil {
		mmCommit.mock.t.Fatalf("Default expectation is already set for the Transaction.Commit method")
	}

	if len(mmCommit.expectations) > 0 {
		mmCommit.mock.t.Fatalf("Some expectations are already set for the Transaction.Commit method")
	}

	mmCommit.mock.funcCommit = f
	return mmCommit.mock
}

// When sets expectation for the Transaction.Commit which will trigger the result defined by the following
// Then helper
func (mmCommit *mTransactionMockCommit) When(ctx context.Context) *TransactionMockCommitExpectation {
	if mmCommit.mock.funcCommit != nil {
		mmCommit.mock.t.Fatalf("TransactionMock.Commit mock is already set by Set")
	}

	expectation := &TransactionMockCommitExpectation{
		mock:   mmCommit.mock,
		params: &TransactionMockCommitParams{ctx},
	}
	mmCommit.expectations = append(mmCommit.expectations, expectation)
	return expectation
}

// Then sets up Transaction.Commit return parameters for the expectation previously defined by the When method
func (e *TransactionMockCommitExpectation) Then(err error) *TransactionMock {
	e.results = &TransactionMockCommitResults{err}
	return e.mock
}

// Times sets number of times Transaction.Commit should be invoked
func (mmCommit *mTransactionMockCommit) Times(n uint64) *mTransactionMockCommit {
	if n == 0 {
		mmCommit.mock.t.Fatalf("Times of TransactionMock.Commit mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmCommit.expectedInvocations, n)
	return mmCommit
}

func (mmCommit *mTransactionMockCommit) invocationsDone() bool {
	if len(mmCommit.expectations) == 0 && mmCommit.defaultExpectation == nil && mmCommit.mock.funcCommit == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmCommit.mock.afterCommitCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmCommit.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// Commit implements transaction.Transaction
func (mmCommit *TransactionMock) Commit(ctx context.Context) (err error) {
	mm_atomic.AddUint64(&mmCommit.beforeCommitCounter, 1)
	defer mm_atomic.AddUint64(&mmCommit.afterCommitCounter, 1)

	if mmCommit.inspectFuncCommit != nil {
		mmCommit.inspectFuncCommit(ctx)
	}

	mm_params := TransactionMockCommitParams{ctx}

	// Record call args
	mmCommit.CommitMock.mutex.Lock()
	mmCommit.CommitMock.callArgs = append(mmCommit.CommitMock.callArgs, &mm_params)
	mmCommit.CommitMock.mutex.Unlock()

	for _, e := range mmCommit.CommitMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmCommit.CommitMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmCommit.CommitMock.defaultExpectation.Counter, 1)
		mm_want := mmCommit.CommitMock.defaultExpectation.params
		mm_want_ptrs := mmCommit.CommitMock.defaultExpectation.paramPtrs

		mm_got := TransactionMockCommitParams{ctx}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmCommit.t.Errorf("TransactionMock.Commit got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmCommit.t.Errorf("TransactionMock.Commit got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmCommit.CommitMock.defaultExpectation.results
		if mm_results == nil {
			mmCommit.t.Fatal("No results are set for the TransactionMock.Commit")
		}
		return (*mm_results).err
	}
	if mmCommit.funcCommit != nil {
		return mmCommit.funcCommit(ctx)
	}
	mmCommit.t.Fatalf("Unexpected call to TransactionMock.Commit. %v", ctx)
	return
}

// CommitAfterCounter returns a count of finished TransactionMock.Commit invocations
func (mmCommit *TransactionMock) CommitAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCommit.afterCommitCounter)
}

// CommitBeforeCounter returns a count of TransactionMock.Commit invocations
func (mmCommit *TransactionMock) CommitBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCommit.beforeCommitCounter)
}

// Calls returns a list of arguments used in each call to TransactionMock.Commit.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmCommit *mTransactionMockCommit) Calls() []*TransactionMockCommitParams {
	mmCommit.mutex.RLock()

	argCopy := make([]*TransactionMockCommitParams, len(mmCommit.callArgs))
	copy(argCopy, mmCommit.callArgs)

	mmCommit.mutex.RUnlock()

	return argCopy
}

// MinimockCommitDone returns true if the count of the Commit invocations corresponds
// the number of defined expectations
func (m *TransactionMock) MinimockCommitDone() bool {
	if m.CommitMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.CommitMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.CommitMock.invocationsDone()
}

// MinimockCommitInspect logs each unmet expectation
func (m *TransactionMock) MinimockCommitInspect() {
	for _, e := range m.CommitMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TransactionMock.Commit with params: %#v", *e.params)
		}
	}

	afterCommitCounter := mm_atomic.LoadUint64(&m.afterCommitCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.CommitMock.defaultExpectation != nil && afterCommitCounter < 1 {
		if m.CommitMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to TransactionMock.Commit")
		} else {
			m.t.Errorf("Expected call to TransactionMock.Commit with params: %#v", *m.CommitMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCommit != nil && afterCommitCounter < 1 {
		m.t.Error("Expected call to TransactionMock.Commit")
	}

	if !m.CommitMock.invocationsDone() && afterCommitCounter > 0 {
		m.t.Errorf("Expected %d calls to TransactionMock.Commit but found %d calls",
			mm_atomic.LoadUint64(&m.CommitMock.expectedInvocations), afterCommitCounter)
	}
}

type mTransactionMockRollback struct {
	optional           bool
	mock               *TransactionMock
	defaultExpectation *TransactionMockRollbackExpectation
	expectations       []*TransactionMockRollbackExpectation

	callArgs []*TransactionMockRollbackParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// TransactionMockRollbackExpectation specifies expectation struct of the Transaction.Rollback
type TransactionMockRollbackExpectation struct {
	mock      *TransactionMock
	params    *TransactionMockRollbackParams
	paramPtrs *TransactionMockRollbackParamPtrs
	results   *TransactionMockRollbackResults
	Counter   uint64
}

// TransactionMockRollbackParams contains parameters of the Transaction.Rollback
type TransactionMockRollbackParams struct {
	ctx context.Context
}

// TransactionMockRollbackParamPtrs contains pointers to parameters of the Transaction.Rollback
type TransactionMockRollbackParamPtrs struct {
	ctx *context.Context
}

// TransactionMockRollbackResults contains results of the Transaction.Rollback
type TransactionMockRollbackResults struct {
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmRollback *mTransactionMockRollback) Optional() *mTransactionMockRollback {
	mmRollback.optional = true
	return mmRollback
}

// Expect sets up expected params for Transaction.Rollback
func (mmRollback *mTransactionMockRollback) Expect(ctx context.Context) *mTransactionMockRollback {
	if mmRollback.mock.funcRollback != nil {
		mmRollback.mock.t.Fatalf("TransactionMock.Rollback mock is already set by Set")
	}

	if mmRollback.defaultExpectation == nil {
		mmRollback.defaultExpectation = &TransactionMockRollbackExpectation{}
	}

	if mmRollback.defaultExpectation.paramPtrs != nil {
		mmRollback.mock.t.Fatalf("TransactionMock.Rollback mock is already set by ExpectParams functions")
	}

	mmRollback.defaultExpectation.params = &TransactionMockRollbackParams{ctx}
	for _, e := range mmRollback.expectations {
		if minimock.Equal(e.params, mmRollback.defaultExpectation.params) {
			mmRollback.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRollback.defaultExpectation.params)
		}
	}

	return mmRollback
}

// ExpectCtxParam1 sets up expected param ctx for Transaction.Rollback
func (mmRollback *mTransactionMockRollback) ExpectCtxParam1(ctx context.Context) *mTransactionMockRollback {
	if mmRollback.mock.funcRollback != nil {
		mmRollback.mock.t.Fatalf("TransactionMock.Rollback mock is already set by Set")
	}

	if mmRollback.defaultExpectation == nil {
		mmRollback.defaultExpectation = &TransactionMockRollbackExpectation{}
	}

	if mmRollback.defaultExpectation.params != nil {
		mmRollback.mock.t.Fatalf("TransactionMock.Rollback mock is already set by Expect")
	}

	if mmRollback.defaultExpectation.paramPtrs == nil {
		mmRollback.defaultExpectation.paramPtrs = &TransactionMockRollbackParamPtrs{}
	}
	mmRollback.defaultExpectation.paramPtrs.ctx = &ctx

	return mmRollback
}

// Inspect accepts an inspector function that has same arguments as the Transaction.Rollback
func (mmRollback *mTransactionMockRollback) Inspect(f func(ctx context.Context)) *mTransactionMockRollback {
	if mmRollback.mock.inspectFuncRollback != nil {
		mmRollback.mock.t.Fatalf("Inspect function is already set for TransactionMock.Rollback")
	}

	mmRollback.mock.inspectFuncRollback = f

	return mmRollback
}

// Return sets up results that will be returned by Transaction.Rollback
func (mmRollback *mTransactionMockRollback) Return(err error) *TransactionMock {
	if mmRollback.mock.funcRollback != nil {
		mmRollback.mock.t.Fatalf("TransactionMock.Rollback mock is already set by Set")
	}

	if mmRollback.defaultExpectation == nil {
		mmRollback.defaultExpectation = &TransactionMockRollbackExpectation{mock: mmRollback.mock}
	}
	mmRollback.defaultExpectation.results = &TransactionMockRollbackResults{err}
	return mmRollback.mock
}

// Set uses given function f to mock the Transaction.Rollback method
func (mmRollback *mTransactionMockRollback) Set(f func(ctx context.Context) (err error)) *TransactionMock {
	if mmRollback.defaultExpectation != nil {
		mmRollback.mock.t.Fatalf("Default expectation is already set for the Transaction.Rollback method")
	}

	if len(mmRollback.expectations) > 0 {
		mmRollback.mock.t.Fatalf("Some expectations are already set for the Transaction.Rollback method")
	}

	mmRollback.mock.funcRollback = f
	return mmRollback.mock
}

// When sets expectation for the Transaction.Rollback which will trigger the result defined by the following
// Then helper
func (mmRollback *mTransactionMockRollback) When(ctx context.Context) *TransactionMockRollbackExpectation {
	if mmRollback.mock.funcRollback != nil {
		mmRollback.mock.t.Fatalf("TransactionMock.Rollback mock is already set by Set")
	}

	expectation := &TransactionMockRollbackExpectation{
		mock:   mmRollback.mock,
		params: &TransactionMockRollbackParams{ctx},
	}
	mmRollback.expectations = append(mmRollback.expectations, expectation)
	return expectation
}

// Then sets up Transaction.Rollback return parameters for the expectation previously defined by the When method
func (e *TransactionMockRollbackExpectation) Then(err error) *TransactionMock {
	e.results = &TransactionMockRollbackResults{err}
	return e.mock
}

// Times sets number of times Transaction.Rollback should be invoked
func (mmRollback *mTransactionMockRollback) Times(n uint64) *mTransactionMockRollback {
	if n == 0 {
		mmRollback.mock.t.Fatalf("Times of TransactionMock.Rollback mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmRollback.expectedInvocations, n)
	return mmRollback
}

func (mmRollback *mTransactionMockRollback) invocationsDone() bool {
	if len(mmRollback.expectations) == 0 && mmRollback.defaultExpectation == nil && mmRollback.mock.funcRollback == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmRollback.mock.afterRollbackCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmRollback.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// Rollback implements transaction.Transaction
func (mmRollback *TransactionMock) Rollback(ctx context.Context) (err error) {
	mm_atomic.AddUint64(&mmRollback.beforeRollbackCounter, 1)
	defer mm_atomic.AddUint64(&mmRollback.afterRollbackCounter, 1)

	if mmRollback.inspectFuncRollback != nil {
		mmRollback.inspectFuncRollback(ctx)
	}

	mm_params := TransactionMockRollbackParams{ctx}

	// Record call args
	mmRollback.RollbackMock.mutex.Lock()
	mmRollback.RollbackMock.callArgs = append(mmRollback.RollbackMock.callArgs, &mm_params)
	mmRollback.RollbackMock.mutex.Unlock()

	for _, e := range mmRollback.RollbackMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmRollback.RollbackMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmRollback.RollbackMock.defaultExpectation.Counter, 1)
		mm_want := mmRollback.RollbackMock.defaultExpectation.params
		mm_want_ptrs := mmRollback.RollbackMock.defaultExpectation.paramPtrs

		mm_got := TransactionMockRollbackParams{ctx}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmRollback.t.Errorf("TransactionMock.Rollback got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRollback.t.Errorf("TransactionMock.Rollback got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmRollback.RollbackMock.defaultExpectation.results
		if mm_results == nil {
			mmRollback.t.Fatal("No results are set for the TransactionMock.Rollback")
		}
		return (*mm_results).err
	}
	if mmRollback.funcRollback != nil {
		return mmRollback.funcRollback(ctx)
	}
	mmRollback.t.Fatalf("Unexpected call to TransactionMock.Rollback. %v", ctx)
	return
}

// RollbackAfterCounter returns a count of finished TransactionMock.Rollback invocations
func (mmRollback *TransactionMock) RollbackAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRollback.afterRollbackCounter)
}

// RollbackBeforeCounter returns a count of TransactionMock.Rollback invocations
func (mmRollback *TransactionMock) RollbackBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRollback.beforeRollbackCounter)
}

// Calls returns a list of arguments used in each call to TransactionMock.Rollback.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRollback *mTransactionMockRollback) Calls() []*TransactionMockRollbackParams {
	mmRollback.mutex.RLock()

	argCopy := make([]*TransactionMockRollbackParams, len(mmRollback.callArgs))
	copy(argCopy, mmRollback.callArgs)

	mmRollback.mutex.RUnlock()

	return argCopy
}

// MinimockRollbackDone returns true if the count of the Rollback invocations corresponds
// the number of defined expectations
func (m *TransactionMock) MinimockRollbackDone() bool {
	if m.RollbackMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.RollbackMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.RollbackMock.invocationsDone()
}

// MinimockRollbackInspect logs each unmet expectation
func (m *TransactionMock) MinimockRollbackInspect() {
	for _, e := range m.RollbackMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TransactionMock.Rollback with params: %#v", *e.params)
		}
	}

	afterRollbackCounter := mm_atomic.LoadUint64(&m.afterRollbackCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.RollbackMock.defaultExpectation != nil && afterRollbackCounter < 1 {
		if m.RollbackMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to TransactionMock.Rollback")
		} else {
			m.t.Errorf("Expected call to TransactionMock.Rollback with params: %#v", *m.RollbackMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRollback != nil && afterRollbackCounter < 1 {
		m.t.Error("Expected call to TransactionMock.Rollback")
	}

	if !m.RollbackMock.invocationsDone() && afterRollbackCounter > 0 {
		m.t.Errorf("Expected %d calls to TransactionMock.Rollback but found %d calls",
			mm_atomic.LoadUint64(&m.RollbackMock.expectedInvocations), afterRollbackCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *TransactionMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockCommitInspect()

			m.MinimockRollbackInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *TransactionMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *TransactionMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockCommitDone() &&
		m.MinimockRollbackDone()
}
