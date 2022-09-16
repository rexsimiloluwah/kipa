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

// provide the bucket item service
func provideBucketItemService(mockBucketItemRepo *mocks.MockIBucketItemRepository, mockBucketRepo *mocks.MockIBucketRepository) IBucketItemService {
	cfg := &config.Config{}
	return NewBucketItemService(cfg, mockBucketItemRepo, mockBucketRepo)
}

func TestBucketItemService_CreateBucketItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		data      dto.CreateBucketItemInputDTO
		userID    primitive.ObjectID
		bucketUID string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		want       *dto.CreateBucketItemOutputDTO
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_create_bucket_item",
			args: args{
				data: dto.CreateBucketItemInputDTO{
					Key: "test",
				},
				userID:    primitive.NewObjectID(),
				bucketUID: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().FindBucketByUID(gomock.Any()).
					Times(1).Return(&models.Bucket{
					ID: primitive.NewObjectID(),
				}, nil)
				bucketItemRepo.EXPECT().FindBucketItemByKeyName(gomock.Any(), gomock.Any()).
					Times(1).Return(nil, models.ErrBucketItemNotFound)
				bucketItemRepo.EXPECT().CreateBucketItem(gomock.Any()).
					Times(1).Return(primitive.NewObjectID(), nil)
			},
			want: &dto.CreateBucketItemOutputDTO{
				Key:       "test",
				BucketUID: "12345",
			},
			wantErr: false,
		},
		{
			name: "should_successfully_create_bucket_item",
			args: args{
				data: dto.CreateBucketItemInputDTO{
					Key: "test",
				},
				userID:    primitive.NewObjectID(),
				bucketUID: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().FindBucketByUID(gomock.Any()).
					Times(1).Return(&models.Bucket{
					ID: primitive.NewObjectID(),
				}, nil)
				bucketItemRepo.EXPECT().FindBucketItemByKeyName(gomock.Any(), gomock.Any()).
					Times(1).Return(nil, models.ErrBucketItemNotFound)
				bucketItemRepo.EXPECT().CreateBucketItem(gomock.Any()).
					Times(1).Return(primitive.ObjectID{}, errors.New("failed to create bucket item"))
			},
			want:       &dto.CreateBucketItemOutputDTO{},
			wantErr:    true,
			wantErrMsg: "failed to create bucket item",
		},
		{
			name: "should_fail_create_bucket_item_empty_key",
			args: args{
				data: dto.CreateBucketItemInputDTO{
					Key: "",
				},
				userID:    primitive.NewObjectID(),
				bucketUID: "12345",
			},
			stubFn:     nil,
			want:       &dto.CreateBucketItemOutputDTO{},
			wantErr:    true,
			wantErrMsg: ErrKeyIsEmpty.Error(),
		},
		{
			name: "should_fail_create_bucket_item_duplicate_key",
			args: args{
				data: dto.CreateBucketItemInputDTO{
					Key: "test",
				},
				userID:    primitive.NewObjectID(),
				bucketUID: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketRepo.EXPECT().FindBucketByUID(gomock.Any()).
					Times(1).Return(&models.Bucket{
					ID: primitive.NewObjectID(),
				}, nil)
				bucketItemRepo.EXPECT().FindBucketItemByKeyName(gomock.Any(), gomock.Any()).
					Times(1).Return(&models.BucketItem{
					ID: primitive.NewObjectID(),
				}, nil)
			},
			want:       &dto.CreateBucketItemOutputDTO{},
			wantErr:    true,
			wantErrMsg: "key 'test' already exists",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketItemSvc := provideBucketItemService(bucketItemRepo, bucketRepo)
			out, err := bucketItemSvc.CreateBucketItem(tc.args.data, tc.args.userID, tc.args.bucketUID)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, tc.args.data.Key, out.Key)
		})
	}
}

func TestBucketItemService_FindBucketItemByID(t *testing.T) {
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
		want       *models.BucketItem
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_find_bucket_item_by_id",
			args: args{
				id: "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().FindBucketItemByID(gomock.Any()).
					Times(1).Return(&models.BucketItem{
					ID:        primitive.NewObjectID(),
					BucketUID: "12345",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_find_bucket_item_by_id_empty_id",
			args: args{
				id: "",
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrBucketItemIDIsEmpty.Error(),
		},
		{
			name: "should_fail_find_bucket_item_by_id",
			args: args{
				id: "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().FindBucketItemByID(gomock.Any()).
					Times(1).Return(nil, errors.New("failed to find bucket item by id"))
			},
			wantErr:    true,
			wantErrMsg: "failed to find bucket item by id",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketItemSvc := provideBucketItemService(bucketItemRepo, bucketRepo)
			out, err := bucketItemSvc.FindBucketItemByID(tc.args.id)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, out.BucketUID, "12345")
		})
	}
}

