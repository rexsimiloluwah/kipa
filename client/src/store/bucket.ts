import { defineStore } from "pinia";
import { BucketDetails, CreateBucketData } from "../common/types/bucket";
import BucketService from "../services/bucket";
import { useToast } from "vue-toastification";
import router from "../router";

const toast = useToast();

export const useBucketStore = defineStore("bucket", {
  state: () => ({
    buckets: [] as Array<BucketDetails>,
    isLoadingBuckets: false as boolean, // loading state for the buckets
  }),
  getters: {
    bucketsCount: (state) => state.buckets.length,
    bucketItemsCount: (state) =>
      state.buckets.reduce(
        (prevValue, currValue) => prevValue + currValue.bucket_items.length,
        0
      ),
    bucket: (state) => {
      return (uid: string) =>
        state.buckets.filter((bucket) => bucket.uid === uid)[0];
    },
  },
  actions: {
    async fetchBuckets() {
      try {
        this.isLoadingBuckets = true;
        const buckets = await BucketService.getUserBuckets();
        // console.log("buckets: ", buckets);
        this.buckets = buckets;
        this.isLoadingBuckets = false;
      } catch (error: any) {
        toast.error(error.error || error.message);
        this.buckets = [];
        this.isLoadingBuckets = false;
      }
    },
    async createBucket(data: CreateBucketData) {
      try {
        await BucketService.createBucket(data);
        toast.success("Bucket Created Successfully! 🎉");
        this.fetchBuckets();
      } catch (error: any) {
        toast.error(error.error || error.message);
      }
    },
    async findBucket() {},
    async updateBucket(uid: string, data: CreateBucketData) {
      try {
        await BucketService.updateBucket(uid, data);
        toast.success("Bucket Updated Successfully! 🎉");
        this.fetchBuckets();
      } catch (error: any) {
        toast.error(error.error || error.message);
      }
    },
    async deleteBucket(uid: string) {
      try {
        await BucketService.deleteBucket(uid);
        toast.success("Bucket Deleted Successfully! 🎉");
        router.push("/dashboard");
        this.fetchBuckets();
      } catch (error: any) {
        toast.error(error.error);
      }
    },
  },
});
