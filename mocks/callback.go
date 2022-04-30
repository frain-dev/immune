// Code generated by MockGen. DO NOT EDIT.
// Source: callback.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	immune "github.com/frain-dev/immune"
	gomock "github.com/golang/mock/gomock"
)

// MockCallbackServer is a mock of CallbackServer interface.
type MockCallbackServer struct {
	ctrl     *gomock.Controller
	recorder *MockCallbackServerMockRecorder
}

// MockCallbackServerMockRecorder is the mock recorder for MockCallbackServer.
type MockCallbackServerMockRecorder struct {
	mock *MockCallbackServer
}

// NewMockCallbackServer creates a new mock instance.
func NewMockCallbackServer(ctrl *gomock.Controller) *MockCallbackServer {
	mock := &MockCallbackServer{ctrl: ctrl}
	mock.recorder = &MockCallbackServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCallbackServer) EXPECT() *MockCallbackServerMockRecorder {
	return m.recorder
}

// ReceiveCallback mocks base method.
func (m *MockCallbackServer) ReceiveCallback(rc chan<- *immune.Signal) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReceiveCallback", rc)
}

// ReceiveCallback indicates an expected call of ReceiveCallback.
func (mr *MockCallbackServerMockRecorder) ReceiveCallback(rc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceiveCallback", reflect.TypeOf((*MockCallbackServer)(nil).ReceiveCallback), rc)
}

// Start mocks base method.
func (m *MockCallbackServer) Start(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockCallbackServerMockRecorder) Start(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockCallbackServer)(nil).Start), ctx)
}

// Stop mocks base method.
func (m *MockCallbackServer) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockCallbackServerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockCallbackServer)(nil).Stop))
}

// MockCallbackSignatureVerifier is a mock of CallbackSignatureVerifier interface.
type MockCallbackSignatureVerifier struct {
	ctrl     *gomock.Controller
	recorder *MockCallbackSignatureVerifierMockRecorder
}

// MockCallbackSignatureVerifierMockRecorder is the mock recorder for MockCallbackSignatureVerifier.
type MockCallbackSignatureVerifierMockRecorder struct {
	mock *MockCallbackSignatureVerifier
}

// NewMockCallbackSignatureVerifier creates a new mock instance.
func NewMockCallbackSignatureVerifier(ctrl *gomock.Controller) *MockCallbackSignatureVerifier {
	mock := &MockCallbackSignatureVerifier{ctrl: ctrl}
	mock.recorder = &MockCallbackSignatureVerifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCallbackSignatureVerifier) EXPECT() *MockCallbackSignatureVerifierMockRecorder {
	return m.recorder
}

// VerifyCallbackSignature mocks base method.
func (m *MockCallbackSignatureVerifier) VerifyCallbackSignature(s *immune.Signal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyCallbackSignature", s)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyCallbackSignature indicates an expected call of VerifyCallbackSignature.
func (mr *MockCallbackSignatureVerifierMockRecorder) VerifyCallbackSignature(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyCallbackSignature", reflect.TypeOf((*MockCallbackSignatureVerifier)(nil).VerifyCallbackSignature), s)
}
