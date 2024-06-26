// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/stone/go_project/ebook/ebook/cmd/cronjob/service/types.go

// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	domain "ebook/cmd/cronjob/domain"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCronJobService is a mock of CronJobService interface.
type MockCronJobService struct {
	ctrl     *gomock.Controller
	recorder *MockCronJobServiceMockRecorder
}

// MockCronJobServiceMockRecorder is the mock recorder for MockCronJobService.
type MockCronJobServiceMockRecorder struct {
	mock *MockCronJobService
}

// NewMockCronJobService creates a new mock instance.
func NewMockCronJobService(ctrl *gomock.Controller) *MockCronJobService {
	mock := &MockCronJobService{ctrl: ctrl}
	mock.recorder = &MockCronJobServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCronJobService) EXPECT() *MockCronJobServiceMockRecorder {
	return m.recorder
}

// AddJob mocks base method.
func (m *MockCronJobService) AddJob(ctx context.Context, j domain.CronJob) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddJob", ctx, j)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddJob indicates an expected call of AddJob.
func (mr *MockCronJobServiceMockRecorder) AddJob(ctx, j interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddJob", reflect.TypeOf((*MockCronJobService)(nil).AddJob), ctx, j)
}

// Preempt mocks base method.
func (m *MockCronJobService) Preempt(ctx context.Context) (domain.CronJob, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Preempt", ctx)
	ret0, _ := ret[0].(domain.CronJob)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Preempt indicates an expected call of Preempt.
func (mr *MockCronJobServiceMockRecorder) Preempt(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Preempt", reflect.TypeOf((*MockCronJobService)(nil).Preempt), ctx)
}

// ResetNextTime mocks base method.
func (m *MockCronJobService) ResetNextTime(ctx context.Context, j domain.CronJob) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetNextTime", ctx, j)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetNextTime indicates an expected call of ResetNextTime.
func (mr *MockCronJobServiceMockRecorder) ResetNextTime(ctx, j interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetNextTime", reflect.TypeOf((*MockCronJobService)(nil).ResetNextTime), ctx, j)
}
