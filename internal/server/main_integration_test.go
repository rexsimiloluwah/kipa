package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"keeper/internal/auth/jwt"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/server/testdb"
	mongoutils "keeper/pkg/mongo"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Repository struct {
	userRepo       repository.IUserRepository
	apiKeyRepo     repository.APIKeyRepository
	bucketRepo     repository.BucketRepository
	bucketItemRepo repository.BucketItemRepository
}

// this will be used to store suite-wide variables
// and register methods on the test suite
type ServerIntegrationTestSuite struct {
	suite.Suite
	Cfg    *config.Config
	DbConn mongoutils.Connection
	Server *Server
	Repo   *Repository
}

const (
	BASE_URL = "/api/v1"
)

// lifecycle methods in the suite
// this runs before all the tests in the suite are run
func (s *ServerIntegrationTestSuite) SetupSuite() {
	s.Cfg = config.NewTest()
	db := mongoutils.NewConnection(s.Cfg)
	s.DbConn = db
	s.Server = NewServer(s.Cfg, s.DbConn.Client)
	s.Server.RegisterRoutes()
	s.Repo = &Repository{
		userRepo: repository.NewUserRepository(s.Cfg, s.DbConn.Client),
	}
	s.DbConn.CleanDB(s.Cfg.DbName)
}

// this runs before each test in the suite is run
func (s *ServerIntegrationTestSuite) SetupTest() {
}

// this runs after all the tests in the suite are run
func (s *ServerIntegrationTestSuite) TearDownSuite() {
	// fmt.Println("this runs after all the tests in the suite are run")
	s.DbConn.CleanDB(s.Cfg.DbName)
	s.DbConn.Disconnect()
}

// this runs after each test in the suite is run
func (s *ServerIntegrationTestSuite) TearDownTest() {
}

// serialize response body into a struct
func serializeResponse(body *bytes.Buffer, out interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, out)
	return err
}

// serialize 'map[string]interface{}' output from response.Data into a struct
func serializeMap(data interface{}, out interface{}) error {
	// convert map to json
	jsonString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonString, out)
	return err
}

// Test case for registering a new user
func (s *ServerIntegrationTestSuite) TestAuth_RegisterUser() {
	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/auth/register", BASE_URL)
	body, err := json.Marshal(&dto.CreateUserInputDTO{
		Firstname: "bola",
		Lastname:  "tinubu",
		Email:     "bolatinubu@gmail.com",
		Username:  "bolatinubu",
		Password:  "Secret123!",
	})
	if err != nil {
		logrus.WithError(err).Error("error marshalling user data")
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		logrus.WithError(err).Fatal("error preparing request")
	}

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)
	// assert
	assert.Equal(s.T(), http.StatusCreated, recorder.Code)
}

// Test 'register existing user'
func (s *ServerIntegrationTestSuite) TestAuth_RegisterExistingUser() {
	newUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/auth/register", BASE_URL)
	body, err := json.Marshal(&dto.CreateUserInputDTO{
		Firstname: newUser.Firstname,
		Lastname:  newUser.Lastname,
		Email:     newUser.Email,
		Username:  newUser.Username,
		Password:  "Secret12345!",
	})
	assert.Nil(s.T(), err)

	// construct request
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json")
	assert.Nil(s.T(), err)

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)
	// assert
	assert.Equal(s.T(), http.StatusBadRequest, recorder.Code)
}

// Test 'login user'
func (s *ServerIntegrationTestSuite) TestAuth_LoginUser() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// arrange
	// construct the endpoint
	url := fmt.Sprintf("%s/auth/login", BASE_URL)
	body, err := json.Marshal(&dto.LoginUserInputDTO{
		Email:    testUser.Email,
		Password: "Secret12345!", //'secret' was used as the test password while seeding
	})
	assert.Nil(s.T(), err)

	// construct request
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json")
	assert.Nil(s.T(), err)

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)
	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
	// serialize body
	var response models.SuccessResponse
	err = serializeResponse(recorder.Body, &response)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), response.Data)

	// serialize token data
	var tokenResponse dto.LoginUserOutputDTO
	err = serializeMap(response.Data.(map[string]interface{}), &tokenResponse)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), tokenResponse.AccessToken)
	assert.NotEmpty(s.T(), tokenResponse.RefreshToken)
}

