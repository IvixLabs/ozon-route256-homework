// Code generated by http://github.com/gojuno/minimock (v3.3.13). DO NOT EDIT.

package mock

//go:generate minimock -i route256/loms/internal/usecase/stock.ReservedStockRepository -o reserved_stock_repository_mock.go -n ReservedStockRepositoryMock -p mock

import (
	"context"
	"route256/loms/internal/model"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// ReservedStockRepositoryMock implements stock.ReservedStockRepository
type ReservedStockRepositoryMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetLocked          func(ctx context.Context, orderID model.OrderID, sku model.Sku) (rp1 *model.ReservedStock, err error)
	inspectFuncGetLocked   func(ctx context.Context, orderID model.OrderID, sku model.Sku)
	afterGetLockedCounter  uint64
	beforeGetLockedCounter uint64
	GetLockedMock          mReservedStockRepositoryMockGetLocked

	funcSave          func(ctx context.Context, rStock *model.ReservedStock) (err error)
	inspectFuncSave   func(ctx context.Context, rStock *model.ReservedStock)
	afterSaveCounter  uint64
	beforeSaveCounter uint64
	SaveMock          mReservedStockRepositoryMockSave
}

// NewReservedStockRepositoryMock returns a mock for stock.ReservedStockRepository
func NewReservedStockRepositoryMock(t minimock.Tester) *ReservedStockRepositoryMock {
	m := &ReservedStockRepositoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetLockedMock = mReservedStockRepositoryMockGetLocked{mock: m}
	m.GetLockedMock.callArgs = []*ReservedStockRepositoryMockGetLockedParams{}

	m.SaveMock = mReservedStockRepositoryMockSave{mock: m}
	m.SaveMock.callArgs = []*ReservedStockRepositoryMockSaveParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mReservedStockRepositoryMockGetLocked struct {
	optional           bool
	mock               *ReservedStockRepositoryMock
	defaultExpectation *ReservedStockRepositoryMockGetLockedExpectation
	expectations       []*ReservedStockRepositoryMockGetLockedExpectation

	callArgs []*ReservedStockRepositoryMockGetLockedParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// ReservedStockRepositoryMockGetLockedExpectation specifies expectation struct of the ReservedStockRepository.GetLocked
type ReservedStockRepositoryMockGetLockedExpectation struct {
	mock      *ReservedStockRepositoryMock
	params    *ReservedStockRepositoryMockGetLockedParams
	paramPtrs *ReservedStockRepositoryMockGetLockedParamPtrs
	results   *ReservedStockRepositoryMockGetLockedResults
	Counter   uint64
}

// ReservedStockRepositoryMockGetLockedParams contains parameters of the ReservedStockRepository.GetLocked
type ReservedStockRepositoryMockGetLockedParams struct {
	ctx     context.Context
	orderID model.OrderID
	sku     model.Sku
}

// ReservedStockRepositoryMockGetLockedParamPtrs contains pointers to parameters of the ReservedStockRepository.GetLocked
type ReservedStockRepositoryMockGetLockedParamPtrs struct {
	ctx     *context.Context
	orderID *model.OrderID
	sku     *model.Sku
}

// ReservedStockRepositoryMockGetLockedResults contains results of the ReservedStockRepository.GetLocked
type ReservedStockRepositoryMockGetLockedResults struct {
	rp1 *model.ReservedStock
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmGetLocked *mReservedStockRepositoryMockGetLocked) Optional() *mReservedStockRepositoryMockGetLocked {
	mmGetLocked.optional = true
	return mmGetLocked
}

// Expect sets up expected params for ReservedStockRepository.GetLocked
func (mmGetLocked *mReservedStockRepositoryMockGetLocked) Expect(ctx context.Context, orderID model.OrderID, sku model.Sku) *mReservedStockRepositoryMockGetLocked {
	if mmGetLocked.mock.funcGetLocked != nil {
		mmGetLocked.mock.t.Fatalf("ReservedStockRepositoryMock.GetLocked mock is already set by Set")
	}

	if mmGetLocked.defaultExpectation == nil {
		mmGetLocked.defaultExpectation = &ReservedStockRepositoryMockGetLockedExpectation{}
	}

	if mmGetLocked.defaultExpectation.paramPtrs != nil {
		mmGetLocked.mock.t.Fatalf("ReservedStockRepositoryMock.GetLocked mock is already set by ExpectParams functions")
	}

	mmGetLocked.defaultExpectation.params = &ReservedStockRepositoryMockGetLockedParams{ctx, orderID, sku}
	for _, e := range mmGetLocked.expectations {
		if minimock.Equal(e.params, mmGetLocked.defaultExpectation.params) {
			mmGetLocked.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetLocked.defaultExpectation.params)
		}
	}

	return mmGetLocked
}

// ExpectCtxParam1 sets up expected param ctx for ReservedStockRepository.GetLocked
func (mmGetLocked *mReservedStockRepositoryMockGetLocked) ExpectCtxParam1(ctx context.Context) *mReservedStockRepositoryMockGetLocked {
	if mmGetLocked.mock.funcGetLocked != nil {
		mmGetLocked.mock.t.Fatalf("ReservedStockRepositoryMock.GetLocked mock is already set by Set")
	}

	if mmGetLocked.defaultExpectation == nil {
		mmGetLocked.defaultExpectation = &ReservedStockRepositoryMockGetLockedExpectation{}
	}

	if mmGetLocked.defaultExpectation.params != nil {
		mmGetLocked.mock.t.Fatalf("ReservedStockRepositoryMock.GetLocked mock is already set by Expect")
	}

	if mmGetLocked.defaultExpectation.paramPtrs == nil {
		mmGetLocked.defaultExpectation.paramPtrs = &ReservedStockRepositoryMockGetLockedParamPtrs{}
	}
	mmGetLocked.defaultExpectation.paramPtrs.ctx = &ctx

	return mmGetLocked
}

// ExpectOrderIDParam2 sets up expected param orderID for ReservedStockRepository.GetLocked
func (mmGetLocked *mReservedStockRepositoryMockGetLocked) ExpectOrderIDParam2(orderID model.OrderID) *mReservedStockRepositoryMockGetLocked {
	if mmGetLocked.mock.funcGetLocked != nil {
		mmGetLocked.mock.t.Fatalf("ReservedStockRepositoryMock.GetLocked mock is already set by Set")
	}

	if mmGetLocked.defaultExpectation == nil {
		mmGetLocked.defaultExpectation = &ReservedStockRepositoryMockGetLockedExpectation{}
	}

	if mmGetLocked.defaultExpectation.params != nil {
		mmGetLocked.mock.t.Fatalf("ReservedStockRepositoryMock.GetLocked mock is already set by Expect")
	}

	if mmGetLocked.defaultExpectation.paramPtrs == nil {
		mmGetLocked.defaultExpectation.paramPtrs = &ReservedStockRepositoryMockGetLockedParamPtrs{}
	}
	mmGetLocked.defaultExpectation.paramPtrs.orderID = &orderID

	return mmGetLocked
}

// ExpectSkuParam3 sets up expected param sku for ReservedStockRepository.GetLocked
func (mmGetLocked *mReservedStockRepositoryMockGetLocked) ExpectSkuParam3(sku model.Sku) *mReservedStockRepositoryMockGetLocked {
	if mmGetLocked.mock.funcGetLocked != nil {
		mmGetLocked.mock.t.Fatalf("ReservedStockRepositoryMock.GetLocked mock is already set by Set")
	}

	if mmGetLocked.defaultExpectation == nil {
		mmGetLocked.defaultExpectation = &ReservedStockRepositoryMockGetLockedExpectation{}
	}

	if mmGetLocked.defaultExpectation.params != nil {
		mmGetLocked.mock.t.Fatalf("ReservedStockRepositoryMock.GetLocked mock is already set by Expect")
	}

	if mmGetLocked.defaultExpectation.paramPtrs == nil {
		mmGetLocked.defaultExpectation.paramPtrs = &ReservedStockRepositoryMockGetLockedParamPtrs{}
	}
	mmGetLocked.defaultExpectation.paramPtrs.sku = &sku

	return mmGetLocked
}

// Inspect accepts an inspector function that has same arguments as the ReservedStockRepository.GetLocked
func (mmGetLocked *mReservedStockRepositoryMockGetLocked) Inspect(f func(ctx context.Context, orderID model.OrderID, sku model.Sku)) *mReservedStockRepositoryMockGetLocked {
	if mmGetLocked.mock.inspectFuncGetLocked != nil {
		mmGetLocked.mock.t.Fatalf("Inspect function is already set for ReservedStockRepositoryMock.GetLocked")
	}

	mmGetLocked.mock.inspectFuncGetLocked = f

	return mmGetLocked
}

// Return sets up results that will be returned by ReservedStockRepository.GetLocked
func (mmGetLocked *mReservedStockRepositoryMockGetLocked) Return(rp1 *model.ReservedStock, err error) *ReservedStockRepositoryMock {
	if mmGetLocked.mock.funcGetLocked != nil {
		mmGetLocked.mock.t.Fatalf("ReservedStockRepositoryMock.GetLocked mock is already set by Set")
	}

	if mmGetLocked.defaultExpectation == nil {
		mmGetLocked.defaultExpectation = &ReservedStockRepositoryMockGetLockedExpectation{mock: mmGetLocked.mock}
	}
	mmGetLocked.defaultExpectation.results = &ReservedStockRepositoryMockGetLockedResults{rp1, err}
	return mmGetLocked.mock
}

// Set uses given function f to mock the ReservedStockRepository.GetLocked method
func (mmGetLocked *mReservedStockRepositoryMockGetLocked) Set(f func(ctx context.Context, orderID model.OrderID, sku model.Sku) (rp1 *model.ReservedStock, err error)) *ReservedStockRepositoryMock {
	if mmGetLocked.defaultExpectation != nil {
		mmGetLocked.mock.t.Fatalf("Default expectation is already set for the ReservedStockRepository.GetLocked method")
	}

	if len(mmGetLocked.expectations) > 0 {
		mmGetLocked.mock.t.Fatalf("Some expectations are already set for the ReservedStockRepository.GetLocked method")
	}

	mmGetLocked.mock.funcGetLocked = f
	return mmGetLocked.mock
}

// When sets expectation for the ReservedStockRepository.GetLocked which will trigger the result defined by the following
// Then helper
func (mmGetLocked *mReservedStockRepositoryMockGetLocked) When(ctx context.Context, orderID model.OrderID, sku model.Sku) *ReservedStockRepositoryMockGetLockedExpectation {
	if mmGetLocked.mock.funcGetLocked != nil {
		mmGetLocked.mock.t.Fatalf("ReservedStockRepositoryMock.GetLocked mock is already set by Set")
	}

	expectation := &ReservedStockRepositoryMockGetLockedExpectation{
		mock:   mmGetLocked.mock,
		params: &ReservedStockRepositoryMockGetLockedParams{ctx, orderID, sku},
	}
	mmGetLocked.expectations = append(mmGetLocked.expectations, expectation)
	return expectation
}

// Then sets up ReservedStockRepository.GetLocked return parameters for the expectation previously defined by the When method
func (e *ReservedStockRepositoryMockGetLockedExpectation) Then(rp1 *model.ReservedStock, err error) *ReservedStockRepositoryMock {
	e.results = &ReservedStockRepositoryMockGetLockedResults{rp1, err}
	return e.mock
}

// Times sets number of times ReservedStockRepository.GetLocked should be invoked
func (mmGetLocked *mReservedStockRepositoryMockGetLocked) Times(n uint64) *mReservedStockRepositoryMockGetLocked {
	if n == 0 {
		mmGetLocked.mock.t.Fatalf("Times of ReservedStockRepositoryMock.GetLocked mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmGetLocked.expectedInvocations, n)
	return mmGetLocked
}

func (mmGetLocked *mReservedStockRepositoryMockGetLocked) invocationsDone() bool {
	if len(mmGetLocked.expectations) == 0 && mmGetLocked.defaultExpectation == nil && mmGetLocked.mock.funcGetLocked == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmGetLocked.mock.afterGetLockedCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmGetLocked.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// GetLocked implements stock.ReservedStockRepository
func (mmGetLocked *ReservedStockRepositoryMock) GetLocked(ctx context.Context, orderID model.OrderID, sku model.Sku) (rp1 *model.ReservedStock, err error) {
	mm_atomic.AddUint64(&mmGetLocked.beforeGetLockedCounter, 1)
	defer mm_atomic.AddUint64(&mmGetLocked.afterGetLockedCounter, 1)

	if mmGetLocked.inspectFuncGetLocked != nil {
		mmGetLocked.inspectFuncGetLocked(ctx, orderID, sku)
	}

	mm_params := ReservedStockRepositoryMockGetLockedParams{ctx, orderID, sku}

	// Record call args
	mmGetLocked.GetLockedMock.mutex.Lock()
	mmGetLocked.GetLockedMock.callArgs = append(mmGetLocked.GetLockedMock.callArgs, &mm_params)
	mmGetLocked.GetLockedMock.mutex.Unlock()

	for _, e := range mmGetLocked.GetLockedMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.rp1, e.results.err
		}
	}

	if mmGetLocked.GetLockedMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetLocked.GetLockedMock.defaultExpectation.Counter, 1)
		mm_want := mmGetLocked.GetLockedMock.defaultExpectation.params
		mm_want_ptrs := mmGetLocked.GetLockedMock.defaultExpectation.paramPtrs

		mm_got := ReservedStockRepositoryMockGetLockedParams{ctx, orderID, sku}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmGetLocked.t.Errorf("ReservedStockRepositoryMock.GetLocked got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.orderID != nil && !minimock.Equal(*mm_want_ptrs.orderID, mm_got.orderID) {
				mmGetLocked.t.Errorf("ReservedStockRepositoryMock.GetLocked got unexpected parameter orderID, want: %#v, got: %#v%s\n", *mm_want_ptrs.orderID, mm_got.orderID, minimock.Diff(*mm_want_ptrs.orderID, mm_got.orderID))
			}

			if mm_want_ptrs.sku != nil && !minimock.Equal(*mm_want_ptrs.sku, mm_got.sku) {
				mmGetLocked.t.Errorf("ReservedStockRepositoryMock.GetLocked got unexpected parameter sku, want: %#v, got: %#v%s\n", *mm_want_ptrs.sku, mm_got.sku, minimock.Diff(*mm_want_ptrs.sku, mm_got.sku))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetLocked.t.Errorf("ReservedStockRepositoryMock.GetLocked got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetLocked.GetLockedMock.defaultExpectation.results
		if mm_results == nil {
			mmGetLocked.t.Fatal("No results are set for the ReservedStockRepositoryMock.GetLocked")
		}
		return (*mm_results).rp1, (*mm_results).err
	}
	if mmGetLocked.funcGetLocked != nil {
		return mmGetLocked.funcGetLocked(ctx, orderID, sku)
	}
	mmGetLocked.t.Fatalf("Unexpected call to ReservedStockRepositoryMock.GetLocked. %v %v %v", ctx, orderID, sku)
	return
}

// GetLockedAfterCounter returns a count of finished ReservedStockRepositoryMock.GetLocked invocations
func (mmGetLocked *ReservedStockRepositoryMock) GetLockedAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetLocked.afterGetLockedCounter)
}

