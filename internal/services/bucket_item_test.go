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

}

func TestBucketItemService_FindBucketItemByKeyName(t *testing.T) {

}

func TestBucketItemService_ListBucketItems(t *testing.T) {

}

func TestBucketItemService_UpdateBucketItem(t *testing.T) {

}

func TestBucketItemService_DeleteBucketItemById(t *testing.T) {

}

func TestBucketItemService_DeleteBucketItemsById(t *testing.T) {

}

func TestBucketItemService_DeleteBucketItems(t *testing.T) {

}

func TestBucketItemService_DeleteBucketItemByKeyName(t *testing.T) {

}