// Test 'fetch authenticated user'
func (s *ServerIntegrationTestSuite) TestAuth_GetAuthUser() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// generate access token
	jwtSvc := jwt.NewJwtService(s.Cfg, s.Repo.userRepo)
	accessToken, err := jwtSvc.GenerateAccessToken(map[string]interface{}{"email": testUser.Email, "id": testUser.ID})
	assert.Nil(s.T(), err)

	// arrange
	// construct the endpoint
	url := fmt.Sprintf("%s/auth/user", BASE_URL)
	request, err := http.NewRequest(http.MethodGet, url, nil)

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	assert.Nil(s.T(), err)

	// act: test getting authenticated user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
}

// Test 'refresh token'
func (s *ServerIntegrationTestSuite) TestAuth_RefreshToken() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// generate refresh token
	jwtSvc := jwt.NewJwtService(s.Cfg, s.Repo.userRepo)
	refreshToken, err := jwtSvc.GenerateRefreshToken(map[string]interface{}{"email": testUser.Email, "id": testUser.ID})
	assert.Nil(s.T(), err)

	// arrange
	// construct the endpoint
	url := fmt.Sprintf("%s/auth/refresh-token", BASE_URL)
	request, err := http.NewRequest(http.MethodPost, url, nil)

	request.Header.Add("x-refresh-token", refreshToken)
	assert.Nil(s.T(), err)

	// act: test getting authenticated user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
}

// Test 'find user'
func (s *ServerIntegrationTestSuite) TestUser_FindUser() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// arrange
	// construct the endpoint
	url := fmt.Sprintf("%s/users/%s", BASE_URL, testUser.ID.Hex())
	request, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(s.T(), err)

	// act: test getting authenticated user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	var response models.SuccessResponse
	err = serializeResponse(recorder.Body, &response)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), response.Data)

	// serialize user data
	var userData models.User
	err = serializeMap(response.Data.(map[string]interface{}), &userData)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), userData.ID.Hex(), testUser.ID.Hex())
	assert.Equal(s.T(), userData.Email, testUser.Email)
}

// Test 'find users'
func (s *ServerIntegrationTestSuite) TestUser_FindUsers() {
	_, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// arrange
	// construct the endpoint
	url := fmt.Sprintf("%s/users", BASE_URL)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(s.T(), err)

	// act: test getting authenticated user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	var response models.SuccessResponse
	err = serializeResponse(recorder.Body, &response)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), response.Data)
}

// Test 'update user'
func (s *ServerIntegrationTestSuite) TestUser_UpdateUser() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// generate access token
	jwtSvc := jwt.NewJwtService(s.Cfg, s.Repo.userRepo)
	accessToken, err := jwtSvc.GenerateAccessToken(map[string]interface{}{"email": testUser.Email, "id": testUser.ID})
	assert.Nil(s.T(), err)

	// arrange
	// construct the endpoint
	url := fmt.Sprintf("%s/user", BASE_URL)
	body, err := json.Marshal(&dto.UpdateUserInputDTO{
		Firstname: "updated-firstname",
		Username:  "updated-username",
	})
	assert.Nil(s.T(), err)
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Add("Content-Type", "application/json")
	assert.Nil(s.T(), err)

	// act: test getting authenticated user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	// check the user
	user, err := s.Repo.userRepo.FindUserByEmail(testUser.Email)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), user.Firstname, "updated-firstname")
	assert.Equal(s.T(), user.Username, "updated-username")
}

