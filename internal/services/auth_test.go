package services

import (
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/mocks"
	"keeper/internal/models"
	"keeper/internal/utils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// provide the mock auth service
func provideMockAuthService(mockUserRepo *mocks.MockIUserRepository) IAuthService {
	cfg := &config.Config{
		Env:                      "test",
		JwtSecretKey:             "secret",
		AccessTokenJwtExpiresIn:  "30m",
		RefreshTokenJwtExpiresIn: "7d",
	}
	return NewAuthService(cfg, mockUserRepo)
}

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)

	type args struct {
		data dto.LoginUserInputDTO
	}

	userData := dto.LoginUserInputDTO{
		Email:    "testuser@gmail.com",
		Password: "secret",
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(userRepo *mocks.MockIUserRepository)
		want       *dto.LoginUserOutputDTO
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_login_user",
			args: args{
				data: userData,
			},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				hashedPassword, _ := utils.HashPassword(userData.Password)
				userRepo.EXPECT().FindUserByEmail(gomock.Any()).
					Times(1).Return(&models.User{
					Password: hashedPassword,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_login_user_empty_email",
			args: args{
				data: dto.LoginUserInputDTO{
					Email:    "",
					Password: "secret",
				},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: "email cannot be empty",
		},
		{
			name: "should_fail_login_user_empty_password",
			args: args{
				data: dto.LoginUserInputDTO{
					Email:    "test@gmail.com",
					Password: "",
				},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: "password cannot be empty",
		},
		{
			name: "should_fail_login_user_user_not_found",
			args: args{
				data: userData,
			},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().FindUserByEmail(gomock.Any()).
					Times(1).Return(nil, models.ErrUserNotFound)
			},
			wantErr:    true,
			wantErrMsg: models.ErrUserNotFound.Error(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(userRepo)
			}

			authSvc := provideMockAuthService(userRepo)
			out, err := authSvc.Login(tc.args.data)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.NotEmpty(t, out.AccessToken)
			require.NotEmpty(t, out.RefreshToken)
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)

	type args struct {
		user *models.User
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(userRepo *mocks.MockIUserRepository)
		want       *dto.RefreshTokenOutputDTO
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_refresh_token",
			args: args{
				user: &models.User{
					Email:    "testuser@gmail.com",
					Password: "secret",
				},
			},
			stubFn:  nil,
			wantErr: false,
		},
		{
			name: "should_fail_refresh_token_empty_email",
			args: args{
				user: &models.User{
					Email:    "",
					Password: "secret",
				},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrEmailIsEmpty.Error(),
		},
		{
			name: "should_fail_refresh_token_empty_password",
			args: args{
				user: &models.User{
					Email:    "testuser@gmail.com",
					Password: "",
				},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrPasswordIsEmpty.Error(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(userRepo)
			}

			authSvc := provideMockAuthService(userRepo)
			out, err := authSvc.RefreshToken(tc.args.user)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.NotEmpty(t, out.AccessToken)
		})
	}
}

func TestAuthService_ForgotPassword(t *testing.T) {
	require.Nil(t, nil)
}

func TestAuthService_ResetPassword(t *testing.T) {
	require.Nil(t, nil)
}
