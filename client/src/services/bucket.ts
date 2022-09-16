import { BucketDetails, CreateBucketData } from "../common/types/bucket";
import axios from "../lib/axios";
import TokenService from "./token";

class BucketService {
  /**
   * Create a new Bucket
   * @param data // accepts the bucket data
   * @returns
   */
  async createBucket(data: CreateBucketData) {
    return new Promise((resolve, reject) => {
      axios
        .post("/bucket?full=true", data)
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Find a bucket by UID
   * @param uid // accepts the bucket UID
   * @returns
   */
  findBucket(uid: string): Promise<BucketDetails> {
    return new Promise((resolve, reject) => {
      axios
        .get(`/bucket/${uid}`)
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Fetch user's buckets
   * @returns
   */
  getUserBuckets(): Promise<Array<BucketDetails>> {
    return new Promise((resolve, reject) => {
      const token = TokenService.getLocalAccessToken();
      const authHeader = { Authorization: `Bearer ${token}` };
      axios
        .get("/buckets", { headers: authHeader })
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Update a user's bucket data
   * @param uid // accepts the UID of the bucket
   * @param data // accepts the data for updating the bucket
   * @returns
   */
  updateBucket(uid: string, data: CreateBucketData) {
    return new Promise((resolve, reject) => {
      axios
        .put(`/bucket/${uid}`, data)
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Delete a user's bucket
   * @param uid // accepts the bucket UID
   * @returns
   */
  deleteBucket(uid: string) {
    return new Promise((resolve, reject) => {
      axios
        .delete(`/bucket/${uid}`)
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }
}

export default new BucketService();
