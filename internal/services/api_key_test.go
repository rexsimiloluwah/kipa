package services

import (
	"errors"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/mocks"
	"keeper/internal/models"
	"keeper/internal/repository"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Provide the mock APIKey service
func provideMockAPIKeyService(mockAPIKeyRepo repository.IAPIKeyRepository) IAPIKeyService {
	mockCfg := &config.Config{}
	return NewAPIKeyService(mockCfg, mockAPIKeyRepo)
}

var testAPIKey models.APIKey = models.APIKey{
	ID:        primitive.NewObjectID(),
	UserID:    primitive.NewObjectID(),
	ExpiresAt: primitive.NewDateTimeFromTime(time.Now()),
	CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
}

func TestAPIKeyService_FindAPIKeyByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiKeyRepo := mocks.NewMockIAPIKeyRepository(ctrl)

	type args struct {
		id string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(apiKeyRepo *mocks.MockIAPIKeyRepository)
		want       *models.APIKey
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_get_apikey_by_id",
			args: args{
				id: "abcd",
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().FindAPIKeyByID(gomock.Any()).
					Times(1).Return(&testAPIKey, nil)
			},
			want:    &testAPIKey,
			wantErr: false,
		},
		{
			name: "should_fail_to_get_apikey_by_id_return_err",
			args: args{
				id: "abcd",
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().FindAPIKeyByID(gomock.Any()).
					Times(1).Return(nil, errors.New("failed to find api key"))
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "failed to find api key",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(apiKeyRepo)
			}

			apiKeySvc := provideMockAPIKeyService(apiKeyRepo)
			apiKey, err := apiKeySvc.FindAPIKeyByID(tc.args.id)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, apiKey.ID, tc.want.ID)
			require.Equal(t, apiKey.UserID, tc.want.UserID)
			require.Equal(t, apiKey.CreatedAt, tc.want.CreatedAt)
			require.Equal(t, apiKey.ExpiresAt, tc.want.ExpiresAt)
			require.Equal(t, apiKey.UpdatedAt, tc.want.UpdatedAt)
		})
	}
}

func TestAPIKeyService_FindUserAPIKeys(t *testing.T) {
	testAPIKeys := []models.APIKey{
		testAPIKey,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiKeyRepo := mocks.NewMockIAPIKeyRepository(ctrl)

	type args struct {
		id string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(apiKeyRepo *mocks.MockIAPIKeyRepository)
		want       []models.APIKey
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_get_user_apikeys",
			args: args{
				id: "abcd",
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().FindUserAPIKeys(gomock.Any()).
					Times(1).Return(testAPIKeys, nil)
			},
			want:    testAPIKeys,
			wantErr: false,
		},
		{
			name: "should_fail_to_get_user_apikeys_by_id_return_err",
			args: args{
				id: "abcd",
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().FindUserAPIKeys(gomock.Any()).
					Times(1).Return(nil, errors.New("failed to find user api keys"))
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "failed to find user api keys",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(apiKeyRepo)
			}

			apiKeySvc := provideMockAPIKeyService(apiKeyRepo)
			apiKeys, err := apiKeySvc.FindUserAPIKeys(tc.args.id)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, len(apiKeys), len(tc.want))
			require.Equal(t, apiKeys[0].ID, tc.want[0].ID)
		})
	}
}

func TestAPIKeyService_CreateAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiKeyRepo := mocks.NewMockIAPIKeyRepository(ctrl)

	type args struct {
		data   dto.CreateAPIKeyInputDTO
		userId primitive.ObjectID
	}

	data := dto.CreateAPIKeyInputDTO{
		Name:      "test",
		KeyType:   "",
		Role:      "",
		ExpiresAt: func() *time.Time { t := time.Now().Add(time.Hour); return &t }(),
	}

	testUserId := primitive.NewObjectID()

	tt := []struct {
		name       string
		args       args
		stubFn     func(apiKeyRepo *mocks.MockIAPIKeyRepository)
		want       dto.CreateAPIKeyOutputDTO
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_create_api_key",
			args: args{
				data:   data,
				userId: testUserId,
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().CreateAPIKey(gomock.Any()).
					Times(1).Return(testUserId, nil)
			},
			want: dto.CreateAPIKeyOutputDTO{
				Name:      data.Name,
				ExpiresAt: primitive.NewDateTimeFromTime(*data.ExpiresAt),
			},
			wantErr: false,
		},
		{
			name: "should_fail_to_create_apikey_return_err",
			args: args{
				data:   data,
				userId: primitive.NewObjectID(),
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().CreateAPIKey(gomock.Any()).
					Times(1).Return(primitive.ObjectID{}, errors.New("failed to create api key"))
			},
			want:       dto.CreateAPIKeyOutputDTO{},
			wantErr:    true,
			wantErrMsg: "failed to create api key",
		},
		{
			name: "should_fail_to_create_apikey_invalid_expires_at",
			args: args{
				data: dto.CreateAPIKeyInputDTO{
					Name:      "test",
					KeyType:   "",
					Role:      "",
					ExpiresAt: func() *time.Time { t := time.Now().Add(-time.Hour); return &t }(),
				},
			},
			stubFn:     nil,
			want:       dto.CreateAPIKeyOutputDTO{},
			wantErr:    true,
			wantErrMsg: "api key expires_at cannot be before now",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(apiKeyRepo)
			}

			apiKeySvc := provideMockAPIKeyService(apiKeyRepo)
			out, err := apiKeySvc.CreateAPIKey(tc.args.data, tc.args.userId)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, tc.want.Name, out.Name)
			require.Equal(t, tc.want.ExpiresAt, out.ExpiresAt)
		})
	}
}

