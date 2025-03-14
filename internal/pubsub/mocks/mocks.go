// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	contract_client "github.com/dTelecom/pubsub-solana/internal/contract_client"
	solana "github.com/gagliardetto/solana-go"
	gomock "github.com/golang/mock/gomock"
)

// MockContractClient is a mock of ContractClient interface.
type MockContractClient struct {
	ctrl     *gomock.Controller
	recorder *MockContractClientMockRecorder
}

// MockContractClientMockRecorder is the mock recorder for MockContractClient.
type MockContractClientMockRecorder struct {
	mock *MockContractClient
}

// NewMockContractClient creates a new mock instance.
func NewMockContractClient(ctrl *gomock.Controller) *MockContractClient {
	mock := &MockContractClient{ctrl: ctrl}
	mock.recorder = &MockContractClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContractClient) EXPECT() *MockContractClientMockRecorder {
	return m.recorder
}

// IncomingMessageSubscribe mocks base method.
func (m *MockContractClient) IncomingMessageSubscribe(ctx context.Context, sender solana.PublicKey, handler func(context.Context, contract_client.MessageData)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IncomingMessageSubscribe", ctx, sender, handler)
}

// IncomingMessageSubscribe indicates an expected call of IncomingMessageSubscribe.
func (mr *MockContractClientMockRecorder) IncomingMessageSubscribe(ctx, sender, handler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncomingMessageSubscribe", reflect.TypeOf((*MockContractClient)(nil).IncomingMessageSubscribe), ctx, sender, handler)
}

// IsSigner mocks base method.
func (m *MockContractClient) IsSigner(key solana.PublicKey) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSigner", key)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSigner indicates an expected call of IsSigner.
func (mr *MockContractClientMockRecorder) IsSigner(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSigner", reflect.TypeOf((*MockContractClient)(nil).IsSigner), key)
}

// MarkAsRead mocks base method.
func (m *MockContractClient) MarkAsRead(ctx context.Context, sender solana.PublicKey, timestamp int64) (solana.Signature, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkAsRead", ctx, sender, timestamp)
	ret0, _ := ret[0].(solana.Signature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkAsRead indicates an expected call of MarkAsRead.
func (mr *MockContractClientMockRecorder) MarkAsRead(ctx, sender, timestamp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkAsRead", reflect.TypeOf((*MockContractClient)(nil).MarkAsRead), ctx, sender, timestamp)
}

// OutgoingMessageSubscribe mocks base method.
func (m *MockContractClient) OutgoingMessageSubscribe(ctx context.Context, receiver solana.PublicKey, handler func(context.Context, contract_client.MessageData)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OutgoingMessageSubscribe", ctx, receiver, handler)
}

// OutgoingMessageSubscribe indicates an expected call of OutgoingMessageSubscribe.
func (mr *MockContractClientMockRecorder) OutgoingMessageSubscribe(ctx, receiver, handler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OutgoingMessageSubscribe", reflect.TypeOf((*MockContractClient)(nil).OutgoingMessageSubscribe), ctx, receiver, handler)
}

// SendMessage mocks base method.
func (m *MockContractClient) SendMessage(ctx context.Context, receiver solana.PublicKey, content []byte) (solana.Signature, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", ctx, receiver, content)
	ret0, _ := ret[0].(solana.Signature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockContractClientMockRecorder) SendMessage(ctx, receiver, content interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockContractClient)(nil).SendMessage), ctx, receiver, content)
}

// MockMessageIdGenerator is a mock of MessageIdGenerator interface.
type MockMessageIdGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockMessageIdGeneratorMockRecorder
}

// MockMessageIdGeneratorMockRecorder is the mock recorder for MockMessageIdGenerator.
type MockMessageIdGeneratorMockRecorder struct {
	mock *MockMessageIdGenerator
}

// NewMockMessageIdGenerator creates a new mock instance.
func NewMockMessageIdGenerator(ctrl *gomock.Controller) *MockMessageIdGenerator {
	mock := &MockMessageIdGenerator{ctrl: ctrl}
	mock.recorder = &MockMessageIdGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageIdGenerator) EXPECT() *MockMessageIdGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockMessageIdGenerator) Generate() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate")
	ret0, _ := ret[0].(string)
	return ret0
}

// Generate indicates an expected call of Generate.
func (mr *MockMessageIdGeneratorMockRecorder) Generate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockMessageIdGenerator)(nil).Generate))
}

// MockDataEncoder is a mock of DataEncoder interface.
type MockDataEncoder struct {
	ctrl     *gomock.Controller
	recorder *MockDataEncoderMockRecorder
}

// MockDataEncoderMockRecorder is the mock recorder for MockDataEncoder.
type MockDataEncoderMockRecorder struct {
	mock *MockDataEncoder
}

// NewMockDataEncoder creates a new mock instance.
func NewMockDataEncoder(ctrl *gomock.Controller) *MockDataEncoder {
	mock := &MockDataEncoder{ctrl: ctrl}
	mock.recorder = &MockDataEncoderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataEncoder) EXPECT() *MockDataEncoderMockRecorder {
	return m.recorder
}

// Decode mocks base method.
func (m *MockDataEncoder) Decode(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decode", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decode indicates an expected call of Decode.
func (mr *MockDataEncoderMockRecorder) Decode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockDataEncoder)(nil).Decode), arg0)
}

// Encode mocks base method.
func (m *MockDataEncoder) Encode(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encode", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Encode indicates an expected call of Encode.
func (mr *MockDataEncoderMockRecorder) Encode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encode", reflect.TypeOf((*MockDataEncoder)(nil).Encode), arg0)
}
