export type CreateBucketItemData = {
  id: number;
  key: string;
  data: any;
  ttl: number;
};

export type BucketItem = {
  id: string;
  user_id: string;
  bucket_id: string;
  bucket_uid: string;
  key: string;
  data: any;
  ttl: number;
  created_at: Date;
  updated_at: Date;
};