func TestAPIKeyService_UpdateAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiKeyRepo := mocks.NewMockIAPIKeyRepository(ctrl)

	type args struct {
		id   string
		data dto.UpdateAPIKeyInputDTO
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(apiKeyRepo *mocks.MockIAPIKeyRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_update_apikey",
			args: args{
				data: dto.UpdateAPIKeyInputDTO{},
				id:   "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().UpdateAPIKey(gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_update_apikey",
			args: args{
				data: dto.UpdateAPIKeyInputDTO{},
				id:   "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().UpdateAPIKey(gomock.Any()).
					Times(1).Return(errors.New("failed to update api key"))
			},
			wantErr:    true,
			wantErrMsg: "failed to update api key",
		},
		{
			name: "should_fail_update_api_key_invalid_id",
			args: args{
				data: dto.UpdateAPIKeyInputDTO{},
				id:   "abcd",
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: models.ErrInvalidObjectID.Error(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(apiKeyRepo)
			}

			apiKeySvc := provideMockAPIKeyService(apiKeyRepo)
			err := apiKeySvc.UpdateAPIKey(tc.args.id, tc.args.data)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}

func TestAPIKeyService_RevokeAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiKeyRepo := mocks.NewMockIAPIKeyRepository(ctrl)

	type args struct {
		id string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(apiKeyRepo *mocks.MockIAPIKeyRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_revoke_apikey",
			args: args{
				id: "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().RevokeAPIKey(gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_revoke_apikey",
			args: args{
				id: "abcd",
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().RevokeAPIKey(gomock.Any()).
					Times(1).Return(errors.New("error revoking api key"))
			},
			wantErr:    true,
			wantErrMsg: "error revoking api key",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(apiKeyRepo)
			}

			apiKeySvc := provideMockAPIKeyService(apiKeyRepo)
			err := apiKeySvc.RevokeAPIKey(tc.args.id)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}

func TestAPIKeyService_RevokeAPIKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiKeyRepo := mocks.NewMockIAPIKeyRepository(ctrl)

	type args struct {
		id []string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(apiKeyRepo *mocks.MockIAPIKeyRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_revoke_apikeys",
			args: args{
				id: []string{"62fa734bfc1cdb7f06a3bf6f"},
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().RevokeAPIKeys(gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_revoke_apikey",
			args: args{
				id: []string{"abcd"},
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().RevokeAPIKeys(gomock.Any()).
					Times(1).Return(errors.New("error revoking api keys"))
			},
			wantErr:    true,
			wantErrMsg: "error revoking api keys",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(apiKeyRepo)
			}

			apiKeySvc := provideMockAPIKeyService(apiKeyRepo)
			err := apiKeySvc.RevokeAPIKeys(tc.args.id)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}

func TestAPIKeyService_DeleteAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiKeyRepo := mocks.NewMockIAPIKeyRepository(ctrl)

	type args struct {
		id string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(apiKeyRepo *mocks.MockIAPIKeyRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_delete_apikey",
			args: args{
				id: "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().DeleteAPIKey(gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_delete_apikey",
			args: args{
				id: "abcd",
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().DeleteAPIKey(gomock.Any()).
					Times(1).Return(errors.New("error deleting api key"))
			},
			wantErr:    true,
			wantErrMsg: "error deleting api key",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(apiKeyRepo)
			}

			apiKeySvc := provideMockAPIKeyService(apiKeyRepo)
			err := apiKeySvc.DeleteAPIKey(tc.args.id)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}

func TestAPIKeyService_DeleteAPIKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiKeyRepo := mocks.NewMockIAPIKeyRepository(ctrl)

	type args struct {
		ids []string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(apiKeyRepo *mocks.MockIAPIKeyRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_delete_apikeys",
			args: args{
				ids: []string{"62fa734bfc1cdb7f06a3bf6f"},
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().DeleteAPIKeys(gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_delete_apikeys",
			args: args{
				ids: []string{"abcd"},
			},
			stubFn: func(apiKeyRepo *mocks.MockIAPIKeyRepository) {
				apiKeyRepo.EXPECT().DeleteAPIKeys(gomock.Any()).
					Times(1).Return(errors.New("error deleting api keys"))
			},
			wantErr:    true,
			wantErrMsg: "error deleting api keys",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(apiKeyRepo)
			}

			apiKeySvc := provideMockAPIKeyService(apiKeyRepo)
			err := apiKeySvc.DeleteAPIKeys(tc.args.ids)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}
