// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/htchan/BookSpider/internal/service (interfaces: Service)

// Package mockservice is a generated GoMock package.
package mockservice

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/htchan/BookSpider/internal/model"
	repo "github.com/htchan/BookSpider/internal/repo"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Book mocks base method.
func (m *MockService) Book(arg0 context.Context, arg1, arg2 string) (*model.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Book", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Book indicates an expected call of Book.
func (mr *MockServiceMockRecorder) Book(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Book", reflect.TypeOf((*MockService)(nil).Book), arg0, arg1, arg2)
}

// BookChapters mocks base method.
func (m *MockService) BookChapters(arg0 context.Context, arg1 *model.Book) (model.Chapters, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BookChapters", arg0, arg1)
	ret0, _ := ret[0].(model.Chapters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BookChapters indicates an expected call of BookChapters.
func (mr *MockServiceMockRecorder) BookChapters(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BookChapters", reflect.TypeOf((*MockService)(nil).BookChapters), arg0, arg1)
}

// BookContent mocks base method.
func (m *MockService) BookContent(arg0 context.Context, arg1 *model.Book) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BookContent", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BookContent indicates an expected call of BookContent.
func (mr *MockServiceMockRecorder) BookContent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BookContent", reflect.TypeOf((*MockService)(nil).BookContent), arg0, arg1)
}

// BookGroup mocks base method.
func (m *MockService) BookGroup(arg0 context.Context, arg1, arg2 string) (*model.Book, *model.BookGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BookGroup", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.Book)
	ret1, _ := ret[1].(*model.BookGroup)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// BookGroup indicates an expected call of BookGroup.
func (mr *MockServiceMockRecorder) BookGroup(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BookGroup", reflect.TypeOf((*MockService)(nil).BookGroup), arg0, arg1, arg2)
}

// BookInfo mocks base method.
func (m *MockService) BookInfo(arg0 context.Context, arg1 *model.Book) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BookInfo", arg0, arg1)
	ret0, _ := ret[0].(string)
	return ret0
}

// BookInfo indicates an expected call of BookInfo.
func (mr *MockServiceMockRecorder) BookInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BookInfo", reflect.TypeOf((*MockService)(nil).BookInfo), arg0, arg1)
}

// CheckAvailability mocks base method.
func (m *MockService) CheckAvailability(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAvailability", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckAvailability indicates an expected call of CheckAvailability.
func (mr *MockServiceMockRecorder) CheckAvailability(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAvailability", reflect.TypeOf((*MockService)(nil).CheckAvailability), arg0)
}

// DBStats mocks base method.
func (m *MockService) DBStats(arg0 context.Context) sql.DBStats {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DBStats", arg0)
	ret0, _ := ret[0].(sql.DBStats)
	return ret0
}

// DBStats indicates an expected call of DBStats.
func (mr *MockServiceMockRecorder) DBStats(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DBStats", reflect.TypeOf((*MockService)(nil).DBStats), arg0)
}

// Download mocks base method.
func (m *MockService) Download(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Download indicates an expected call of Download.
func (mr *MockServiceMockRecorder) Download(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockService)(nil).Download), arg0)
}

// DownloadBook mocks base method.
func (m *MockService) DownloadBook(arg0 context.Context, arg1 *model.Book) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadBook", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadBook indicates an expected call of DownloadBook.
func (mr *MockServiceMockRecorder) DownloadBook(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadBook", reflect.TypeOf((*MockService)(nil).DownloadBook), arg0, arg1)
}

// Explore mocks base method.
func (m *MockService) Explore(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Explore", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Explore indicates an expected call of Explore.
func (mr *MockServiceMockRecorder) Explore(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Explore", reflect.TypeOf((*MockService)(nil).Explore), arg0)
}

// ExploreBook mocks base method.
func (m *MockService) ExploreBook(arg0 context.Context, arg1 *model.Book) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExploreBook", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExploreBook indicates an expected call of ExploreBook.
func (mr *MockServiceMockRecorder) ExploreBook(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExploreBook", reflect.TypeOf((*MockService)(nil).ExploreBook), arg0, arg1)
}

// Name mocks base method.
func (m *MockService) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockServiceMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockService)(nil).Name))
}