// Test 'delete user'
func (s *ServerIntegrationTestSuite) TestUser_DeleteUser() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// generate access token
	jwtSvc := jwt.NewJwtService(s.Cfg, s.Repo.userRepo)
	accessToken, err := jwtSvc.GenerateAccessToken(map[string]interface{}{"email": testUser.Email, "id": testUser.ID})
	assert.Nil(s.T(), err)

	// arrange
	// construct the endpoint
	url := fmt.Sprintf("%s/user", BASE_URL)
	request, err := http.NewRequest(http.MethodDelete, url, nil)

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	assert.Nil(s.T(), err)

	// act: test getting authenticated user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	// check the user
	_, err = s.Repo.userRepo.FindUserByEmail(testUser.Email)
	assert.Equal(s.T(), err, models.ErrUserNotFound)
}

// Test 'create api key'
func (s *ServerIntegrationTestSuite) TestAPIKey_Create() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// generate access token
	jwtSvc := jwt.NewJwtService(s.Cfg, s.Repo.userRepo)
	accessToken, err := jwtSvc.GenerateAccessToken(map[string]interface{}{"email": testUser.Email, "id": testUser.ID})
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/api_key", BASE_URL)
	body, err := json.Marshal(&dto.CreateAPIKeyInputDTO{
		Name:      "test-key",
		KeyType:   "test",
		ExpiresAt: func() *time.Time { t := time.Now().Add(time.Hour); return &t }(),
	})
	assert.Nil(s.T(), err)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	assert.Nil(s.T(), err)

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)
	// assert
	assert.Equal(s.T(), http.StatusCreated, recorder.Code)
}

// Test 'list users api keys'
func (s *ServerIntegrationTestSuite) TestAPIKey_FindUserAPIKeys() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	_, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/api_keys", BASE_URL)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))
	assert.Nil(s.T(), err)

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)
	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	var response models.SuccessResponse
	err = serializeResponse(recorder.Body, &response)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), response.Data)
	assert.Len(s.T(), response.Data, 1)
}

// Test 'update api key'
func (s *ServerIntegrationTestSuite) TestAPIKey_UpdateAPIKey() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	testAPIKey, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/api_key/%s", BASE_URL, testAPIKey.ID.Hex())
	body, err := json.Marshal(&dto.UpdateAPIKeyInputDTO{
		Name: "updated-key-name",
	})
	assert.Nil(s.T(), err)

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))
	request.Header.Add("Content-Type", "application/json")

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	// check the api key
	// apiKey, err := s.Repo.apiKeyRepo.FindAPIKeyByID(testAPIKey.ID.Hex())
	// assert.Nil(s.T(), err)
	// assert.Equal(s.T(), apiKey.Name, "updated-key-name")
}

// Test 'find api key by id'
// TODO: fix 'invalid memory address or nil pointer dereference error'

// Test 'delete api key'
func (s *ServerIntegrationTestSuite) TestAPIKey_Delete() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	testAPIKey, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/api_keys", BASE_URL)
	body, err := json.Marshal(&dto.APIKeysIDsInputDTO{
		Ids: []string{testAPIKey.ID.Hex()},
	})
	assert.Nil(s.T(), err)

	request, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(body))
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))
	request.Header.Add("Content-Type", "application/json")

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
}

// Test 'revoke api key'
func (s *ServerIntegrationTestSuite) TestAPIKey_Revoke() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	testAPIKey, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/api_keys/revoke", BASE_URL)
	body, err := json.Marshal(&dto.APIKeysIDsInputDTO{
		Ids: []string{testAPIKey.ID.Hex()},
	})
	assert.Nil(s.T(), err)

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))
	request.Header.Add("Content-Type", "application/json")

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
}

// Test 'create bucket'
func (s *ServerIntegrationTestSuite) TestBucket_Create() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	_, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/bucket", BASE_URL)
	body, err := json.Marshal(&dto.CreateBucketInputDTO{
		Name:        "test-bucket",
		Description: "test bucket",
	})
	assert.Nil(s.T(), err)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))
	request.Header.Add("Content-Type", "application/json")

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusCreated, recorder.Code)
}

