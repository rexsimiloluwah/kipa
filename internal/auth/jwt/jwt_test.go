package jwt

import (
	"errors"
	"keeper/internal/auth"
	"keeper/internal/config"
	"keeper/internal/mocks"
	"keeper/internal/models"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestJwtService_ParseExpiresIn(t *testing.T) {
	type args struct {
		expiresIn string
	}

	tt := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "should_successfully_parse_expires_in_duration",
			args: args{
				expiresIn: "15d",
			},
			want:    time.Now().Add(time.Duration(15) * time.Hour * 24).Unix(),
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			duration, err := parseExpiresInTime(tc.args.expiresIn)
			if tc.wantErr {
				require.NotNil(t, err)
				return
			}
			require.Equal(t, tc.want, duration)
		})
	}
}

func TestJwtService_GenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	jwtSvc := NewJwtService(&config.Config{}, mockUserRepo)

	type args struct {
		payload   map[string]interface{}
		expiresIn string
		secret    string
	}

	tt := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_generate_token",
			args: args{
				payload: map[string]interface{}{
					"email": "test@gmail.com",
					"id":    1,
				},
				secret:    "secret",
				expiresIn: "30m",
			},
			wantErr: false,
		},
		{
			name: "should_fail_generate_token_invalid_expires_in",
			args: args{
				payload: map[string]interface{}{
					"email": "test@gmail.com",
					"id":    1,
				},
				secret:    "secret",
				expiresIn: "30x",
			},
			wantErr:    true,
			wantErrMsg: ErrInvalidExpiresIn.Error(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			token, err := jwtSvc.GenerateToken(tc.args.payload, tc.args.expiresIn, tc.args.secret)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, tc.wantErrMsg, err.Error())
				return
			}

			require.Nil(t, err)
			require.Equal(t, len(strings.Split(token, ".")), 3)
		})
	}
}

func TestJwtService_GenerateAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		JwtSecretKey:            "secret",
		AccessTokenJwtExpiresIn: "30m",
	}
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	jwtSvc := NewJwtService(cfg, mockUserRepo)

	type jwtArgs struct {
		payload map[string]interface{}
	}

	testPayload := map[string]interface{}{"email": "test-user@gmail.com", "id": 2}

	tt := []struct {
		name       string
		args       jwtArgs
		stubFn     func()
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_generate_jwt_successfully",
			args: jwtArgs{
				payload: testPayload,
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "should_fail_invalid_expires_in",
			args: jwtArgs{
				payload: testPayload,
			},
			stubFn: func() {
				jwtSvc.cfg.AccessTokenJwtExpiresIn = "0"
			},
			wantErr:    true,
			wantErrMsg: ErrInvalidExpiresIn.Error(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn()
			}
			token, err := jwtSvc.GenerateAccessToken(tc.args.payload)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, tc.wantErrMsg, err.Error())
				return
			}

			require.Nil(t, err)
			require.Equal(t, len(strings.Split(token, ".")), 3)
		})
	}
}

func TestJwtService_GenerateRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		EmailVerificationTokenSecretKey: "secret",
		RefreshTokenJwtExpiresIn:        "7d",
	}
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	jwtSvc := NewJwtService(cfg, mockUserRepo)

	type jwtArgs struct {
		payload map[string]interface{}
	}

	testPayload := map[string]interface{}{"email": "test-user@gmail.com", "id": 2}

	tt := []struct {
		name       string
		args       jwtArgs
		stubFn     func()
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_generate_jwt_successfully",
			args: jwtArgs{
				payload: testPayload,
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "should_fail_invalid_expires_in",
			args: jwtArgs{
				payload: testPayload,
			},
			stubFn: func() {
				jwtSvc.cfg.RefreshTokenJwtExpiresIn = "0"
			},
			wantErr:    true,
			wantErrMsg: ErrInvalidExpiresIn.Error(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn()
			}
			token, err := jwtSvc.GenerateRefreshToken(tc.args.payload)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, tc.wantErrMsg, err.Error())
				return
			}

			require.Nil(t, err)
			require.Equal(t, len(strings.Split(token, ".")), 3)
		})
	}
}

func TestJwtService_GenerateEmailVerificationToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		JwtSecretKey:                    "secret",
		EmailVerificationTokenExpiresIn: "7d",
	}
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	jwtSvc := NewJwtService(cfg, mockUserRepo)

	type jwtArgs struct {
		payload map[string]interface{}
	}

	testPayload := map[string]interface{}{"email": "test-user@gmail.com", "id": 2}

	tt := []struct {
		name       string
		args       jwtArgs
		stubFn     func()
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_generate_jwt_successfully",
			args: jwtArgs{
				payload: testPayload,
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "should_fail_invalid_expires_in",
			args: jwtArgs{
				payload: testPayload,
			},
			stubFn: func() {
				jwtSvc.cfg.EmailVerificationTokenExpiresIn = "0"
			},
			wantErr:    true,
			wantErrMsg: ErrInvalidExpiresIn.Error(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn()
			}
			token, err := jwtSvc.GenerateEmailVerificationToken(tc.args.payload)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, tc.wantErrMsg, err.Error())
				return
			}

			require.Nil(t, err)
			require.Equal(t, len(strings.Split(token, ".")), 3)
		})
	}
}

