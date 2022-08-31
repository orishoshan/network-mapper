// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock_client is a generated GoMock package.
package mock_client

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	client "github.com/otterize/otternose/sniffer/pkg/client"
)

// MockMapperClient is a mock of MapperClient interface.
type MockMapperClient struct {
	ctrl     *gomock.Controller
	recorder *MockMapperClientMockRecorder
}

// MockMapperClientMockRecorder is the mock recorder for MockMapperClient.
type MockMapperClientMockRecorder struct {
	mock *MockMapperClient
}

// NewMockMapperClient creates a new mock instance.
func NewMockMapperClient(ctrl *gomock.Controller) *MockMapperClient {
	mock := &MockMapperClient{ctrl: ctrl}
	mock.recorder = &MockMapperClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMapperClient) EXPECT() *MockMapperClientMockRecorder {
	return m.recorder
}

// ReportCaptureResults mocks base method.
func (m *MockMapperClient) ReportCaptureResults(ctx context.Context, results client.CaptureResults) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReportCaptureResults", ctx, results)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReportCaptureResults indicates an expected call of ReportCaptureResults.
func (mr *MockMapperClientMockRecorder) ReportCaptureResults(ctx, results interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportCaptureResults", reflect.TypeOf((*MockMapperClient)(nil).ReportCaptureResults), ctx, results)
}

// ReportSocketScanResults mocks base method.
func (m *MockMapperClient) ReportSocketScanResults(ctx context.Context, results client.SocketScanResults) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReportSocketScanResults", ctx, results)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReportSocketScanResults indicates an expected call of ReportSocketScanResults.
func (mr *MockMapperClientMockRecorder) ReportSocketScanResults(ctx, results interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportSocketScanResults", reflect.TypeOf((*MockMapperClient)(nil).ReportSocketScanResults), ctx, results)
}
