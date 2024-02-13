import axios from "../lib/axios";
import {
  BucketItem,
  CreateBucketItemData,
  UpdateBucketItemData,
} from "../common/types/bucket_item";
import { AxiosResponse } from "axios";
import {
  PaginatedSuccessResponse,
  SuccessResponse,
} from "../common/types/response";

class BucketItemService {
  /**
   * Create a new bucket item
   * @param bucketUID // Bucket UID
   * @param data // Bucket item payload
   * @returns
   */
  createBucketItem(
    bucketUID: string,
    data: CreateBucketItemData
  ): Promise<AxiosResponse<SuccessResponse<BucketItem>>> {
    return new Promise((resolve, reject) => {
      axios
        .post(`/item/${bucketUID}`, data)
        .then((response) => {
          resolve(response);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * List all the bucket items
   * @param bucketUID // Bucket UID
   * @param queryParams // Query params for pagination, filtering, and sorting
   * @returns
   */
  listBucketItems(
    bucketUID: string,
    queryParams?: { [key: string]: any }
  ): Promise<AxiosResponse<PaginatedSuccessResponse<BucketItem[]>>> {
    return new Promise((resolve, reject) => {
      axios
        .get(`/items`, {
          params: {
            bucket_uid: bucketUID,
            ...queryParams,
          },
        })
        .then((response) => {
          resolve(response);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Find bucket item by key
   * @param bucketUID // Bucket UID
   * @param key // Bucket item key
   * @returns
   */
  findBucketItemByKey(
    bucketUID: string,
    key: string
  ): Promise<AxiosResponse<SuccessResponse<BucketItem>>> {
    return new Promise((resolve, reject) => {
      axios
        .get(`/item/${bucketUID}/${key}`)
        .then((response) => {
          resolve(response);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Update bucket item by key
   * @param bucketUID // Bucket UID
   * @param key // Bucket item key
   * @param data // Updated bucket item payload
   * @returns
   */
  updateBucketItemByKey(
    bucketUID: string,
    key: string,
    data: UpdateBucketItemData
  ): Promise<AxiosResponse<SuccessResponse<BucketItem>>> {
    return new Promise((resolve, reject) => {
      axios
        .put(`/item/${bucketUID}/${key}`, data)
        .then((response) => {
          resolve(response);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Delete bucket item by key
   * @param bucketUID // Bucket UID
   * @param key // Bucket item key
   * @returns
   */
  deleteBucketItemByKey(
    bucketUID: string,
    key: string
  ): Promise<AxiosResponse<SuccessResponse<BucketItem>>> {
    return new Promise((resolve, reject) => {
      axios
        .delete(`/item/${bucketUID}/${key}`)
        .then((response) => {
          resolve(response);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }
}

export default new BucketItemService();
