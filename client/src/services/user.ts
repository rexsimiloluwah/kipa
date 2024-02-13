import axios from "../lib/axios";
import { UpdateUserData } from "../common/types/user";

class UserService {
  /**
   * Update user
   * @param data // Updated user payload
   * @returns
   */
  updateUser(data: UpdateUserData) {
    return new Promise((resolve, reject) => {
      axios
        .put(`/user`, data)
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
   * Update user password
   * @param data // Updated password payload
   * @returns
   */
  updateUserPassword(data: { password: string }) {
    return new Promise((resolve, reject) => {
      axios
        .put(`/user/password`, data)
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
   * Delete user
   * @returns
   */
  deleteUser() {
    return new Promise((resolve, reject) => {
      axios
        .delete(`/user`)
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

export default new UserService();