// Test 'find bucket by uid'
func (s *ServerIntegrationTestSuite) TestBucket_FindBucketByUID() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// seed api key for authorization
	_, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket
	testBucket, err := testdb.SeedBucket(s.DbConn.Client, s.Cfg, testUser.ID, models.BUCKET_PERMISSIONS)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/bucket/%s", BASE_URL, testBucket.UID)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
	// serialize body
	var response models.SuccessResponse
	err = serializeResponse(recorder.Body, &response)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), response.Data)

	// serialize bucket data
	var bucketDetailsOutput dto.BucketDetailsOutput
	err = serializeMap(response.Data.(map[string]interface{}), &bucketDetailsOutput)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testBucket.ID.Hex(), bucketDetailsOutput.ID.Hex())
	assert.Equal(s.T(), testBucket.Name, bucketDetailsOutput.Name)
	assert.Len(s.T(), bucketDetailsOutput.BucketItems, 0)
}

// Test 'update bucket'
func (s *ServerIntegrationTestSuite) TestBucket_Update() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// seed api key for authorization
	_, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket
	testBucket, err := testdb.SeedBucket(s.DbConn.Client, s.Cfg, testUser.ID, models.BUCKET_PERMISSIONS)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/bucket/%s", BASE_URL, testBucket.UID)
	body, err := json.Marshal(&dto.UpdateBucketInputDTO{
		Name: "updated-bucket-name",
	})
	assert.Nil(s.T(), err)

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))
	request.Header.Add("Content-Type", "application/json")

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	// check bucket
	// bucket, err := s.Repo.bucketRepo.FindBucketByUID(testBucket.UID)
	// assert.Nil(s.T(), err)
	// assert.NotNil(s.T(), bucket)
	// assert.Equal(s.T(), bucket.Name, "updated-bucket-name")
}

// Test 'delete bucket'
func (s *ServerIntegrationTestSuite) TestBucket_Delete() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// seed api key for authorization
	_, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket
	testBucket, err := testdb.SeedBucket(s.DbConn.Client, s.Cfg, testUser.ID, models.BUCKET_PERMISSIONS)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/bucket/%s", BASE_URL, testBucket.UID)

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
}

// Test 'list user buckets'
func (s *ServerIntegrationTestSuite) TestBucket_ListUserBuckets() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// seed api key for authorization
	_, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket
	testBucket, err := testdb.SeedBucket(s.DbConn.Client, s.Cfg, testUser.ID, models.BUCKET_PERMISSIONS)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/buckets", BASE_URL)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	// serialize body
	var response models.SuccessResponse
	err = serializeResponse(recorder.Body, &response)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), response.Data)

	// serialize bucket data
	var bucketsDetailsOutput []dto.BucketDetailsOutput
	err = serializeMap(response.Data.([]interface{}), &bucketsDetailsOutput)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), bucketsDetailsOutput[0].ID.Hex(), testBucket.ID.Hex())
}

// Test 'create bucket item'
func (s *ServerIntegrationTestSuite) TestBucketItem_Create() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// seed api key for authorization
	_, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket
	testBucket, err := testdb.SeedBucket(s.DbConn.Client, s.Cfg, testUser.ID, models.BUCKET_PERMISSIONS)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/item/%s", BASE_URL, testBucket.UID)
	body, err := json.Marshal(&dto.CreateBucketItemInputDTO{
		Key:  "test-key",
		Data: "string value",
	})
	assert.Nil(s.T(), err)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))
	request.Header.Add("Content-Type", "application/json")

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusCreated, recorder.Code)
}

// Test 'list bucket items'
func (s *ServerIntegrationTestSuite) TestBucketItem_ListBucketItems() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// seed api key for authorization
	_, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket
	testBucket, err := testdb.SeedBucket(s.DbConn.Client, s.Cfg, testUser.ID, models.BUCKET_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket item
	testBucketItem, err := testdb.SeedBucketItem(
		s.DbConn.Client,
		s.Cfg, testUser.ID,
		testBucket.ID,
		testBucket.UID,
		"string",
	)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/items/%s", BASE_URL, testBucketItem.BucketUID)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	// serialize body
	var response models.SuccessResponse
	err = serializeResponse(recorder.Body, &response)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), response.Data)

	// serialize bucket data
	var bucketItems []models.BucketItem
	err = serializeMap(response.Data.([]interface{}), &bucketItems)
	assert.Nil(s.T(), err)
	assert.Len(s.T(), bucketItems, 1)
	assert.Equal(s.T(), bucketItems[0].ID.Hex(), testBucketItem.ID.Hex())
}

