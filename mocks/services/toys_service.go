// Code generated by MockGen. DO NOT EDIT.
// Source: services.go
//
// Generated by this command:
//
//	mockgen -source=services.go -destination=../../mocks/services/toys_service.go -exclude_interfaces=MastersService,CategoriesService,TagsService,SsoService -package=mockservices
//

// Package mockservices is a generated GoMock package.
package mockservices

import (
	context "context"
	reflect "reflect"

	entities "github.com/DKhorkov/hmtm-toys/internal/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockToysService is a mock of ToysService interface.
type MockToysService struct {
	ctrl     *gomock.Controller
	recorder *MockToysServiceMockRecorder
	isgomock struct{}
}

// MockToysServiceMockRecorder is the mock recorder for MockToysService.
type MockToysServiceMockRecorder struct {
	mock *MockToysService
}

// NewMockToysService creates a new mock instance.
func NewMockToysService(ctrl *gomock.Controller) *MockToysService {
	mock := &MockToysService{ctrl: ctrl}
	mock.recorder = &MockToysServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockToysService) EXPECT() *MockToysServiceMockRecorder {
	return m.recorder
}

// AddToy mocks base method.
func (m *MockToysService) AddToy(ctx context.Context, toyData entities.AddToyDTO) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToy", ctx, toyData)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddToy indicates an expected call of AddToy.
func (mr *MockToysServiceMockRecorder) AddToy(ctx, toyData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToy", reflect.TypeOf((*MockToysService)(nil).AddToy), ctx, toyData)
}

// DeleteToy mocks base method.
func (m *MockToysService) DeleteToy(ctx context.Context, id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteToy", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteToy indicates an expected call of DeleteToy.
func (mr *MockToysServiceMockRecorder) DeleteToy(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteToy", reflect.TypeOf((*MockToysService)(nil).DeleteToy), ctx, id)
}

// GetAllToys mocks base method.
func (m *MockToysService) GetAllToys(ctx context.Context) ([]entities.Toy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllToys", ctx)
	ret0, _ := ret[0].([]entities.Toy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllToys indicates an expected call of GetAllToys.
func (mr *MockToysServiceMockRecorder) GetAllToys(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllToys", reflect.TypeOf((*MockToysService)(nil).GetAllToys), ctx)
}

// GetMasterToys mocks base method.
func (m *MockToysService) GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMasterToys", ctx, masterID)
	ret0, _ := ret[0].([]entities.Toy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMasterToys indicates an expected call of GetMasterToys.
func (mr *MockToysServiceMockRecorder) GetMasterToys(ctx, masterID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMasterToys", reflect.TypeOf((*MockToysService)(nil).GetMasterToys), ctx, masterID)
}

// GetToyByID mocks base method.
func (m *MockToysService) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetToyByID", ctx, id)
	ret0, _ := ret[0].(*entities.Toy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetToyByID indicates an expected call of GetToyByID.
func (mr *MockToysServiceMockRecorder) GetToyByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetToyByID", reflect.TypeOf((*MockToysService)(nil).GetToyByID), ctx, id)
}

// UpdateToy mocks base method.
func (m *MockToysService) UpdateToy(ctx context.Context, toyData entities.UpdateToyDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateToy", ctx, toyData)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateToy indicates an expected call of UpdateToy.
func (mr *MockToysServiceMockRecorder) UpdateToy(ctx, toyData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateToy", reflect.TypeOf((*MockToysService)(nil).UpdateToy), ctx, toyData)
}