func TestBucketItemService_FindBucketItemByKeyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		bucketUID string
		key       string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		want       *models.BucketItem
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_find_bucket_item_by_keyname",
			args: args{
				bucketUID: "12345",
				key:       "key",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().FindBucketItemByKeyName(gomock.Any(), gomock.Any()).
					Times(1).Return(&models.BucketItem{
					ID:        primitive.NewObjectID(),
					BucketUID: "12345",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_find_bucket_item_by_keyname_empty_bucket_uid",
			args: args{
				bucketUID: "",
				key:       "key",
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrBucketUIDIsEmpty.Error(),
		},
		{
			name: "should_fail_find_bucket_item_by_keyname_empty_key",
			args: args{
				bucketUID: "12345",
				key:       "",
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrKeyIsEmpty.Error(),
		},
		{
			name: "should_fail_find_bucket_by_keyname",
			args: args{
				bucketUID: "12345",
				key:       "key",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().FindBucketItemByKeyName(gomock.Any(), gomock.Any()).
					Times(1).Return(nil, errors.New("failed to find bucket item by key name"))
			},
			wantErr:    true,
			wantErrMsg: "failed to find bucket item by key name",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketItemSvc := provideBucketItemService(bucketItemRepo, bucketRepo)
			out, err := bucketItemSvc.FindBucketItemByKeyName(tc.args.bucketUID, tc.args.key)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Equal(t, out.BucketUID, "12345")
		})
	}
}

func TestBucketItemService_ListBucketItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		bucketUID string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		want       []models.BucketItem
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_list_bucket_items",
			args: args{
				bucketUID: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
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
			name: "should_fail_list_bucket_items_empty_bucket_uid",
			args: args{
				bucketUID: "",
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrBucketUIDIsEmpty.Error(),
		},
		{
			name: "should_fail_list_bucket_items",
			args: args{
				bucketUID: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().FindBucketItems(gomock.Any()).
					Times(1).Return([]models.BucketItem{}, errors.New("failed to list bucket items"))
			},
			wantErr:    true,
			wantErrMsg: "failed to list bucket items",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketItemSvc := provideBucketItemService(bucketItemRepo, bucketRepo)
			out, err := bucketItemSvc.ListBucketItems(tc.args.bucketUID)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
			require.Len(t, out, 1)
		})
	}
}

