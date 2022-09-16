import axios from "../lib/axios";
import { UpdateUserData } from "../common/types/user";

class UserService {
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
