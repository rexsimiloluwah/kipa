import axios from "../lib/axios";
import { User } from "../common/types/user";
import TokenService from "./token";

class AuthService {
  getAuthUser(): Promise<User> {
    return new Promise((resolve, reject) => {
      const token = TokenService.getLocalAccessToken();
      if (token) {
        const authHeader = { Authorization: `Bearer ${token}` };
        axios
          .get("/auth/user", { headers: authHeader })
          .then((response) => {
            const { data } = response.data;
            resolve(data);
          })
          .catch((error) => {
            reject(error.response.data);
          });
      }
      return;
    });
  }

  login(data: { email: string; password: string }) {
    return new Promise((resolve, reject) => {
      axios
        .post("/auth/login", data)
        .then((response) => {
          const { data } = response;
          TokenService.setLocalAccessToken(data.data.access_token);
          TokenService.setLocalRefreshToken(data.data.refresh_token);
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  register(data: {
    firstname: string;
    lastname: string;
    username?: string;
    email: string;
    password: string;
  }) {
    return new Promise((resolve, reject) => {
      axios
        .post("/auth/register", data)
        .then((response) => {
          const { data } = response;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }
}

export default new AuthService();