func TestBucketItemService_UpdateBucketItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		data      dto.UpdateBucketItemInputDTO
		bucketUID string
		key       string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_update_bucket_item_by_keyname",
			args: args{
				bucketUID: "12345",
				key:       "key",
				data: dto.UpdateBucketItemInputDTO{
					Key: "updated-key",
				},
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().UpdateBucketItem(gomock.Any(), gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_update_bucket_item_by_keyname_empty_bucket_uid",
			args: args{
				bucketUID: "",
				key:       "key",
				data: dto.UpdateBucketItemInputDTO{
					Key: "updated-key",
				},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrBucketUIDIsEmpty.Error(),
		},
		{
			name: "should_fail_update_bucket_item_by_keyname_empty_keyname",
			args: args{
				bucketUID: "12345",
				key:       "",
				data: dto.UpdateBucketItemInputDTO{
					Key: "updated-key",
				},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrKeyIsEmpty.Error(),
		},
		{
			name: "should_fail_update_bucket_item_by_keyname",
			args: args{
				bucketUID: "12345",
				key:       "key",
				data: dto.UpdateBucketItemInputDTO{
					Key: "updated-key",
				},
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().UpdateBucketItem(gomock.Any(), gomock.Any()).
					Times(1).Return(errors.New("failed to update bucket item"))
			},
			wantErr:    true,
			wantErrMsg: "failed to update bucket item",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketItemSvc := provideBucketItemService(bucketItemRepo, bucketRepo)
			err := bucketItemSvc.UpdateBucketItemByKeyName(tc.args.data, tc.args.bucketUID, tc.args.key)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}

func TestBucketItemService_DeleteBucketItemById(t *testing.T) {
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
			name: "should_successfully_delete_bucket_item_by_id",
			args: args{
				id: "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItemById(gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_delete_bucket_item_by_id_empty_id",
			args: args{
				id: "",
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrBucketItemIDIsEmpty.Error(),
		},
		{
			name: "should_fail_delete_bucket_item_by_id",
			args: args{
				id: "62fa734bfc1cdb7f06a3bf6f",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItemById(gomock.Any()).
					Times(1).Return(errors.New("failed to delete bucket item by id"))
			},
			wantErr:    true,
			wantErrMsg: "failed to delete bucket item by id",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketItemSvc := provideBucketItemService(bucketItemRepo, bucketRepo)
			err := bucketItemSvc.DeleteBucketItemById(tc.args.id)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}

func TestBucketItemService_DeleteBucketItemsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		ids []string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_delete_bucket_items_by_id",
			args: args{
				ids: []string{
					"62fa734bfc1cdb7f06a3bf6f",
				},
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItemsById(gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_delete_bucket_items_by_id",
			args: args{
				ids: []string{
					"62fa734bfc1cdb7f06a3bf6f",
				},
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItemsById(gomock.Any()).
					Times(1).Return(errors.New("failed to delete bucket items"))
			},
			wantErr:    true,
			wantErrMsg: "failed to delete bucket items",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketItemSvc := provideBucketItemService(bucketItemRepo, bucketRepo)
			err := bucketItemSvc.DeleteBucketItemsById(tc.args.ids)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}

func TestBucketItemService_DeleteBucketItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		bucketUID string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_delete_bucket_items",
			args: args{
				bucketUID: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItems(gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "should_fail_delete_bucket_items_empty_bucket_uid",
			args:       args{bucketUID: ""},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrBucketUIDIsEmpty.Error(),
		},
		{
			name: "should_fail_delete_bucket_items",
			args: args{
				bucketUID: "12345",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItems(gomock.Any()).
					Times(1).Return(errors.New("failed to delete bucket items"))
			},
			wantErr:    true,
			wantErrMsg: "failed to delete bucket items",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketItemSvc := provideBucketItemService(bucketItemRepo, bucketRepo)
			err := bucketItemSvc.DeleteBucketItems(tc.args.bucketUID)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}

func TestBucketItemService_DeleteBucketItemByKeyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		bucketUID string
		key       string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_delete_bucket_item_by_keyname",
			args: args{
				bucketUID: "12345",
				key:       "key",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItemByKeyName(gomock.Any(), gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_delete_bucket_item_by_keyname_empty_bucket_uid",
			args: args{
				bucketUID: "",
				key:       "key",
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrBucketUIDIsEmpty.Error(),
		},
		{
			name: "should_fail_delete_bucket_item_by_keyname_empty_key",
			args: args{
				bucketUID: "12345",
				key:       "",
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrKeyIsEmpty.Error(),
		},
		{
			name: "should_fail_delete_bucket_item_by_keyname",
			args: args{
				bucketUID: "12345",
				key:       "key",
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItemByKeyName(gomock.Any(), gomock.Any()).
					Times(1).Return(errors.New("failed to delete bucket item"))
			},
			wantErr:    true,
			wantErrMsg: "failed to delete bucket item",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketItemSvc := provideBucketItemService(bucketItemRepo, bucketRepo)
			err := bucketItemSvc.DeleteBucketItemByKeyName(tc.args.bucketUID, tc.args.key)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}

func TestBucketItemService_DeleteBucketItemsByKeyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketRepo := mocks.NewMockIBucketRepository(ctrl)
	bucketItemRepo := mocks.NewMockIBucketItemRepository(ctrl)

	type args struct {
		bucketUID string
		keys      []string
	}

	tt := []struct {
		name       string
		args       args
		stubFn     func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_successfully_delete_bucket_items_by_keynames",
			args: args{
				bucketUID: "1234",
				keys:      []string{"key1"},
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItemByKeyName(gomock.Any(), gomock.Any()).
					Times(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "should_fail_delete_bucket_items_by_keyname_empty_bucket_uid",
			args: args{
				bucketUID: "",
				keys:      []string{"key1"},
			},
			stubFn:     nil,
			wantErr:    true,
			wantErrMsg: ErrBucketUIDIsEmpty.Error(),
		},
		{
			name: "should_fail_delete_bucket_item_by_keyname",
			args: args{
				bucketUID: "12345",
				keys:      []string{"key1"},
			},
			stubFn: func(bucketRepo *mocks.MockIBucketRepository, bucketItemRepo *mocks.MockIBucketItemRepository) {
				bucketItemRepo.EXPECT().DeleteBucketItemByKeyName(gomock.Any(), gomock.Any()).
					Times(1).Return(errors.New("failed to delete bucket item"))
			},
			wantErr:    true,
			wantErrMsg: "failed to delete bucket item",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn(bucketRepo, bucketItemRepo)
			}

			bucketItemSvc := provideBucketItemService(bucketItemRepo, bucketRepo)
			err := bucketItemSvc.DeleteBucketItemsByKeyName(tc.args.bucketUID, tc.args.keys)
			if tc.wantErr {
				require.NotNil(t, err)
				require.Equal(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.Nil(t, err)
		})
	}
}
