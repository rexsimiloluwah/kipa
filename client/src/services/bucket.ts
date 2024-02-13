import { BucketDetails, CreateBucketData } from "../common/types/bucket";
import axios from "../lib/axios";
import TokenService from "./token";

class BucketService {
  /**
   * Create a new Bucket
   * @param data // Bucket payload
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
   * @param uid // Bucket UID
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
      const token = TokenService.getAccessTokenCookie();
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
   * @param uid // Bucket UID
   * @param data // Updated bucket payload
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
   * @param uid // Bucket UID
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
