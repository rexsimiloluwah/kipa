package services

import (
	"errors"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/mocks"
	"keeper/internal/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// provide the bucket service
func provideBucketService(mockBucketRepo *mocks.MockIBucketRepository, mockBucketItemRepo *mocks.MockIBucketItemRepository) IBucketService {
	cfg := &config.Config{}
	return NewBucketService(cfg, mockBucketRepo, mockBucketItemRepo)
}

func TestBucketService_CreateBucket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		data   dto.CreateBucketInputDTO
		userID primitive.ObjectID
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository)
		want       *dto.CreateBucketOutputDTO
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_create_new_bucket",
			args: args{
				data: dto.CreateBucketInputDTO{
					Name: "testbucket",
				},
				userID: primitive.NewObjectID(),
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository) {
				bucketRepo.EXPECT().CreateBucket(gomock.Any()).
					Times(1).Return(primitive.NewObjectID(), nil)
			},
			want: &dto.CreateBucketOutputDTO{
				Name: "12345",
			},
			wantErr: false,
		},
		{
			name: "should_fail_create_new_bucket_empty_name",
			args: args{
				data: dto.CreateBucketInputDTO{
					Name: "",
				},
				userID: primitive.NewObjectID(),
			},
			stubFn:     nil,
			want:       nil,
			wantErr:    true,
			wantErrMsg: ErrBucketNameIsEmpty.Error(),
		},
		{
			name: "should_fail_create_bucket",
			args: args{
				data: dto.CreateBucketInputDTO{
					Name: "testbucket",
				},
				userID: primitive.NewObjectID(),
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository) {
				bucketRepo.EXPECT().CreateBucket(gomock.Any()).
					Times(1).Return(primitive.ObjectID{}, errors.New("error"))
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "error saving bucket to database: error",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo)
			}

			bucketSvc := provideBucketService(bucketRepo, bucketItemRepo)
			out, err := bucketSvc.CreateBucket(tc.args.data, tc.args.userID)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, tc.args.data.Name, out.Name)
		})
	}
}

func TestBucketService_FindBucketByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		id string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_find_bucket_by_id",
			args: args{
				id: "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().FindBucketByID(gomock.Any()).
					Times(1).Return(&models.Bucket{
					ID:  primitive.NewObjectID(),
					UID: "12345",
				}, nil)
				bucketItemRepo.EXPECT().FindBucketItems(gomock.Any()).
					Times(1).Return([]models.BucketItem{
					{
						ID: primitive.NewObjectID(),
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_find_bucket_by_id_empty_id",
			args: args{
				id: "",
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrBucketIDIsEmpty.Error(),
		},
		{
			name: "should_fail_find_bucket_by_id_bucket_not_found_error",
			args: args{
				id: "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().FindBucketByID(gomock.Any()).
					Times(1).Return(nil, errors.New("failed to find bucket"))
			},
			wantErr:    true,
			wantErrMsg: "failed to find bucket",
		},
		{
			name: "should_fail_find_bucket_by_id_find_bucket_items_error",
			args: args{
				id: "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().FindBucketByID(gomock.Any()).
					Times(1).Return(&models.Bucket{
					ID:  primitive.NewObjectID(),
					UID: "12345",
				}, nil)
				bucketItemRepo.EXPECT().FindBucketItems(gomock.Any()).
					Times(1).Return(nil, errors.New("failed to find bucket items"))
			},
			wantErr:    true,
			wantErrMsg: "failed to find bucket items",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketSvc := provideBucketService(bucketRepo, bucketItemRepo)
			out, err := bucketSvc.FindBucketByID(tc.args.id)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, out.UID, "12345")
		})
	}
}

func TestBucketService_FindBucketByUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		uid string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		want       *dto.BucketDetailsOutput
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_find_bucket_by_uid",
			args: args{
				uid: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().FindBucketByUID(gomock.Any()).
					Times(1).Return(&models.Bucket{
					ID:  primitive.NewObjectID(),
					UID: "12345",
				}, nil)
				bucketItemRepo.EXPECT().FindBucketItems(gomock.Any()).
					Times(1).Return([]models.BucketItem{
					{
						ID: primitive.NewObjectID(),
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_find_bucket_by_uid_empty_uid",
			args: args{
				uid: "",
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrBucketUIDIsEmpty.Error(),
		},
		{
			name: "should_fail_find_bucket_by_uid_bucket_not_found_error",
			args: args{
				uid: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().FindBucketByUID(gomock.Any()).
					Times(1).Return(nil, errors.New("failed to find bucket by uid"))
			},
			wantErr:    true,
			wantErrMsg: "failed to find bucket by uid",
		},
		{
			name: "should_fail_find_bucket_by_uid_find_bucket_items_error",
			args: args{
				uid: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().FindBucketByUID(gomock.Any()).
					Times(1).Return(&models.Bucket{
					ID:  primitive.NewObjectID(),
					UID: "12345",
				}, nil)
				bucketItemRepo.EXPECT().FindBucketItems(gomock.Any()).
					Times(1).Return([]models.BucketItem{}, errors.New("failed to find bucket items"))
			},
			wantErr:    true,
			wantErrMsg: "failed to find bucket items",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketSvc := provideBucketService(bucketRepo, bucketItemRepo)
			out, err := bucketSvc.FindBucketByUID(tc.args.uid)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, out.UID, "12345")
		})
	}
}