// GetLockedBeforeCounter returns a count of ReservedStockRepositoryMock.GetLocked invocations
func (mmGetLocked *ReservedStockRepositoryMock) GetLockedBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetLocked.beforeGetLockedCounter)
}

// Calls returns a list of arguments used in each call to ReservedStockRepositoryMock.GetLocked.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetLocked *mReservedStockRepositoryMockGetLocked) Calls() []*ReservedStockRepositoryMockGetLockedParams {
	mmGetLocked.mutex.RLock()

	argCopy := make([]*ReservedStockRepositoryMockGetLockedParams, len(mmGetLocked.callArgs))
	copy(argCopy, mmGetLocked.callArgs)

	mmGetLocked.mutex.RUnlock()

	return argCopy
}

// MinimockGetLockedDone returns true if the count of the GetLocked invocations corresponds
// the number of defined expectations
func (m *ReservedStockRepositoryMock) MinimockGetLockedDone() bool {
	if m.GetLockedMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.GetLockedMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.GetLockedMock.invocationsDone()
}

// MinimockGetLockedInspect logs each unmet expectation
func (m *ReservedStockRepositoryMock) MinimockGetLockedInspect() {
	for _, e := range m.GetLockedMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ReservedStockRepositoryMock.GetLocked with params: %#v", *e.params)
		}
	}

	afterGetLockedCounter := mm_atomic.LoadUint64(&m.afterGetLockedCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.GetLockedMock.defaultExpectation != nil && afterGetLockedCounter < 1 {
		if m.GetLockedMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ReservedStockRepositoryMock.GetLocked")
		} else {
			m.t.Errorf("Expected call to ReservedStockRepositoryMock.GetLocked with params: %#v", *m.GetLockedMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetLocked != nil && afterGetLockedCounter < 1 {
		m.t.Error("Expected call to ReservedStockRepositoryMock.GetLocked")
	}

	if !m.GetLockedMock.invocationsDone() && afterGetLockedCounter > 0 {
		m.t.Errorf("Expected %d calls to ReservedStockRepositoryMock.GetLocked but found %d calls",
			mm_atomic.LoadUint64(&m.GetLockedMock.expectedInvocations), afterGetLockedCounter)
	}
}

type mReservedStockRepositoryMockSave struct {
	optional           bool
	mock               *ReservedStockRepositoryMock
	defaultExpectation *ReservedStockRepositoryMockSaveExpectation
	expectations       []*ReservedStockRepositoryMockSaveExpectation

	callArgs []*ReservedStockRepositoryMockSaveParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// ReservedStockRepositoryMockSaveExpectation specifies expectation struct of the ReservedStockRepository.Save
type ReservedStockRepositoryMockSaveExpectation struct {
	mock      *ReservedStockRepositoryMock
	params    *ReservedStockRepositoryMockSaveParams
	paramPtrs *ReservedStockRepositoryMockSaveParamPtrs
	results   *ReservedStockRepositoryMockSaveResults
	Counter   uint64
}

// ReservedStockRepositoryMockSaveParams contains parameters of the ReservedStockRepository.Save
type ReservedStockRepositoryMockSaveParams struct {
	ctx    context.Context
	rStock *model.ReservedStock
}

// ReservedStockRepositoryMockSaveParamPtrs contains pointers to parameters of the ReservedStockRepository.Save
type ReservedStockRepositoryMockSaveParamPtrs struct {
	ctx    *context.Context
	rStock **model.ReservedStock
}

// ReservedStockRepositoryMockSaveResults contains results of the ReservedStockRepository.Save
type ReservedStockRepositoryMockSaveResults struct {
	err error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmSave *mReservedStockRepositoryMockSave) Optional() *mReservedStockRepositoryMockSave {
	mmSave.optional = true
	return mmSave
}

// Expect sets up expected params for ReservedStockRepository.Save
func (mmSave *mReservedStockRepositoryMockSave) Expect(ctx context.Context, rStock *model.ReservedStock) *mReservedStockRepositoryMockSave {
	if mmSave.mock.funcSave != nil {
		mmSave.mock.t.Fatalf("ReservedStockRepositoryMock.Save mock is already set by Set")
	}

	if mmSave.defaultExpectation == nil {
		mmSave.defaultExpectation = &ReservedStockRepositoryMockSaveExpectation{}
	}

	if mmSave.defaultExpectation.paramPtrs != nil {
		mmSave.mock.t.Fatalf("ReservedStockRepositoryMock.Save mock is already set by ExpectParams functions")
	}

	mmSave.defaultExpectation.params = &ReservedStockRepositoryMockSaveParams{ctx, rStock}
	for _, e := range mmSave.expectations {
		if minimock.Equal(e.params, mmSave.defaultExpectation.params) {
			mmSave.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmSave.defaultExpectation.params)
		}
	}

	return mmSave
}

// ExpectCtxParam1 sets up expected param ctx for ReservedStockRepository.Save
func (mmSave *mReservedStockRepositoryMockSave) ExpectCtxParam1(ctx context.Context) *mReservedStockRepositoryMockSave {
	if mmSave.mock.funcSave != nil {
		mmSave.mock.t.Fatalf("ReservedStockRepositoryMock.Save mock is already set by Set")
	}

	if mmSave.defaultExpectation == nil {
		mmSave.defaultExpectation = &ReservedStockRepositoryMockSaveExpectation{}
	}

	if mmSave.defaultExpectation.params != nil {
		mmSave.mock.t.Fatalf("ReservedStockRepositoryMock.Save mock is already set by Expect")
	}

	if mmSave.defaultExpectation.paramPtrs == nil {
		mmSave.defaultExpectation.paramPtrs = &ReservedStockRepositoryMockSaveParamPtrs{}
	}
	mmSave.defaultExpectation.paramPtrs.ctx = &ctx

	return mmSave
}

// ExpectRStockParam2 sets up expected param rStock for ReservedStockRepository.Save
func (mmSave *mReservedStockRepositoryMockSave) ExpectRStockParam2(rStock *model.ReservedStock) *mReservedStockRepositoryMockSave {
	if mmSave.mock.funcSave != nil {
		mmSave.mock.t.Fatalf("ReservedStockRepositoryMock.Save mock is already set by Set")
	}

	if mmSave.defaultExpectation == nil {
		mmSave.defaultExpectation = &ReservedStockRepositoryMockSaveExpectation{}
	}

	if mmSave.defaultExpectation.params != nil {
		mmSave.mock.t.Fatalf("ReservedStockRepositoryMock.Save mock is already set by Expect")
	}

	if mmSave.defaultExpectation.paramPtrs == nil {
		mmSave.defaultExpectation.paramPtrs = &ReservedStockRepositoryMockSaveParamPtrs{}
	}
	mmSave.defaultExpectation.paramPtrs.rStock = &rStock

	return mmSave
}

// Inspect accepts an inspector function that has same arguments as the ReservedStockRepository.Save
func (mmSave *mReservedStockRepositoryMockSave) Inspect(f func(ctx context.Context, rStock *model.ReservedStock)) *mReservedStockRepositoryMockSave {
	if mmSave.mock.inspectFuncSave != nil {
		mmSave.mock.t.Fatalf("Inspect function is already set for ReservedStockRepositoryMock.Save")
	}

	mmSave.mock.inspectFuncSave = f

	return mmSave
}

// Return sets up results that will be returned by ReservedStockRepository.Save
func (mmSave *mReservedStockRepositoryMockSave) Return(err error) *ReservedStockRepositoryMock {
	if mmSave.mock.funcSave != nil {
		mmSave.mock.t.Fatalf("ReservedStockRepositoryMock.Save mock is already set by Set")
	}

	if mmSave.defaultExpectation == nil {
		mmSave.defaultExpectation = &ReservedStockRepositoryMockSaveExpectation{mock: mmSave.mock}
	}
	mmSave.defaultExpectation.results = &ReservedStockRepositoryMockSaveResults{err}
	return mmSave.mock
}

// Set uses given function f to mock the ReservedStockRepository.Save method
func (mmSave *mReservedStockRepositoryMockSave) Set(f func(ctx context.Context, rStock *model.ReservedStock) (err error)) *ReservedStockRepositoryMock {
	if mmSave.defaultExpectation != nil {
		mmSave.mock.t.Fatalf("Default expectation is already set for the ReservedStockRepository.Save method")
	}

	if len(mmSave.expectations) > 0 {
		mmSave.mock.t.Fatalf("Some expectations are already set for the ReservedStockRepository.Save method")
	}

	mmSave.mock.funcSave = f
	return mmSave.mock
}

// When sets expectation for the ReservedStockRepository.Save which will trigger the result defined by the following
// Then helper
func (mmSave *mReservedStockRepositoryMockSave) When(ctx context.Context, rStock *model.ReservedStock) *ReservedStockRepositoryMockSaveExpectation {
	if mmSave.mock.funcSave != nil {
		mmSave.mock.t.Fatalf("ReservedStockRepositoryMock.Save mock is already set by Set")
	}

	expectation := &ReservedStockRepositoryMockSaveExpectation{
		mock:   mmSave.mock,
		params: &ReservedStockRepositoryMockSaveParams{ctx, rStock},
	}
	mmSave.expectations = append(mmSave.expectations, expectation)
	return expectation
}

// Then sets up ReservedStockRepository.Save return parameters for the expectation previously defined by the When method
func (e *ReservedStockRepositoryMockSaveExpectation) Then(err error) *ReservedStockRepositoryMock {
	e.results = &ReservedStockRepositoryMockSaveResults{err}
	return e.mock
}

// Times sets number of times ReservedStockRepository.Save should be invoked
func (mmSave *mReservedStockRepositoryMockSave) Times(n uint64) *mReservedStockRepositoryMockSave {
	if n == 0 {
		mmSave.mock.t.Fatalf("Times of ReservedStockRepositoryMock.Save mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmSave.expectedInvocations, n)
	return mmSave
}

func (mmSave *mReservedStockRepositoryMockSave) invocationsDone() bool {
	if len(mmSave.expectations) == 0 && mmSave.defaultExpectation == nil && mmSave.mock.funcSave == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmSave.mock.afterSaveCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmSave.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// Save implements stock.ReservedStockRepository
func (mmSave *ReservedStockRepositoryMock) Save(ctx context.Context, rStock *model.ReservedStock) (err error) {
	mm_atomic.AddUint64(&mmSave.beforeSaveCounter, 1)
	defer mm_atomic.AddUint64(&mmSave.afterSaveCounter, 1)

	if mmSave.inspectFuncSave != nil {
		mmSave.inspectFuncSave(ctx, rStock)
	}

	mm_params := ReservedStockRepositoryMockSaveParams{ctx, rStock}

	// Record call args
	mmSave.SaveMock.mutex.Lock()
	mmSave.SaveMock.callArgs = append(mmSave.SaveMock.callArgs, &mm_params)
	mmSave.SaveMock.mutex.Unlock()

	for _, e := range mmSave.SaveMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmSave.SaveMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmSave.SaveMock.defaultExpectation.Counter, 1)
		mm_want := mmSave.SaveMock.defaultExpectation.params
		mm_want_ptrs := mmSave.SaveMock.defaultExpectation.paramPtrs

		mm_got := ReservedStockRepositoryMockSaveParams{ctx, rStock}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmSave.t.Errorf("ReservedStockRepositoryMock.Save got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.rStock != nil && !minimock.Equal(*mm_want_ptrs.rStock, mm_got.rStock) {
				mmSave.t.Errorf("ReservedStockRepositoryMock.Save got unexpected parameter rStock, want: %#v, got: %#v%s\n", *mm_want_ptrs.rStock, mm_got.rStock, minimock.Diff(*mm_want_ptrs.rStock, mm_got.rStock))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmSave.t.Errorf("ReservedStockRepositoryMock.Save got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmSave.SaveMock.defaultExpectation.results
		if mm_results == nil {
			mmSave.t.Fatal("No results are set for the ReservedStockRepositoryMock.Save")
		}
		return (*mm_results).err
	}
	if mmSave.funcSave != nil {
		return mmSave.funcSave(ctx, rStock)
	}
	mmSave.t.Fatalf("Unexpected call to ReservedStockRepositoryMock.Save. %v %v", ctx, rStock)
	return
}

// SaveAfterCounter returns a count of finished ReservedStockRepositoryMock.Save invocations
func (mmSave *ReservedStockRepositoryMock) SaveAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSave.afterSaveCounter)
}

// SaveBeforeCounter returns a count of ReservedStockRepositoryMock.Save invocations
func (mmSave *ReservedStockRepositoryMock) SaveBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSave.beforeSaveCounter)
}

// Calls returns a list of arguments used in each call to ReservedStockRepositoryMock.Save.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmSave *mReservedStockRepositoryMockSave) Calls() []*ReservedStockRepositoryMockSaveParams {
	mmSave.mutex.RLock()

	argCopy := make([]*ReservedStockRepositoryMockSaveParams, len(mmSave.callArgs))
	copy(argCopy, mmSave.callArgs)

	mmSave.mutex.RUnlock()

	return argCopy
}

