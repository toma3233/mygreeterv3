package server

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAccountsClient is a mock of AccountsClient interface
type MockAccountsClient struct {
	mock.Mock
}

// NewMockAccountsClient creates a new mock instance
func NewMockAccountsClient() *MockAccountsClient {
	return &MockAccountsClient{}
}

// BeginCreate implements mocked method
func (m *MockAccountsClient) BeginCreate(ctx context.Context, resourceGroupName string, accountName string, parameters armstorage.AccountCreateParameters, options *armstorage.AccountsClientBeginCreateOptions) (armstorage.AccountsClientBeginCreateResponse, error) {
	args := m.Called(ctx, resourceGroupName, accountName, parameters, options)
	return args.Get(0).(armstorage.AccountsClientBeginCreateResponse), args.Error(1)
}

// GetProperties implements mocked method
func (m *MockAccountsClient) GetProperties(ctx context.Context, resourceGroupName string, accountName string, options *armstorage.AccountsClientGetPropertiesOptions) (armstorage.AccountsClientGetPropertiesResponse, error) {
	args := m.Called(ctx, resourceGroupName, accountName, options)
	return args.Get(0).(armstorage.AccountsClientGetPropertiesResponse), args.Error(1)
}

// Delete implements mocked method
func (m *MockAccountsClient) Delete(ctx context.Context, resourceGroupName string, accountName string, options *armstorage.AccountsClientDeleteOptions) (armstorage.AccountsClientDeleteResponse, error) {
	args := m.Called(ctx, resourceGroupName, accountName, options)
	return args.Get(0).(armstorage.AccountsClientDeleteResponse), args.Error(1)
}

// Update implements mocked method
func (m *MockAccountsClient) Update(ctx context.Context, resourceGroupName string, accountName string, parameters armstorage.AccountUpdateParameters, options *armstorage.AccountsClientUpdateOptions) (armstorage.AccountsClientUpdateResponse, error) {
	args := m.Called(ctx, resourceGroupName, accountName, parameters, options)
	return args.Get(0).(armstorage.AccountsClientUpdateResponse), args.Error(1)
}

// NewListByResourceGroupPager implements mocked method
func (m *MockAccountsClient) NewListByResourceGroupPager(resourceGroupName string, options *armstorage.AccountsClientListByResourceGroupOptions) *armstorage.AccountsClientListByResourceGroupPager {
	args := m.Called(resourceGroupName, options)
	return args.Get(0).(*armstorage.AccountsClientListByResourceGroupPager)
}

// TestCreateStorageAccount tests the CreateStorageAccount method
func TestCreateStorageAccount(t *testing.T) {
	mockAccountsClient := NewMockAccountsClient()
	server := Server{AccountsClient: mockAccountsClient}
	ctx := context.Background()
	rgName := "testRG"
	region := "westus"
	accountName := "testAccount"
	parameters := armstorage.AccountCreateParameters{
		Location: to.Ptr(region),
		SKU: &armstorage.SKU{
			Name: to.Ptr(armstorage.SKUNameStandardGRS),
		},
		Kind: to.Ptr(armstorage.KindStorageV2),
	}

	mockAccountsClient.On("BeginCreate", ctx, rgName, accountName, parameters, mock.Anything).Return(armstorage.AccountsClientBeginCreateResponse{}, nil)

	_, err := server.CreateStorageAccount(ctx, &pb.CreateStorageAccountRequest{RgName: rgName, Region: region})
	assert.NoError(t, err)
	mockAccountsClient.AssertExpectations(t)
}

// TestReadStorageAccount tests the ReadStorageAccount method
func TestReadStorageAccount(t *testing.T) {
	mockAccountsClient := NewMockAccountsClient()
	server := Server{AccountsClient: mockAccountsClient}
	ctx := context.Background()
	rgName := "testRG"
	accountName := "testAccount"

	mockAccountsClient.On("GetProperties", ctx, rgName, accountName, mock.Anything).Return(armstorage.AccountsClientGetPropertiesResponse{}, nil)

	_, err := server.ReadStorageAccount(ctx, &pb.ReadStorageAccountRequest{RgName: rgName, SaName: accountName})
	assert.NoError(t, err)
	mockAccountsClient.AssertExpectations(t)
}

// TestUpdateStorageAccount tests the UpdateStorageAccount method
func TestUpdateStorageAccount(t *testing.T) {
	mockAccountsClient := NewMockAccountsClient()
	server := Server{AccountsClient: mockAccountsClient}
	ctx := context.Background()
	rgName := "testRG"
	accountName := "testAccount"
	parameters := armstorage.AccountUpdateParameters{}

	mockAccountsClient.On("Update", ctx, rgName, accountName, parameters, mock.Anything).Return(armstorage.AccountsClientUpdateResponse{}, nil)

	_, err := server.UpdateStorageAccount(ctx, &pb.UpdateStorageAccountRequest{RgName: rgName, SaName: accountName})
	assert.NoError(t, err)
	mockAccountsClient.AssertExpectations(t)
}

// TestDeleteStorageAccount tests the DeleteStorageAccount method
func TestDeleteStorageAccount(t *testing.T) {
	mockAccountsClient := NewMockAccountsClient()
	server := Server{AccountsClient: mockAccountsClient}
	ctx := context.Background()
	rgName := "testRG"
	accountName := "testAccount"

	mockAccountsClient.On("Delete", ctx, rgName, accountName, mock.Anything).Return(armstorage.AccountsClientDeleteResponse{}, nil)

	_, err := server.DeleteStorageAccount(ctx, &pb.DeleteStorageAccountRequest{RgName: rgName, SaName: accountName})
	assert.NoError(t, err)
	mockAccountsClient.AssertExpectations(t)
}

// TestListStorageAccounts tests the ListStorageAccounts method
func TestListStorageAccounts(t *testing.T) {
	mockAccountsClient := NewMockAccountsClient()
	server := Server{AccountsClient: mockAccountsClient}
	ctx := context.Background()
	rgName := "testRG"

	mockPager := &armstorage.AccountsClientListByResourceGroupPager{}
	mockAccountsClient.On("NewListByResourceGroupPager", rgName, mock.Anything).Return(mockPager)

	_, err := server.ListStorageAccounts(ctx, &pb.ListStorageAccountRequest{RgName: rgName})
	assert.NoError(t, err)
	mockAccountsClient.AssertExpectations(t)
}