func TestBucketService_ListUserBuckets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		userID string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_list_user_buckets",
			args: args{
				userID: "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().FindBucketsByUserID(gomock.Any()).
					Times(1).Return([]models.Bucket{
					{
						ID:  primitive.NewObjectID(),
						UID: "12345",
					},
				}, nil)
				bucketItemRepo.EXPECT().FindBucketItems(gomock.Any()).
					Times(1).Return([]models.BucketItem{
					{
						ID: primitive.NewObjectID(),
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_list_user_buckets_find_user_buckets_error",
			args: args{
				userID: "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().FindBucketsByUserID(gomock.Any()).
					Times(1).Return([]models.Bucket{}, errors.New("failed to find buckets"))
			},
			wantErr:    true,
			wantErrMsg: "failed to find buckets",
		},
		{
			name: "should_fail_list_user_buckets_find_bucket_items_error",
			args: args{
				userID: "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().FindBucketsByUserID(gomock.Any()).
					Times(1).Return([]models.Bucket{
					{
						ID:  primitive.NewObjectID(),
						UID: "12345",
					},
				}, nil)
				bucketItemRepo.EXPECT().FindBucketItems(gomock.Any()).
					Times(1).Return([]models.BucketItem{}, errors.New("failed to find bucket items"))
			},
			wantErr:    true,
			wantErrMsg: "failed to find bucket items",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketSvc := provideBucketService(bucketRepo, bucketItemRepo)
			out, err := bucketSvc.ListUserBuckets(tc.args.userID)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, out[0].UID, "12345")
		})
	}
}

func TestBucketService_UpdateBucket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		uid  string
		data dto.UpdateBucketInputDTO
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_update_bucket",
			args: args{
				uid: "12345",
				data: dto.UpdateBucketInputDTO{
					Name: "updated-bucket-name",
				},
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().UpdateBucket(gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_update_bucket_empty_uid",
			args: args{
				uid: "",
				data: dto.UpdateBucketInputDTO{
					Name: "updated-bucket-name",
				},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrBucketUIDIsEmpty.Error(),
		},
		{
			name: "should_fail_update_bucket",
			args: args{
				uid: "12345",
				data: dto.UpdateBucketInputDTO{
					Name: "updated-bucket-name",
				},
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().UpdateBucket(gomock.Any()).
					Times(1).Return(errors.New("failed to update bucket"))
			},
			wantErr:    true,
			wantErrMsg: "failed to update bucket",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketSvc := provideBucketService(bucketRepo, bucketItemRepo)
			err := bucketSvc.UpdateBucket(tc.args.uid, tc.args.data)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}

func TestBucketService_DeleteBucket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		uid string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_delete_bucket",
			args: args{
				uid: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItems(gomock.Any()).
					Times(1).Return(nil)
				bucketRepo.EXPECT().DeleteBucketByUID(gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_delete_bucket_empty_uid",
			args: args{
				uid: "",
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrBucketUIDIsEmpty.Error(),
		},
		{
			name: "should_fail_delete_bucket_could_not_delete_bucket_items",
			args: args{
				uid: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItems(gomock.Any()).
					Times(1).Return(errors.New("failed to delete bucket items"))
			},
			wantErr:    true,
			wantErrMsg: "failed to delete bucket items",
		},
		{
			name: "should_fail_delete_bucket",
			args: args{
				uid: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItems(gomock.Any()).
					Times(1).Return(nil)
				bucketRepo.EXPECT().DeleteBucketByUID(gomock.Any()).
					Times(1).Return(errors.New("failed to delete bucket"))
			},
			wantErr:    true,
			wantErrMsg: "failed to delete bucket",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketSvc := provideBucketService(bucketRepo, bucketItemRepo)
			err := bucketSvc.DeleteBucket(tc.args.uid)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}