func TestJwtService_GenerateResetPasswordToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		ResetPasswordTokenSecretKey: "secret",
		ResetPasswordTokenExpiresIn: "7d",
	}
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	jwtSvc := NewJwtService(cfg, mockUserRepo)

	type jwtArgs struct {
		payload map[string]interface{}
	}

	testPayload := map[string]interface{}{"email": "test-user@gmail.com", "id": 2}

	tt := []struct {
		name       string
		args       jwtArgs
		stubFn     func()
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_generate_jwt_successfully",
			args: jwtArgs{
				payload: testPayload,
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "should_fail_invalid_expires_in",
			args: jwtArgs{
				payload: testPayload,
			},
			stubFn: func() {
				jwtSvc.cfg.ResetPasswordTokenExpiresIn = "0"
			},
			wantErr:    true,
			wantErrMsg: ErrInvalidExpiresIn.Error(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn()
			}
			token, err := jwtSvc.GenerateResetPasswordToken(tc.args.payload)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, tc.wantErrMsg, err.Error())
				return
			}

			require.Nil(t, err)
			require.Equal(t, len(strings.Split(token, ".")), 3)
		})
	}
}

func TestJwtService_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		JwtSecretKey:            "secret",
		AccessTokenJwtExpiresIn: "30m",
	}
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	jwtSvc := NewJwtService(cfg, mockUserRepo)

	type args struct {
		tokenString string
	}

	// test token string
	payload := map[string]interface{}{
		"sub": 1,
	}
	testTokenString, _ := jwtSvc.GenerateAccessToken(payload)

	tt := []struct {
		name       string
		args       args
		want       *jwt.Token
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_fail_on_invalid_token",
			args: args{
				tokenString: "abcde",
			},
			want:       &jwt.Token{},
			wantErr:    true,
			wantErrMsg: "token contains an invalid number of segments",
		},
		{
			name: "should_successfully_validate_token",
			args: args{
				tokenString: testTokenString,
			},
			want: &jwt.Token{
				Valid: true,
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			token, err := jwtSvc.ValidateToken(tc.args.tokenString)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}
			require.Nil(t, err)
			require.True(t, token.Valid)
		})
	}
}

func TestJwtService_DecodeToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		JwtSecretKey:            "secret",
		AccessTokenJwtExpiresIn: "15m",
	}
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	jwtSvc := NewJwtService(cfg, mockUserRepo)

	// test token string
	payload := map[string]interface{}{
		"sub": "1",
	}
	testTokenString, _ := jwtSvc.GenerateAccessToken(payload)

	type args struct {
		tokenString string
	}

	tt := []struct {
		name       string
		args       args
		want       *JwtCustomClaims
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_fail_on_invalid_token",
			args: args{
				tokenString: "abcde",
			},
			want:       &JwtCustomClaims{},
			wantErr:    true,
			wantErrMsg: "error parsing token: token contains an invalid number of segments",
		},
		{
			name: "should_successfully_decode_token",
			args: args{
				tokenString: testTokenString,
			},
			want: &JwtCustomClaims{
				Payload: payload,
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := jwtSvc.DecodeToken(tc.args.tokenString, cfg.JwtSecretKey)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}
			require.Nil(t, err)
			require.Equal(t, claims.Payload, tc.want.Payload)
		})
	}
}

func TestJwtService_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type args struct {
		credential *auth.Credential
	}

	cfg := &config.Config{
		JwtSecretKey:            "secret",
		AccessTokenJwtExpiresIn: "15m",
	}
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	jwtSvc := NewJwtService(cfg, mockUserRepo)

	testUser := &models.User{
		ID:    primitive.NewObjectID(),
		Email: "test-user@gmail.com",
	}

	// generate a test access token
	payload := map[string]interface{}{"id": testUser.ID.Hex(), "email": testUser.Email}
	accessToken, _ := jwtSvc.GenerateAccessToken(payload)

	tt := []struct {
		name       string
		args       args
		stubFn     func(userRepo *mocks.MockIUserRepository)
		want       *auth.AuthResponse
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_authenticate_user",
			args: args{
				credential: &auth.Credential{
					Type: auth.CredentialTypeJWT,
					JWT:  accessToken,
				},
			},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().FindUserById(gomock.Any()).
					Times(1).Return(testUser, nil)
			},
		},
		{
			name: "should_fail_user_does_not_exist",
			args: args{
				credential: &auth.Credential{
					Type: auth.CredentialTypeJWT,
					JWT:  accessToken,
				},
			},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().FindUserById(gomock.Any()).
					Times(1).Return(nil, errors.New("user does not exist"))
			},
			wantErr:    true,
			wantErrMsg: "user does not exist",
		},
		{
			name: "should_fail_invalid_credential_type",
			args: args{
				credential: &auth.Credential{
					Type: auth.CredentialTypeAPIKey,
				},
			},
			wantErr:    true,
			wantErrMsg: "credential must be of jwt type",
		},
	}

	for _, tc := range tt {
		if tc.stubFn != nil {
			tc.stubFn(mockUserRepo)
		}
		authResponse, err := jwtSvc.Authenticate(tc.args.credential)
		if tc.wantErr {
			require.NotNil(t, err)
			require.Equal(t, err.Error(), tc.wantErrMsg)
			return
		}
		require.Nil(t, err)
		require.Equal(t, authResponse.AuthMode, auth.CredentialTypeJWT)
	}
}
