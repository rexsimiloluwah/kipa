package jwt

import (
	"fmt"
	"keeper/internal/config"
	"keeper/internal/mocks"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

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
				fmt.Println(err)
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

	tt := []struct {
		name       string
		args       jwtArgs
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_generate_jwt_successfully",
			args: jwtArgs{
				payload: map[string]interface{}{"email": "test-user@gmail.com", "id": 2},
			},
			wantErr:    false,
			wantErrMsg: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
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
		JwtSecretKey:             "secret",
		RefreshTokenJwtExpiresIn: "7d",
	}
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	jwtSvc := NewJwtService(cfg, mockUserRepo)

	type jwtArgs struct {
		payload map[string]interface{}
	}

	tt := []struct {
		name       string
		args       jwtArgs
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_generate_jwt_successfully",
			args: jwtArgs{
				payload: map[string]interface{}{"email": "test-user@gmail.com", "id": 2},
			},
			wantErr:    false,
			wantErrMsg: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
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
		JwtSecretKey: "secret",
	}
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	jwtSvc := NewJwtService(cfg, mockUserRepo)

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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := jwtSvc.DecodeToken(tc.args.tokenString)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}
			require.Nil(t, err)
			require.Equal(t, claims, tc.want)
		})
	}
}

func TestJwtService_Authenticate(t *testing.T) {
	require.True(t, true)
}
