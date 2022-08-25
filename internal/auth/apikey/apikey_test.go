package apikey

import (
	"errors"
	"keeper/internal/auth"
	"keeper/internal/mocks"
	"keeper/internal/models"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAPIKeyService_Authenticate(t *testing.T) {
	testUser := &models.User{
		ID: primitive.NewObjectID(),
	}
	testAPIKey := &models.APIKey{
		ID:        primitive.NewObjectID(),
		UserID:    testUser.ID,
		ExpiresAt: primitive.NewDateTimeFromTime(time.Now().Add(time.Hour)),
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		Hash:      "qiskLFPd-LcZ0y5hTufkyavp2Ky6LI1Sk_1yYcCOfP8=",
		MaskID:    "Ml7nXwRH3Nw3uX3x",
		Salt:      "JaOhInSZpNeq8DYNdGmfAxBl",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		credential *auth.Credential
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(mockAPIKeyRepo *mocks.MockIAPIKeyRepository, mockUserRepo *mocks.MockIUserRepository)
		want       *auth.AuthResponse
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_authenticate_apikey",
			args: args{
				credential: &auth.Credential{
					Type:   auth.CredentialTypeAPIKey,
					APIKey: "KP.Ml7nXwRH3Nw3uX3x.3ELBwprFqyNuAWKqd5dufxoeRRCgsvZ3grt1M9lKGr4NAwdr",
				},
			},
			stubFn: func(mockAPIKeyRepo *mocks.MockIAPIKeyRepository, mockUserRepo *mocks.MockIUserRepository) {
				// stubs
				mockAPIKeyRepo.EXPECT().FindAPIKeyByMaskID(gomock.Any()).Times(1).Return(testAPIKey, nil)
				mockUserRepo.EXPECT().FindUserById(gomock.Any()).Times(1).Return(testUser, nil)
			},
			wantErr:    false,
			wantErrMsg: "",
			want: &auth.AuthResponse{
				AuthMode: auth.CredentialTypeAPIKey,
				Credential: auth.Credential{
					Type:   auth.CredentialTypeAPIKey,
					APIKey: "KP.Ml7nXwRH3Nw3uX3x.3ELBwprFqyNuAWKqd5dufxoeRRCgsvZ3grt1M9lKGr4NAwdr",
				},
				User: testUser,
			},
		},
		{
			name: "should_fail_authenticate_apikey_does_not_exist",
			args: args{
				credential: &auth.Credential{
					Type:   auth.CredentialTypeAPIKey,
					APIKey: "KP.Ml7nXwRH3Nw3uX3x.3ELBwprFqyNuAWKqd5dufxoeRRCgsvZ3grt1M9lKGr4NAwdr",
				},
			},
			stubFn: func(mockAPIKeyRepo *mocks.MockIAPIKeyRepository, mockUserRepo *mocks.MockIUserRepository) {
				// stubs
				mockAPIKeyRepo.EXPECT().FindAPIKeyByMaskID(gomock.Any()).Times(1).Return(nil, errors.New("api key does not exist"))
			},
			wantErr:    true,
			wantErrMsg: ErrAPIKeyDoesNotExist.Error(),
			want:       nil,
		},
		{
			name: "should_fail_authenticate_apikey_cannot_decode_hash",
			args: args{
				credential: &auth.Credential{
					Type:   auth.CredentialTypeAPIKey,
					APIKey: "KP.Ml7nXwRH3Nw3uX3x.3ELBwprFqyNuAWKqd5dufxoeRRCgsvZ3grt1M9lKGr4NAwdr",
				},
			},
			stubFn: func(mockAPIKeyRepo *mocks.MockIAPIKeyRepository, mockUserRepo *mocks.MockIUserRepository) {
				// stubs
				mockAPIKeyRepo.EXPECT().FindAPIKeyByMaskID(gomock.Any()).Times(1).Return(&models.APIKey{
					Hash:   "abcde",
					MaskID: "Ml7nXwRH3Nw3uX3x",
					Salt:   "JaOhInSZpNeq8DYNdGmfAxBl",
				}, nil)
			},
			wantErr:    true,
			wantErrMsg: ErrFailedToDecodeAPIKeyHash.Error(),
			want:       nil,
		},
		{
			name: "should_fail_authenticate_apikey_invalid_credential_type",
			args: args{
				credential: &auth.Credential{
					Type: auth.CredentialTypeJWT,
				},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrAPIKeyCredentialType.Error(),
			want:       nil,
		},
		{
			name: "should_fail_authenticate_apikey_invalid_apikey_length",
			args: args{
				credential: &auth.Credential{
					Type:   auth.CredentialTypeAPIKey,
					APIKey: "KP.Ml7nXwRH3Nw3uX3x",
				},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrInvalidAPIKeyLength.Error(),
			want:       nil,
		},
		{
			name: "should_fail_authenticate_apikey_invalid_apikey",
			args: args{
				credential: &auth.Credential{
					Type:   auth.CredentialTypeAPIKey,
					APIKey: "KP.Ml7nXwRH3Nw3uX3x.abcdefg",
				},
			},
			stubFn: func(mockAPIKeyRepo *mocks.MockIAPIKeyRepository, mockUserRepo *mocks.MockIUserRepository) {
				// stubs
				mockAPIKeyRepo.EXPECT().FindAPIKeyByMaskID(gomock.Any()).Times(1).Return(testAPIKey, nil)
			},
			wantErr:    true,
			wantErrMsg: ErrInvalidAPIKey.Error(),
			want:       nil,
		},
		{
			name: "should_return_expired_apikey_error",
			args: args{
				credential: &auth.Credential{
					Type:   auth.CredentialTypeAPIKey,
					APIKey: "KP.Ml7nXwRH3Nw3uX3x.3ELBwprFqyNuAWKqd5dufxoeRRCgsvZ3grt1M9lKGr4NAwdr",
				},
			},
			stubFn: func(mockAPIKeyRepo *mocks.MockIAPIKeyRepository, mockUserRepo *mocks.MockIUserRepository) {
				// stubs
				testAPIKey.ExpiresAt = primitive.NewDateTimeFromTime(time.Now().Add(-time.Hour))
				mockAPIKeyRepo.EXPECT().FindAPIKeyByMaskID(gomock.Any()).Times(1).Return(testAPIKey, nil)
			},
			wantErr:    true,
			wantErrMsg: ErrExpiredAPIKey.Error(),
		},
		{
			name: "should_return_revoked_apikey_error",
			args: args{
				credential: &auth.Credential{
					Type:   auth.CredentialTypeAPIKey,
					APIKey: "KP.Ml7nXwRH3Nw3uX3x.3ELBwprFqyNuAWKqd5dufxoeRRCgsvZ3grt1M9lKGr4NAwdr",
				},
			},
			stubFn: func(mockAPIKeyRepo *mocks.MockIAPIKeyRepository, mockUserRepo *mocks.MockIUserRepository) {
				// stubs
				testAPIKey.ExpiresAt = primitive.NewDateTimeFromTime(time.Now().Add(time.Hour))
				testAPIKey.Revoked = true
				mockAPIKeyRepo.EXPECT().FindAPIKeyByMaskID(gomock.Any()).Times(1).Return(testAPIKey, nil)
			},
			wantErr:    true,
			wantErrMsg: ErrRevokedAPIKey.Error(),
		},
		{
			name: "should_fail_authenticate_user_not_found",
			args: args{
				credential: &auth.Credential{
					Type:   auth.CredentialTypeAPIKey,
					APIKey: "KP.Ml7nXwRH3Nw3uX3x.3ELBwprFqyNuAWKqd5dufxoeRRCgsvZ3grt1M9lKGr4NAwdr",
				},
			},
			stubFn: func(mockAPIKeyRepo *mocks.MockIAPIKeyRepository, mockUserRepo *mocks.MockIUserRepository) {
				// stubs
				testAPIKey.Revoked = false
				mockAPIKeyRepo.EXPECT().FindAPIKeyByMaskID(gomock.Any()).Times(1).Return(testAPIKey, nil)
				mockUserRepo.EXPECT().FindUserById(gomock.Any()).Times(1).Return(nil, models.ErrUserNotFound)
			},
			wantErr:    true,
			wantErrMsg: models.ErrUserNotFound.Error(),
			want:       nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockAPIKeyRepo := mocks.NewMockIAPIKeyRepository(ctrl)
			mockUserRepo := mocks.NewMockIUserRepository(ctrl)

			apiKeyChainSvc := NewAPIKeyService(mockAPIKeyRepo, mockUserRepo)
			if tc.stubFn != nil {
				tc.stubFn(mockAPIKeyRepo, mockUserRepo)
			}

			resp, err := apiKeyChainSvc.Authenticate(tc.args.credential)
			if tc.wantErr {
				require.Equal(t, tc.wantErrMsg, err.Error())
				return
			}
			require.Nil(t, err)
			require.Equal(t, tc.want, resp)
		})
	}
}
