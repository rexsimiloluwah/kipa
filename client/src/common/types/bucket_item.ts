export type CreateBucketItemData = {
  key: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data: any;
  ttl: number;
};

export type UpdateBucketItemData = {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data?: any;
  ttl?: number;
};

export type BucketItem = {
  id: string;
  user_id: string;
  bucket_id: string;
  bucket_uid: string;
  key: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data: any;
  type: string;
  ttl: number;
  created_at: Date;
  updated_at: Date;
};
