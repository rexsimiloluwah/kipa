// Code generated by MockGen. DO NOT EDIT.
// Source: repository/main.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "keeper/internal/models"
	utils "keeper/internal/utils"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	bson "go.mongodb.org/mongo-driver/bson"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

// MockIUserRepository is a mock of IUserRepository interface.
type MockIUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIUserRepositoryMockRecorder
}

// MockIUserRepositoryMockRecorder is the mock recorder for MockIUserRepository.
type MockIUserRepositoryMockRecorder struct {
	mock *MockIUserRepository
}

// NewMockIUserRepository creates a new mock instance.
func NewMockIUserRepository(ctrl *gomock.Controller) *MockIUserRepository {
	mock := &MockIUserRepository{ctrl: ctrl}
	mock.recorder = &MockIUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserRepository) EXPECT() *MockIUserRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockIUserRepository) CreateUser(user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockIUserRepositoryMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockIUserRepository)(nil).CreateUser), user)
}

// DeleteUser mocks base method.
func (m *MockIUserRepository) DeleteUser(userId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockIUserRepositoryMockRecorder) DeleteUser(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockIUserRepository)(nil).DeleteUser), userId)
}

// FindAllUsers mocks base method.
func (m *MockIUserRepository) FindAllUsers() ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllUsers")
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllUsers indicates an expected call of FindAllUsers.
func (mr *MockIUserRepositoryMockRecorder) FindAllUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllUsers", reflect.TypeOf((*MockIUserRepository)(nil).FindAllUsers))
}

// FindUserByEmail mocks base method.
func (m *MockIUserRepository) FindUserByEmail(email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByEmail", email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByEmail indicates an expected call of FindUserByEmail.
func (mr *MockIUserRepositoryMockRecorder) FindUserByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByEmail", reflect.TypeOf((*MockIUserRepository)(nil).FindUserByEmail), email)
}

// FindUserById mocks base method.
func (m *MockIUserRepository) FindUserById(id string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserById", id)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserById indicates an expected call of FindUserById.
func (mr *MockIUserRepositoryMockRecorder) FindUserById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserById", reflect.TypeOf((*MockIUserRepository)(nil).FindUserById), id)
}

// UpdateUser mocks base method.
func (m *MockIUserRepository) UpdateUser(user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockIUserRepositoryMockRecorder) UpdateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockIUserRepository)(nil).UpdateUser), user)
}

// MockIAPIKeyRepository is a mock of IAPIKeyRepository interface.
type MockIAPIKeyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIAPIKeyRepositoryMockRecorder
}

// MockIAPIKeyRepositoryMockRecorder is the mock recorder for MockIAPIKeyRepository.
type MockIAPIKeyRepositoryMockRecorder struct {
	mock *MockIAPIKeyRepository
}

// NewMockIAPIKeyRepository creates a new mock instance.
func NewMockIAPIKeyRepository(ctrl *gomock.Controller) *MockIAPIKeyRepository {
	mock := &MockIAPIKeyRepository{ctrl: ctrl}
	mock.recorder = &MockIAPIKeyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAPIKeyRepository) EXPECT() *MockIAPIKeyRepositoryMockRecorder {
	return m.recorder
}

