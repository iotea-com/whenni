package secrets

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func NewMockClient() *MockClient {
	mockClient := &MockClient{}

	return mockClient
}

func (m *MockClient) ReadKeyValue(ctx context.Context, path string) (map[string]any, error) {
	args := m.Called(ctx, path)
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *MockClient) WriteKeyValue(ctx context.Context, path string, data map[string]any) error {
	args := m.Called(ctx, path, data)
	return args.Error(1)
}

func (m *MockClient) DeleteKeyValue(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(1)
}

func (m *MockClient) RetrieveCert(serialNumber string) (*string, *string, error) {
	args := m.Called(serialNumber)
	return args.Get(0).(*string), args.Get(1).(*string), args.Error(2)
}

func (m *MockClient) RetrieveRootCaCert() (*string, error) {
	args := m.Called()
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockClient) CreateCert(certRequest CertRequest) (string, string, string, error) {
	// args := m.Called()
	return "test", "test", "test", nil
}

func (m *MockClient) RevokeCert(path string) error {
	args := m.Called(path)
	return args.Error(1)
}

func (m *MockClient) GetBoolFromMap(secrets map[string]any, key string) bool {
	args := m.Called()
	return args.Get(0).(bool)
}

func (m *MockClient) GetStringFromMap(secrets map[string]any, key string) string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockClient) GetIntFromMap(secrets map[string]any, key string) int {
	args := m.Called()
	return args.Get(0).(int)
}
