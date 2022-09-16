import { BucketItem } from "./bucket_item";

export type Bucket = {
  id: string;
  uid: string;
  user_id: string;
  name: string;
  description: string;
  permissions: Array<string>;
  created_at: Date;
  updated_at: Date;
};

export type BucketDetails = Bucket & {
  bucket_items: Array<BucketItem>;
};

export type CreateBucketData = {
  name: string;
  description?: string;
  permissions?: Array<string>;
};
