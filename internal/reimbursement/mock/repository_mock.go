// Code generated by MockGen. DO NOT EDIT.
// Source: internal/reimbursement/repository.go
//
// Generated by this command:
//
//	mockgen -source internal/reimbursement/repository.go -destination internal/reimbursement/mock/repository_mock.go -package=mocks -typed
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	reimbursement "github.com/vnnyx/employee-management/internal/reimbursement"
	entity "github.com/vnnyx/employee-management/internal/reimbursement/entity"
	database "github.com/vnnyx/employee-management/pkg/database"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
	isgomock struct{}
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// FindReimbursementByPeriod mocks base method.
func (m *MockRepository) FindReimbursementByPeriod(ctx context.Context, startDate, endDate time.Time, opts ...entity.FindReimbursementOptions) (entity.FindReimbursementResult, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, startDate, endDate}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindReimbursementByPeriod", varargs...)
	ret0, _ := ret[0].(entity.FindReimbursementResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindReimbursementByPeriod indicates an expected call of FindReimbursementByPeriod.
func (mr *MockRepositoryMockRecorder) FindReimbursementByPeriod(ctx, startDate, endDate any, opts ...any) *MockRepositoryFindReimbursementByPeriodCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, startDate, endDate}, opts...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindReimbursementByPeriod", reflect.TypeOf((*MockRepository)(nil).FindReimbursementByPeriod), varargs...)
	return &MockRepositoryFindReimbursementByPeriodCall{Call: call}
}

// MockRepositoryFindReimbursementByPeriodCall wrap *gomock.Call
type MockRepositoryFindReimbursementByPeriodCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryFindReimbursementByPeriodCall) Return(arg0 entity.FindReimbursementResult, arg1 error) *MockRepositoryFindReimbursementByPeriodCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryFindReimbursementByPeriodCall) Do(f func(context.Context, time.Time, time.Time, ...entity.FindReimbursementOptions) (entity.FindReimbursementResult, error)) *MockRepositoryFindReimbursementByPeriodCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryFindReimbursementByPeriodCall) DoAndReturn(f func(context.Context, time.Time, time.Time, ...entity.FindReimbursementOptions) (entity.FindReimbursementResult, error)) *MockRepositoryFindReimbursementByPeriodCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// FindReimbursementByUserIDPeriod mocks base method.
func (m *MockRepository) FindReimbursementByUserIDPeriod(ctx context.Context, userID string, startDate, endDate time.Time) ([]entity.Reimbursement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindReimbursementByUserIDPeriod", ctx, userID, startDate, endDate)
	ret0, _ := ret[0].([]entity.Reimbursement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindReimbursementByUserIDPeriod indicates an expected call of FindReimbursementByUserIDPeriod.
func (mr *MockRepositoryMockRecorder) FindReimbursementByUserIDPeriod(ctx, userID, startDate, endDate any) *MockRepositoryFindReimbursementByUserIDPeriodCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindReimbursementByUserIDPeriod", reflect.TypeOf((*MockRepository)(nil).FindReimbursementByUserIDPeriod), ctx, userID, startDate, endDate)
	return &MockRepositoryFindReimbursementByUserIDPeriodCall{Call: call}
}

// MockRepositoryFindReimbursementByUserIDPeriodCall wrap *gomock.Call
type MockRepositoryFindReimbursementByUserIDPeriodCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryFindReimbursementByUserIDPeriodCall) Return(arg0 []entity.Reimbursement, arg1 error) *MockRepositoryFindReimbursementByUserIDPeriodCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryFindReimbursementByUserIDPeriodCall) Do(f func(context.Context, string, time.Time, time.Time) ([]entity.Reimbursement, error)) *MockRepositoryFindReimbursementByUserIDPeriodCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryFindReimbursementByUserIDPeriodCall) DoAndReturn(f func(context.Context, string, time.Time, time.Time) ([]entity.Reimbursement, error)) *MockRepositoryFindReimbursementByUserIDPeriodCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// StoreNewReimbursement mocks base method.
func (m *MockRepository) StoreNewReimbursement(ctx context.Context, arg1 entity.Reimbursement) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreNewReimbursement", ctx, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreNewReimbursement indicates an expected call of StoreNewReimbursement.
func (mr *MockRepositoryMockRecorder) StoreNewReimbursement(ctx, arg1 any) *MockRepositoryStoreNewReimbursementCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreNewReimbursement", reflect.TypeOf((*MockRepository)(nil).StoreNewReimbursement), ctx, arg1)
	return &MockRepositoryStoreNewReimbursementCall{Call: call}
}

// MockRepositoryStoreNewReimbursementCall wrap *gomock.Call
type MockRepositoryStoreNewReimbursementCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryStoreNewReimbursementCall) Return(arg0 error) *MockRepositoryStoreNewReimbursementCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryStoreNewReimbursementCall) Do(f func(context.Context, entity.Reimbursement) error) *MockRepositoryStoreNewReimbursementCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryStoreNewReimbursementCall) DoAndReturn(f func(context.Context, entity.Reimbursement) error) *MockRepositoryStoreNewReimbursementCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WithTx mocks base method.
func (m *MockRepository) WithTx(tx database.DBTx) reimbursement.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", tx)
	ret0, _ := ret[0].(reimbursement.Repository)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockRepositoryMockRecorder) WithTx(tx any) *MockRepositoryWithTxCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockRepository)(nil).WithTx), tx)
	return &MockRepositoryWithTxCall{Call: call}
}

// MockRepositoryWithTxCall wrap *gomock.Call
type MockRepositoryWithTxCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepositoryWithTxCall) Return(arg0 reimbursement.Repository) *MockRepositoryWithTxCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepositoryWithTxCall) Do(f func(database.DBTx) reimbursement.Repository) *MockRepositoryWithTxCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepositoryWithTxCall) DoAndReturn(f func(database.DBTx) reimbursement.Repository) *MockRepositoryWithTxCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