// PatchDownloadStatus mocks base method.
func (m *MockService) PatchDownloadStatus(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchDownloadStatus", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PatchDownloadStatus indicates an expected call of PatchDownloadStatus.
func (mr *MockServiceMockRecorder) PatchDownloadStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchDownloadStatus", reflect.TypeOf((*MockService)(nil).PatchDownloadStatus), arg0)
}

// PatchMissingRecords mocks base method.
func (m *MockService) PatchMissingRecords(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchMissingRecords", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PatchMissingRecords indicates an expected call of PatchMissingRecords.
func (mr *MockServiceMockRecorder) PatchMissingRecords(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchMissingRecords", reflect.TypeOf((*MockService)(nil).PatchMissingRecords), arg0)
}

// Process mocks base method.
func (m *MockService) Process(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockServiceMockRecorder) Process(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockService)(nil).Process), arg0)
}

// ProcessBook mocks base method.
func (m *MockService) ProcessBook(arg0 context.Context, arg1 *model.Book) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessBook", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessBook indicates an expected call of ProcessBook.
func (mr *MockServiceMockRecorder) ProcessBook(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessBook", reflect.TypeOf((*MockService)(nil).ProcessBook), arg0, arg1)
}

// QueryBooks mocks base method.
func (m *MockService) QueryBooks(arg0 context.Context, arg1, arg2 string, arg3, arg4 int) ([]model.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryBooks", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]model.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryBooks indicates an expected call of QueryBooks.
func (mr *MockServiceMockRecorder) QueryBooks(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryBooks", reflect.TypeOf((*MockService)(nil).QueryBooks), arg0, arg1, arg2, arg3, arg4)
}

// RandomBooks mocks base method.
func (m *MockService) RandomBooks(arg0 context.Context, arg1 int) ([]model.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RandomBooks", arg0, arg1)
	ret0, _ := ret[0].([]model.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RandomBooks indicates an expected call of RandomBooks.
func (mr *MockServiceMockRecorder) RandomBooks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RandomBooks", reflect.TypeOf((*MockService)(nil).RandomBooks), arg0, arg1)
}

// Stats mocks base method.
func (m *MockService) Stats(arg0 context.Context) repo.Summary {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats", arg0)
	ret0, _ := ret[0].(repo.Summary)
	return ret0
}

// Stats indicates an expected call of Stats.
func (mr *MockServiceMockRecorder) Stats(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*MockService)(nil).Stats), arg0)
}

// Update mocks base method.
func (m *MockService) Update(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockServiceMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockService)(nil).Update), arg0)
}

// UpdateBook mocks base method.
func (m *MockService) UpdateBook(arg0 context.Context, arg1 *model.Book) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBook", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBook indicates an expected call of UpdateBook.
func (mr *MockServiceMockRecorder) UpdateBook(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBook", reflect.TypeOf((*MockService)(nil).UpdateBook), arg0, arg1)
}

// ValidateBookEnd mocks base method.
func (m *MockService) ValidateBookEnd(arg0 context.Context, arg1 *model.Book) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateBookEnd", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateBookEnd indicates an expected call of ValidateBookEnd.
func (mr *MockServiceMockRecorder) ValidateBookEnd(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateBookEnd", reflect.TypeOf((*MockService)(nil).ValidateBookEnd), arg0, arg1)
}

// ValidateEnd mocks base method.
func (m *MockService) ValidateEnd(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateEnd", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateEnd indicates an expected call of ValidateEnd.
func (mr *MockServiceMockRecorder) ValidateEnd(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateEnd", reflect.TypeOf((*MockService)(nil).ValidateEnd), arg0)
}
