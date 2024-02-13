import { defineStore } from "pinia";
import {
  BucketItem,
  CreateBucketItemData,
  UpdateBucketItemData,
} from "../common/types/bucket_item";
import BucketItemService from "../services/bucket_item";
import { useToast } from "vue-toastification";
import { PaginatedSuccessResponse } from "../common/types/response";

const toast = useToast();

export const useBucketItemStore = defineStore("bucket_item", {
  state: () => ({
    activeBucketItem: null as BucketItem | null,
    activeBucketItemsResponse: null as PaginatedSuccessResponse<
      BucketItem[]
    > | null,
  }),
  getters: {},
  actions: {
    async fetchBucketItems(
      bucketUID: string,
      queryParams?: { [key: string]: any }
    ) {
      try {
        const result = await BucketItemService.listBucketItems(
          bucketUID,
          queryParams
        );
        this.activeBucketItemsResponse = result.data;
      } catch (error: any) {
        toast.error(
          `failed to fetch bucket items: ${error.error || error.message}`
        );
      }
    },
    async createBucketItem(bucketUID: string, data: CreateBucketItemData) {
      try {
        await BucketItemService.createBucketItem(bucketUID, data);
        toast.success("Bucket item created successfully!");
        // update the bucket items
        await this.fetchBucketItems(bucketUID);
      } catch (error: any) {
        toast.error(
          `failed to create bucket item: ${error.error || error.message}`
        );
      }
    },
    async setActiveBucketItem(bucketItem: BucketItem) {
      this.activeBucketItem = bucketItem;
    },
    async updateBucketItem(
      bucketUID: string,
      key: string,
      data: UpdateBucketItemData
    ) {
      try {
        await BucketItemService.updateBucketItemByKey(bucketUID, key, data);
        toast.success("Bucket item updated successfully!");
        // update the bucket items state
        await this.fetchBucketItems(bucketUID);
      } catch (error: any) {
        toast.error(
          `failed to update bucket item: ${error.error || error.message}`
        );
      }
    },
    async deleteBucketItem(bucketUID: string, key: string) {
      try {
        await BucketItemService.deleteBucketItemByKey(bucketUID, key);
        toast.success("Bucket item deleted successfully!");
        // update the bucket items state
        await this.fetchBucketItems(bucketUID);
      } catch (error: any) {
        toast.error(
          `failed to delete bucket item: ${error.error || error.message}`
        );
      }
    },
    async clearActiveBucketItem() {
      this.activeBucketItem = null;
    },
  },
});