// Test 'find bucket item by key name'
func (s *ServerIntegrationTestSuite) TestBucketItem_FindByKeyName() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// seed api key for authorization
	_, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket
	testBucket, err := testdb.SeedBucket(s.DbConn.Client, s.Cfg, testUser.ID, models.BUCKET_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket item
	testBucketItem, err := testdb.SeedBucketItem(
		s.DbConn.Client,
		s.Cfg, testUser.ID,
		testBucket.ID,
		testBucket.UID, "string",
	)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/item/%s/%s?full=true", BASE_URL, testBucket.UID, testBucketItem.Key)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	// serialize body
	var response models.SuccessResponse
	err = serializeResponse(recorder.Body, &response)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), response.Data)

	// serialize bucket data
	var bucketItem models.BucketItem
	err = serializeMap(response.Data, &bucketItem)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), bucketItem.ID.Hex(), testBucketItem.ID.Hex())
}

// Test 'update bucket item by key name'
func (s *ServerIntegrationTestSuite) TestBucketItem_UpdateByKeyName() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// seed api key for authorization
	_, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket
	testBucket, err := testdb.SeedBucket(s.DbConn.Client, s.Cfg, testUser.ID, models.BUCKET_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket item
	testBucketItem, err := testdb.SeedBucketItem(
		s.DbConn.Client,
		s.Cfg, testUser.ID,
		testBucket.ID,
		testBucket.UID, "string",
	)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/item/%s/%s", BASE_URL, testBucket.UID, testBucketItem.Key)
	body, err := json.Marshal(&dto.UpdateBucketItemInputDTO{
		Data: "updated string value",
	})
	assert.Nil(s.T(), err)

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))
	request.Header.Add("Content-Type", "application/json")

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
}

// Test 'bucket item atomic operations' i.e. increment or decrement
func (s *ServerIntegrationTestSuite) TestBucketItem_IncrementOperation() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// seed api key for authorization
	_, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket
	testBucket, err := testdb.SeedBucket(s.DbConn.Client, s.Cfg, testUser.ID, models.BUCKET_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket item
	testBucketItem, err := testdb.SeedBucketItem(
		s.DbConn.Client,
		s.Cfg, testUser.ID,
		testBucket.ID,
		testBucket.UID, 20,
	)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/item/%s/%s", BASE_URL, testBucket.UID, testBucketItem.Key)

	request, err := http.NewRequest(http.MethodPut, url, strings.NewReader("+5"))
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))

	// act: test bucket item increment operation
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)
	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
}

// Test 'delete bucket item by keyname'
func (s *ServerIntegrationTestSuite) TestBucketItem_DeleteByKeyName() {
	testUser, err := testdb.SeedUser(s.DbConn.Client, s.Cfg)
	assert.Nil(s.T(), err)

	// seed api key for authorization
	_, key, err := testdb.SeedAPIKey(s.DbConn.Client, s.Cfg, testUser.ID, models.APIKEY_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket
	testBucket, err := testdb.SeedBucket(s.DbConn.Client, s.Cfg, testUser.ID, models.BUCKET_PERMISSIONS)
	assert.Nil(s.T(), err)

	// seed a new bucket item
	testBucketItem, err := testdb.SeedBucketItem(
		s.DbConn.Client,
		s.Cfg, testUser.ID,
		testBucket.ID,
		testBucket.UID, "string",
	)
	assert.Nil(s.T(), err)

	// construct the endpoint
	// arrange
	url := fmt.Sprintf("%s/item/%s/%s", BASE_URL, testBucket.UID, testBucketItem.Key)

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	assert.Nil(s.T(), err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))

	// act: test registering a new user
	// send request and record response
	recorder := httptest.NewRecorder()
	s.Server.Server.ServeHTTP(recorder, request)

	// assert
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestServerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ServerIntegrationTestSuite))
}
