package services

import (
	"errors"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/mocks"
	"keeper/internal/models"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// provide user service
func provideUserService(mockUserRepo *mocks.MockIUserRepository) IUserService {
	cfg := &config.Config{
		Env: "test",
	}
	return NewUserService(cfg, mockUserRepo)
}

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)

	input := dto.CreateUserInputDTO{
		Firstname: "bola",
		Lastname:  "tinubu",
		Email:     "bolatinubu@gmail.com",
		Password:  "secret",
		Username:  "testuser",
	}

	type args struct {
		data dto.CreateUserInputDTO
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(userRepo *mocks.MockIUserRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_register_new_user",
			args: args{
				data: input,
			},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().CreateUser(gomock.Any()).
					Times(1).Return(nil)
				userRepo.EXPECT().FindUserByEmail(gomock.Any()).
					Times(1).Return(nil, models.ErrUserNotFound)
			},
		},
		{
			name: "should_fail_user_already_exists",
			args: args{
				data: input,
			},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().FindUserByEmail(gomock.Any()).
					Times(1).Return(&models.User{
					ID: primitive.NewObjectID(),
				}, nil)
			},
			wantErr:    true,
			wantErrMsg: models.ErrUserAlreadyExists.Error(),
		},
		{
			name: "should_fail_register_new_user",
			args: args{
				data: input,
			},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().CreateUser(gomock.Any()).
					Times(1).Return(errors.New("failed to register new user"))
				userRepo.EXPECT().FindUserByEmail(gomock.Any()).
					Times(1).Return(nil, models.ErrUserNotFound)
			},
			wantErr:    true,
			wantErrMsg: "failed to register new user",
		},
		{
			name: "should_fail_register_new_user_empty_email",
			args: args{
				data: dto.CreateUserInputDTO{
					Password:  "secret",
					Firstname: "bola",
					Lastname:  "tinubu",
				},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: "email must not be empty",
		},
		{
			name: "should_fail_register_new_user_empty_password",
			args: args{
				data: dto.CreateUserInputDTO{
					Email:     "bolatinubu@gmail.com",
					Firstname: "bola",
					Lastname:  "tinubu",
				},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: "password must not be empty",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(userRepo)
			}

			userSvc := provideUserService(userRepo)
			err := userSvc.Register(tc.args.data)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}

func TestUserService_FindUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)

	type args struct {
		id string
	}

	testUser := &models.User{
		ID:        primitive.NewObjectID(),
		Firstname: "bola",
		Lastname:  "tinubu",
		Email:     "bolatinubu@gmail.com",
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(userRepo *mocks.MockIUserRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_find_user_by_id",
			args: args{
				id: "12345",
			},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().FindUserById(gomock.Any()).
					Times(1).Return(testUser, nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_find_user_by_id",
			args: args{
				id: "12345",
			},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().FindUserById(gomock.Any()).
					Times(1).Return(nil, errors.New("failed to find user by id"))
			},
			wantErr:    true,
			wantErrMsg: "failed to find user by id",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(userRepo)
			}

			userSvc := provideUserService(userRepo)
			user, err := userSvc.FindUserByID(tc.args.id)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, user.ID, testUser.ID)
			require.Equal(t, user.Email, testUser.Email)
			require.Equal(t, user.CreatedAt, testUser.CreatedAt)
			require.Equal(t, user.UpdatedAt, testUser.UpdatedAt)
		})
	}
}

func TestUserService_FindAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)

	testUsers := []models.User{
		{
			ID:        primitive.NewObjectID(),
			Firstname: "bola",
			Lastname:  "tinubu",
			Email:     "bolatinubu@gmail.com",
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	tt := []struct {
		name       string
		stubFn     func(userRepo *mocks.MockIUserRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_find_all_users",
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().FindAllUsers().
					Times(1).Return(testUsers, nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_find_all_users",
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().FindAllUsers().
					Times(1).Return([]models.User{}, errors.New("failed to find all users"))
			},
			wantErr:    true,
			wantErrMsg: "failed to find all users",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(userRepo)
			}

			userSvc := provideUserService(userRepo)
			users, err := userSvc.FindAllUsers()
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, len(users), len(testUsers))
			require.Equal(t, users[0].ID, testUsers[0].ID)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)

	type args struct {
		id   string
		data dto.UpdateUserInputDTO
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(userRepo *mocks.MockIUserRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_update_user",
			args: args{
				id: "62fa734bfc1cdb7f06a3bf6f",
				data: dto.UpdateUserInputDTO{
					Firstname: "peter",
					Lastname:  "obi",
				},
			},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().UpdateUser(gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_update_user",
			args: args{
				id: "62fa734bfc1cdb7f06a3bf6f",
				data: dto.UpdateUserInputDTO{
					Firstname: "peter",
					Lastname:  "obi",
				},
			},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().UpdateUser(gomock.Any()).
					Times(1).Return(errors.New("failed to update user"))
			},
			wantErr:    true,
			wantErrMsg: "failed to update user",
		},
		{
			name: "should_fail_update_user_invalid_id",
			args: args{
				id: "abcde",
				data: dto.UpdateUserInputDTO{
					Firstname: "peter",
					Lastname:  "obi",
				},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: models.ErrInvalidObjectID.Error(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(userRepo)
			}

			userSvc := provideUserService(userRepo)
			err := userSvc.UpdateUser(tc.args.id, tc.args.data)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}

func TestUserService_UpdateUserPassword(t *testing.T) {
	require.Nil(t, nil)
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockIUserRepository(ctrl)

	type args struct {
		id string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(userRepo *mocks.MockIUserRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_delete_user",
			args: args{id: "62fa734bfc1cdb7f06a3bf6f"},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().DeleteUser(gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_delete_user",
			args: args{id: "abcde"},
			stubFn: func(userRepo *mocks.MockIUserRepository) {
				userRepo.EXPECT().DeleteUser(gomock.Any()).
					Times(1).Return(errors.New("failed to delete user"))
			},
			wantErr:    true,
			wantErrMsg: "failed to delete user",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(userRepo)
			}

			userSvc := provideUserService(userRepo)
			err := userSvc.DeleteUser(tc.args.id)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}
