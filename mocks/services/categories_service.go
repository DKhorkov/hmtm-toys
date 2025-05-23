// Code generated by MockGen. DO NOT EDIT.
// Source: services.go
//
// Generated by this command:
//
//	mockgen -source=services.go -destination=../../mocks/services/categories_service.go -exclude_interfaces=TagsService,MastersService,ToysService,SsoService -package=mockservices
//

// Package mockservices is a generated GoMock package.
package mockservices

import (
	context "context"
	reflect "reflect"

	entities "github.com/DKhorkov/hmtm-toys/internal/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockCategoriesService is a mock of CategoriesService interface.
type MockCategoriesService struct {
	ctrl     *gomock.Controller
	recorder *MockCategoriesServiceMockRecorder
	isgomock struct{}
}

// MockCategoriesServiceMockRecorder is the mock recorder for MockCategoriesService.
type MockCategoriesServiceMockRecorder struct {
	mock *MockCategoriesService
}

// NewMockCategoriesService creates a new mock instance.
func NewMockCategoriesService(ctrl *gomock.Controller) *MockCategoriesService {
	mock := &MockCategoriesService{ctrl: ctrl}
	mock.recorder = &MockCategoriesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCategoriesService) EXPECT() *MockCategoriesServiceMockRecorder {
	return m.recorder
}

// GetAllCategories mocks base method.
func (m *MockCategoriesService) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCategories", ctx)
	ret0, _ := ret[0].([]entities.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCategories indicates an expected call of GetAllCategories.
func (mr *MockCategoriesServiceMockRecorder) GetAllCategories(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCategories", reflect.TypeOf((*MockCategoriesService)(nil).GetAllCategories), ctx)
}

// GetCategoryByID mocks base method.
func (m *MockCategoriesService) GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategoryByID", ctx, id)
	ret0, _ := ret[0].(*entities.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategoryByID indicates an expected call of GetCategoryByID.
func (mr *MockCategoriesServiceMockRecorder) GetCategoryByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategoryByID", reflect.TypeOf((*MockCategoriesService)(nil).GetCategoryByID), ctx, id)
}
