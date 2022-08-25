package services

import (
	"errors"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/mocks"
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

}

func TestBucketService_FindBucketByUID(t *testing.T) {

}

func TestBucketService_ListUserBuckets(t *testing.T) {

}

func TestBucketService_UpdateBucket(t *testing.T) {

}

func TestBucketService_DeleteBucket(t *testing.T) {

}
