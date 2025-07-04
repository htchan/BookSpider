// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/htchan/BookSpider/internal/vendorservice (interfaces: VendorService)
//
// Generated by this command:
//
//	mockgen -destination=../mock/vendorservice/vendor_service.go -package=mockvendorservice . VendorService
//

// Package mockvendorservice is a generated GoMock package.
package mockvendorservice

import (
	reflect "reflect"

	vendor "github.com/htchan/BookSpider/internal/vendorservice"
	gomock "go.uber.org/mock/gomock"
)

// MockVendorService is a mock of VendorService interface.
type MockVendorService struct {
	ctrl     *gomock.Controller
	recorder *MockVendorServiceMockRecorder
	isgomock struct{}
}

// MockVendorServiceMockRecorder is the mock recorder for MockVendorService.
type MockVendorServiceMockRecorder struct {
	mock *MockVendorService
}

// NewMockVendorService creates a new mock instance.
func NewMockVendorService(ctrl *gomock.Controller) *MockVendorService {
	mock := &MockVendorService{ctrl: ctrl}
	mock.recorder = &MockVendorServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVendorService) EXPECT() *MockVendorServiceMockRecorder {
	return m.recorder
}

// AvailabilityURL mocks base method.
func (m *MockVendorService) AvailabilityURL() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AvailabilityURL")
	ret0, _ := ret[0].(string)
	return ret0
}

// AvailabilityURL indicates an expected call of AvailabilityURL.
func (mr *MockVendorServiceMockRecorder) AvailabilityURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AvailabilityURL", reflect.TypeOf((*MockVendorService)(nil).AvailabilityURL))
}

// BookURL mocks base method.
func (m *MockVendorService) BookURL(bookID string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BookURL", bookID)
	ret0, _ := ret[0].(string)
	return ret0
}

// BookURL indicates an expected call of BookURL.
func (mr *MockVendorServiceMockRecorder) BookURL(bookID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BookURL", reflect.TypeOf((*MockVendorService)(nil).BookURL), bookID)
}

// ChapterListURL mocks base method.
func (m *MockVendorService) ChapterListURL(bookID string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChapterListURL", bookID)
	ret0, _ := ret[0].(string)
	return ret0
}

// ChapterListURL indicates an expected call of ChapterListURL.
func (mr *MockVendorServiceMockRecorder) ChapterListURL(bookID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChapterListURL", reflect.TypeOf((*MockVendorService)(nil).ChapterListURL), bookID)
}

// ChapterURL mocks base method.
func (m *MockVendorService) ChapterURL(resources ...string) string {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range resources {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ChapterURL", varargs...)
	ret0, _ := ret[0].(string)
	return ret0
}

// ChapterURL indicates an expected call of ChapterURL.
func (mr *MockVendorServiceMockRecorder) ChapterURL(resources ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChapterURL", reflect.TypeOf((*MockVendorService)(nil).ChapterURL), resources...)
}

// FindMissingIds mocks base method.
func (m *MockVendorService) FindMissingIds(ids []int) []int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMissingIds", ids)
	ret0, _ := ret[0].([]int)
	return ret0
}

// FindMissingIds indicates an expected call of FindMissingIds.
func (mr *MockVendorServiceMockRecorder) FindMissingIds(ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMissingIds", reflect.TypeOf((*MockVendorService)(nil).FindMissingIds), ids)
}

// IsAvailable mocks base method.
func (m *MockVendorService) IsAvailable(body string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAvailable", body)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsAvailable indicates an expected call of IsAvailable.
func (mr *MockVendorServiceMockRecorder) IsAvailable(body any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAvailable", reflect.TypeOf((*MockVendorService)(nil).IsAvailable), body)
}

// ParseBook mocks base method.
func (m *MockVendorService) ParseBook(body string) (*vendor.BookInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseBook", body)
	ret0, _ := ret[0].(*vendor.BookInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseBook indicates an expected call of ParseBook.
func (mr *MockVendorServiceMockRecorder) ParseBook(body any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseBook", reflect.TypeOf((*MockVendorService)(nil).ParseBook), body)
}

// ParseChapter mocks base method.
func (m *MockVendorService) ParseChapter(body string) (*vendor.ChapterInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseChapter", body)
	ret0, _ := ret[0].(*vendor.ChapterInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseChapter indicates an expected call of ParseChapter.
func (mr *MockVendorServiceMockRecorder) ParseChapter(body any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseChapter", reflect.TypeOf((*MockVendorService)(nil).ParseChapter), body)
}

// ParseChapterList mocks base method.
func (m *MockVendorService) ParseChapterList(bookID, body string) (vendor.ChapterList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseChapterList", bookID, body)
	ret0, _ := ret[0].(vendor.ChapterList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseChapterList indicates an expected call of ParseChapterList.
func (mr *MockVendorServiceMockRecorder) ParseChapterList(bookID, body any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseChapterList", reflect.TypeOf((*MockVendorService)(nil).ParseChapterList), bookID, body)
}