// CreateAPIKey mocks base method.
func (m *MockIAPIKeyRepository) CreateAPIKey(apiKey *models.APIKey) (primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAPIKey", apiKey)
	ret0, _ := ret[0].(primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAPIKey indicates an expected call of CreateAPIKey.
func (mr *MockIAPIKeyRepositoryMockRecorder) CreateAPIKey(apiKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAPIKey", reflect.TypeOf((*MockIAPIKeyRepository)(nil).CreateAPIKey), apiKey)
}

// DeleteAPIKey mocks base method.
func (m *MockIAPIKeyRepository) DeleteAPIKey(apiKeyID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAPIKey", apiKeyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAPIKey indicates an expected call of DeleteAPIKey.
func (mr *MockIAPIKeyRepositoryMockRecorder) DeleteAPIKey(apiKeyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAPIKey", reflect.TypeOf((*MockIAPIKeyRepository)(nil).DeleteAPIKey), apiKeyID)
}

// DeleteAPIKeys mocks base method.
func (m *MockIAPIKeyRepository) DeleteAPIKeys(apiKeyIDs []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAPIKeys", apiKeyIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAPIKeys indicates an expected call of DeleteAPIKeys.
func (mr *MockIAPIKeyRepositoryMockRecorder) DeleteAPIKeys(apiKeyIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAPIKeys", reflect.TypeOf((*MockIAPIKeyRepository)(nil).DeleteAPIKeys), apiKeyIDs)
}

// FindAPIKeyByHash mocks base method.
func (m *MockIAPIKeyRepository) FindAPIKeyByHash(hash string) (*models.APIKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAPIKeyByHash", hash)
	ret0, _ := ret[0].(*models.APIKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAPIKeyByHash indicates an expected call of FindAPIKeyByHash.
func (mr *MockIAPIKeyRepositoryMockRecorder) FindAPIKeyByHash(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAPIKeyByHash", reflect.TypeOf((*MockIAPIKeyRepository)(nil).FindAPIKeyByHash), hash)
}

// FindAPIKeyByID mocks base method.
func (m *MockIAPIKeyRepository) FindAPIKeyByID(apiKeyID string) (*models.APIKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAPIKeyByID", apiKeyID)
	ret0, _ := ret[0].(*models.APIKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAPIKeyByID indicates an expected call of FindAPIKeyByID.
func (mr *MockIAPIKeyRepositoryMockRecorder) FindAPIKeyByID(apiKeyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAPIKeyByID", reflect.TypeOf((*MockIAPIKeyRepository)(nil).FindAPIKeyByID), apiKeyID)
}

// FindAPIKeyByMaskID mocks base method.
func (m *MockIAPIKeyRepository) FindAPIKeyByMaskID(maskID string) (*models.APIKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAPIKeyByMaskID", maskID)
	ret0, _ := ret[0].(*models.APIKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAPIKeyByMaskID indicates an expected call of FindAPIKeyByMaskID.
func (mr *MockIAPIKeyRepositoryMockRecorder) FindAPIKeyByMaskID(maskID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAPIKeyByMaskID", reflect.TypeOf((*MockIAPIKeyRepository)(nil).FindAPIKeyByMaskID), maskID)
}

// FindUserAPIKeys mocks base method.
func (m *MockIAPIKeyRepository) FindUserAPIKeys(userID string) ([]models.APIKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserAPIKeys", userID)
	ret0, _ := ret[0].([]models.APIKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserAPIKeys indicates an expected call of FindUserAPIKeys.
func (mr *MockIAPIKeyRepositoryMockRecorder) FindUserAPIKeys(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserAPIKeys", reflect.TypeOf((*MockIAPIKeyRepository)(nil).FindUserAPIKeys), userID)
}

// RevokeAPIKey mocks base method.
func (m *MockIAPIKeyRepository) RevokeAPIKey(apiKeyID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevokeAPIKey", apiKeyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RevokeAPIKey indicates an expected call of RevokeAPIKey.
func (mr *MockIAPIKeyRepositoryMockRecorder) RevokeAPIKey(apiKeyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeAPIKey", reflect.TypeOf((*MockIAPIKeyRepository)(nil).RevokeAPIKey), apiKeyID)
}

// RevokeAPIKeys mocks base method.
func (m *MockIAPIKeyRepository) RevokeAPIKeys(apiKeyIDs []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevokeAPIKeys", apiKeyIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// RevokeAPIKeys indicates an expected call of RevokeAPIKeys.
func (mr *MockIAPIKeyRepositoryMockRecorder) RevokeAPIKeys(apiKeyIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeAPIKeys", reflect.TypeOf((*MockIAPIKeyRepository)(nil).RevokeAPIKeys), apiKeyIDs)
}

// UpdateAPIKey mocks base method.
func (m *MockIAPIKeyRepository) UpdateAPIKey(apiKey *models.APIKey) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAPIKey", apiKey)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAPIKey indicates an expected call of UpdateAPIKey.
func (mr *MockIAPIKeyRepositoryMockRecorder) UpdateAPIKey(apiKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAPIKey", reflect.TypeOf((*MockIAPIKeyRepository)(nil).UpdateAPIKey), apiKey)
}

// MockIBucketRepository is a mock of IBucketRepository interface.
type MockIBucketRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIBucketRepositoryMockRecorder
}

// MockIBucketRepositoryMockRecorder is the mock recorder for MockIBucketRepository.
type MockIBucketRepositoryMockRecorder struct {
	mock *MockIBucketRepository
}

// NewMockIBucketRepository creates a new mock instance.
func NewMockIBucketRepository(ctrl *gomock.Controller) *MockIBucketRepository {
	mock := &MockIBucketRepository{ctrl: ctrl}
	mock.recorder = &MockIBucketRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBucketRepository) EXPECT() *MockIBucketRepositoryMockRecorder {
	return m.recorder
}

// CreateBucket mocks base method.
func (m *MockIBucketRepository) CreateBucket(bucket *models.Bucket) (primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBucket", bucket)
	ret0, _ := ret[0].(primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBucket indicates an expected call of CreateBucket.
func (mr *MockIBucketRepositoryMockRecorder) CreateBucket(bucket interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBucket", reflect.TypeOf((*MockIBucketRepository)(nil).CreateBucket), bucket)
}

// DeleteBucketByID mocks base method.
func (m *MockIBucketRepository) DeleteBucketByID(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBucketByID", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBucketByID indicates an expected call of DeleteBucketByID.
func (mr *MockIBucketRepositoryMockRecorder) DeleteBucketByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBucketByID", reflect.TypeOf((*MockIBucketRepository)(nil).DeleteBucketByID), id)
}

// DeleteBucketByUID mocks base method.
func (m *MockIBucketRepository) DeleteBucketByUID(uid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBucketByUID", uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBucketByUID indicates an expected call of DeleteBucketByUID.
func (mr *MockIBucketRepositoryMockRecorder) DeleteBucketByUID(uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBucketByUID", reflect.TypeOf((*MockIBucketRepository)(nil).DeleteBucketByUID), uid)
}

// FindBucketByID mocks base method.
func (m *MockIBucketRepository) FindBucketByID(id string) (*models.Bucket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBucketByID", id)
	ret0, _ := ret[0].(*models.Bucket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBucketByID indicates an expected call of FindBucketByID.
func (mr *MockIBucketRepositoryMockRecorder) FindBucketByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBucketByID", reflect.TypeOf((*MockIBucketRepository)(nil).FindBucketByID), id)
}

// FindBucketByUID mocks base method.
func (m *MockIBucketRepository) FindBucketByUID(uid string) (*models.Bucket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBucketByUID", uid)
	ret0, _ := ret[0].(*models.Bucket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBucketByUID indicates an expected call of FindBucketByUID.
func (mr *MockIBucketRepositoryMockRecorder) FindBucketByUID(uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBucketByUID", reflect.TypeOf((*MockIBucketRepository)(nil).FindBucketByUID), uid)
}

// FindBucketsByUserID mocks base method.
func (m *MockIBucketRepository) FindBucketsByUserID(userID string) ([]models.Bucket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBucketsByUserID", userID)
	ret0, _ := ret[0].([]models.Bucket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBucketsByUserID indicates an expected call of FindBucketsByUserID.
func (mr *MockIBucketRepositoryMockRecorder) FindBucketsByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBucketsByUserID", reflect.TypeOf((*MockIBucketRepository)(nil).FindBucketsByUserID), userID)
}

// FindBucketsByUserIDPaged mocks base method.
func (m *MockIBucketRepository) FindBucketsByUserIDPaged(userID string, filter bson.M, findOpts *options.FindOptions, paginationParams utils.PaginationParams) ([]models.Bucket, utils.PageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBucketsByUserIDPaged", userID, filter, findOpts, paginationParams)
	ret0, _ := ret[0].([]models.Bucket)
	ret1, _ := ret[1].(utils.PageInfo)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindBucketsByUserIDPaged indicates an expected call of FindBucketsByUserIDPaged.
func (mr *MockIBucketRepositoryMockRecorder) FindBucketsByUserIDPaged(userID, filter, findOpts, paginationParams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBucketsByUserIDPaged", reflect.TypeOf((*MockIBucketRepository)(nil).FindBucketsByUserIDPaged), userID, filter, findOpts, paginationParams)
}

// UpdateBucket mocks base method.
func (m *MockIBucketRepository) UpdateBucket(bucket *models.Bucket) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBucket", bucket)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBucket indicates an expected call of UpdateBucket.
func (mr *MockIBucketRepositoryMockRecorder) UpdateBucket(bucket interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBucket", reflect.TypeOf((*MockIBucketRepository)(nil).UpdateBucket), bucket)
}

// MockIBucketItemRepository is a mock of IBucketItemRepository interface.
type MockIBucketItemRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIBucketItemRepositoryMockRecorder
}

// MockIBucketItemRepositoryMockRecorder is the mock recorder for MockIBucketItemRepository.
type MockIBucketItemRepositoryMockRecorder struct {
	mock *MockIBucketItemRepository
}

// NewMockIBucketItemRepository creates a new mock instance.
func NewMockIBucketItemRepository(ctrl *gomock.Controller) *MockIBucketItemRepository {
	mock := &MockIBucketItemRepository{ctrl: ctrl}
	mock.recorder = &MockIBucketItemRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBucketItemRepository) EXPECT() *MockIBucketItemRepositoryMockRecorder {
	return m.recorder
}

// CreateBucketItem mocks base method.
func (m *MockIBucketItemRepository) CreateBucketItem(bucketItem *models.BucketItem) (primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBucketItem", bucketItem)
	ret0, _ := ret[0].(primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBucketItem indicates an expected call of CreateBucketItem.
func (mr *MockIBucketItemRepositoryMockRecorder) CreateBucketItem(bucketItem interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBucketItem", reflect.TypeOf((*MockIBucketItemRepository)(nil).CreateBucketItem), bucketItem)
}

// DeleteBucketItemById mocks base method.
func (m *MockIBucketItemRepository) DeleteBucketItemById(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBucketItemById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBucketItemById indicates an expected call of DeleteBucketItemById.
func (mr *MockIBucketItemRepositoryMockRecorder) DeleteBucketItemById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBucketItemById", reflect.TypeOf((*MockIBucketItemRepository)(nil).DeleteBucketItemById), id)
}

// DeleteBucketItemByKeyName mocks base method.
func (m *MockIBucketItemRepository) DeleteBucketItemByKeyName(bucketUID, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBucketItemByKeyName", bucketUID, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBucketItemByKeyName indicates an expected call of DeleteBucketItemByKeyName.
func (mr *MockIBucketItemRepositoryMockRecorder) DeleteBucketItemByKeyName(bucketUID, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBucketItemByKeyName", reflect.TypeOf((*MockIBucketItemRepository)(nil).DeleteBucketItemByKeyName), bucketUID, key)
}

// DeleteBucketItems mocks base method.
func (m *MockIBucketItemRepository) DeleteBucketItems(bucketUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBucketItems", bucketUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBucketItems indicates an expected call of DeleteBucketItems.
func (mr *MockIBucketItemRepositoryMockRecorder) DeleteBucketItems(bucketUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBucketItems", reflect.TypeOf((*MockIBucketItemRepository)(nil).DeleteBucketItems), bucketUID)
}

// DeleteBucketItemsById mocks base method.
func (m *MockIBucketItemRepository) DeleteBucketItemsById(ids []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBucketItemsById", ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBucketItemsById indicates an expected call of DeleteBucketItemsById.
func (mr *MockIBucketItemRepositoryMockRecorder) DeleteBucketItemsById(ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBucketItemsById", reflect.TypeOf((*MockIBucketItemRepository)(nil).DeleteBucketItemsById), ids)
}

// FindBucketItemByID mocks base method.
func (m *MockIBucketItemRepository) FindBucketItemByID(id string) (*models.BucketItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBucketItemByID", id)
	ret0, _ := ret[0].(*models.BucketItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBucketItemByID indicates an expected call of FindBucketItemByID.
func (mr *MockIBucketItemRepositoryMockRecorder) FindBucketItemByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBucketItemByID", reflect.TypeOf((*MockIBucketItemRepository)(nil).FindBucketItemByID), id)
}

// FindBucketItemByKeyName mocks base method.
func (m *MockIBucketItemRepository) FindBucketItemByKeyName(bucketUID, key string) (*models.BucketItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBucketItemByKeyName", bucketUID, key)
	ret0, _ := ret[0].(*models.BucketItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBucketItemByKeyName indicates an expected call of FindBucketItemByKeyName.
func (mr *MockIBucketItemRepositoryMockRecorder) FindBucketItemByKeyName(bucketUID, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBucketItemByKeyName", reflect.TypeOf((*MockIBucketItemRepository)(nil).FindBucketItemByKeyName), bucketUID, key)
}

// FindBucketItems mocks base method.
func (m *MockIBucketItemRepository) FindBucketItems(bucketUID string) ([]models.BucketItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBucketItems", bucketUID)
	ret0, _ := ret[0].([]models.BucketItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBucketItems indicates an expected call of FindBucketItems.
func (mr *MockIBucketItemRepositoryMockRecorder) FindBucketItems(bucketUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBucketItems", reflect.TypeOf((*MockIBucketItemRepository)(nil).FindBucketItems), bucketUID)
}

// FindBucketItemsPaged mocks base method.
func (m *MockIBucketItemRepository) FindBucketItemsPaged(filter bson.M, opts *options.FindOptions, paginationParams utils.PaginationParams) ([]models.BucketItem, utils.PageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBucketItemsPaged", filter, opts, paginationParams)
	ret0, _ := ret[0].([]models.BucketItem)
	ret1, _ := ret[1].(utils.PageInfo)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindBucketItemsPaged indicates an expected call of FindBucketItemsPaged.
func (mr *MockIBucketItemRepositoryMockRecorder) FindBucketItemsPaged(filter, opts, paginationParams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBucketItemsPaged", reflect.TypeOf((*MockIBucketItemRepository)(nil).FindBucketItemsPaged), filter, opts, paginationParams)
}

// IncrementIntItem mocks base method.
func (m *MockIBucketItemRepository) IncrementIntItem(bucketUID, key string, amount int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementIntItem", bucketUID, key, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementIntItem indicates an expected call of IncrementIntItem.
func (mr *MockIBucketItemRepositoryMockRecorder) IncrementIntItem(bucketUID, key, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementIntItem", reflect.TypeOf((*MockIBucketItemRepository)(nil).IncrementIntItem), bucketUID, key, amount)
}

// UpdateBucketItem mocks base method.
func (m *MockIBucketItemRepository) UpdateBucketItem(bucketItem *models.BucketItem, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBucketItem", bucketItem, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBucketItem indicates an expected call of UpdateBucketItem.
func (mr *MockIBucketItemRepositoryMockRecorder) UpdateBucketItem(bucketItem, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBucketItem", reflect.TypeOf((*MockIBucketItemRepository)(nil).UpdateBucketItem), bucketItem, key)
}