// MinimockSaveDone returns true if the count of the Save invocations corresponds
// the number of defined expectations
func (m *ReservedStockRepositoryMock) MinimockSaveDone() bool {
	if m.SaveMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.SaveMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.SaveMock.invocationsDone()
}

// MinimockSaveInspect logs each unmet expectation
func (m *ReservedStockRepositoryMock) MinimockSaveInspect() {
	for _, e := range m.SaveMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ReservedStockRepositoryMock.Save with params: %#v", *e.params)
		}
	}

	afterSaveCounter := mm_atomic.LoadUint64(&m.afterSaveCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.SaveMock.defaultExpectation != nil && afterSaveCounter < 1 {
		if m.SaveMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ReservedStockRepositoryMock.Save")
		} else {
			m.t.Errorf("Expected call to ReservedStockRepositoryMock.Save with params: %#v", *m.SaveMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSave != nil && afterSaveCounter < 1 {
		m.t.Error("Expected call to ReservedStockRepositoryMock.Save")
	}

	if !m.SaveMock.invocationsDone() && afterSaveCounter > 0 {
		m.t.Errorf("Expected %d calls to ReservedStockRepositoryMock.Save but found %d calls",
			mm_atomic.LoadUint64(&m.SaveMock.expectedInvocations), afterSaveCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ReservedStockRepositoryMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetLockedInspect()

			m.MinimockSaveInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ReservedStockRepositoryMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *ReservedStockRepositoryMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetLockedDone() &&
		m.MinimockSaveDone()
}
