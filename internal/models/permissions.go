package models

type BucketPermission string
type APIKeyPermission string
type BucketPermissionsList []BucketPermission
type APIKeyPermissionsList []APIKeyPermission

var (
	BucketPermissionPublicReadBucket   = BucketPermission("public:read")
	BucketPermissionPublicWriteBucket  = BucketPermission("public:write")
	BucketPermissionPublicDeleteBucket = BucketPermission("public:delete")
	BucketPermissionPublicReadItem     = BucketPermission("public:read:item")
	BucketPermissionPublicWriteItem    = BucketPermission("public:write:item")
	BucketPermissionPublicDeleteItem   = BucketPermission("public:delete:item")
	APIKeyPermissionReadBucket         = APIKeyPermission("read:bucket")
	APIKeyPermissionWriteBucket        = APIKeyPermission("write:bucket")
	APIKeyPermissionDeleteBucket       = APIKeyPermission("delete:bucket")
	APIKeyPermissionReadUser           = APIKeyPermission("read:user")
	APIKeyPermissionWriteUser          = APIKeyPermission("write:user")
	APIKeyPermissionDeleteUser         = APIKeyPermission("delete:user")
	APIKeyPermissionReadItem           = APIKeyPermission("read:item")
	APIKeyPermissionWriteItem          = APIKeyPermission("write:item")
	APIKeyPermissionDeleteItem         = APIKeyPermission("delete:item")
	BUCKET_PERMISSIONS                 = BucketPermissionsList{
		BucketPermissionPublicWriteItem,
		BucketPermissionPublicWriteBucket,
		BucketPermissionPublicDeleteBucket,
		BucketPermissionPublicDeleteItem,
		BucketPermissionPublicReadItem,
		BucketPermissionPublicReadBucket,
	}
	APIKEY_PERMISSIONS = APIKeyPermissionsList{
		APIKeyPermissionReadBucket,
		APIKeyPermissionWriteBucket,
		APIKeyPermissionDeleteBucket,
		APIKeyPermissionReadUser,
		APIKeyPermissionWriteUser,
		APIKeyPermissionDeleteUser,
		APIKeyPermissionReadItem,
		APIKeyPermissionWriteItem,
		APIKeyPermissionDeleteItem,
	}
)

func (b BucketPermission) String() string {
	return string(b)
}

func (a APIKeyPermission) String() string {
	return string(a)
}

// check if the bucket permissions list contains a specific permission
// Returns true if the permission is present in the bucket permissions list
func (b BucketPermissionsList) Contains(permission BucketPermission) bool {
	for _, p := range b {
		if p == permission {
			return true
		}
	}
	return false
}

// check if the apikey permissions list contains a specific permission
// Returns true if the permission is present in the apikey permissions list
func (a APIKeyPermissionsList) Contains(permission APIKeyPermission) bool {
	for _, p := range a {
		if p == permission {
			return true
		}
	}
	return false
}
