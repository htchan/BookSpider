// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/htchan/BookSpider/internal/vendorservice (interfaces: VendorService)

// Package mockvendorservice is a generated GoMock package.
package mockvendorservice

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	vendor "github.com/htchan/BookSpider/internal/vendorservice"
)

// MockVendorService is a mock of VendorService interface.
type MockVendorService struct {
	ctrl     *gomock.Controller
	recorder *MockVendorServiceMockRecorder
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
func (m *MockVendorService) BookURL(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BookURL", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// BookURL indicates an expected call of BookURL.
func (mr *MockVendorServiceMockRecorder) BookURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BookURL", reflect.TypeOf((*MockVendorService)(nil).BookURL), arg0)
}

// ChapterListURL mocks base method.
func (m *MockVendorService) ChapterListURL(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChapterListURL", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// ChapterListURL indicates an expected call of ChapterListURL.
func (mr *MockVendorServiceMockRecorder) ChapterListURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChapterListURL", reflect.TypeOf((*MockVendorService)(nil).ChapterListURL), arg0)
}

// ChapterURL mocks base method.
func (m *MockVendorService) ChapterURL(arg0 ...string) string {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ChapterURL", varargs...)
	ret0, _ := ret[0].(string)
	return ret0
}

// ChapterURL indicates an expected call of ChapterURL.
func (mr *MockVendorServiceMockRecorder) ChapterURL(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChapterURL", reflect.TypeOf((*MockVendorService)(nil).ChapterURL), arg0...)
}

// FindMissingIds mocks base method.
func (m *MockVendorService) FindMissingIds(arg0 []int) []int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMissingIds", arg0)
	ret0, _ := ret[0].([]int)
	return ret0
}

// FindMissingIds indicates an expected call of FindMissingIds.
func (mr *MockVendorServiceMockRecorder) FindMissingIds(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMissingIds", reflect.TypeOf((*MockVendorService)(nil).FindMissingIds), arg0)
}

// IsAvailable mocks base method.
func (m *MockVendorService) IsAvailable(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAvailable", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsAvailable indicates an expected call of IsAvailable.
func (mr *MockVendorServiceMockRecorder) IsAvailable(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAvailable", reflect.TypeOf((*MockVendorService)(nil).IsAvailable), arg0)
}

// ParseBook mocks base method.
func (m *MockVendorService) ParseBook(arg0 string) (*vendor.BookInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseBook", arg0)
	ret0, _ := ret[0].(*vendor.BookInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseBook indicates an expected call of ParseBook.
func (mr *MockVendorServiceMockRecorder) ParseBook(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseBook", reflect.TypeOf((*MockVendorService)(nil).ParseBook), arg0)
}

// ParseChapter mocks base method.
func (m *MockVendorService) ParseChapter(arg0 string) (*vendor.ChapterInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseChapter", arg0)
	ret0, _ := ret[0].(*vendor.ChapterInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseChapter indicates an expected call of ParseChapter.
func (mr *MockVendorServiceMockRecorder) ParseChapter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseChapter", reflect.TypeOf((*MockVendorService)(nil).ParseChapter), arg0)
}

// ParseChapterList mocks base method.
func (m *MockVendorService) ParseChapterList(arg0, arg1 string) (vendor.ChapterList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseChapterList", arg0, arg1)
	ret0, _ := ret[0].(vendor.ChapterList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseChapterList indicates an expected call of ParseChapterList.
func (mr *MockVendorServiceMockRecorder) ParseChapterList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseChapterList", reflect.TypeOf((*MockVendorService)(nil).ParseChapterList), arg0, arg1)
}
