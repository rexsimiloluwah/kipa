import { CreateUserData, LoginUserData } from "../common/types/user";
import axios from "../lib/axios";
import TokenService from "./token";

class AuthService {
  /**
   * Fetch the authenticated user
   * @returns
   */
  getAuthUser() {
    return new Promise((resolve, reject) => {
      const token = TokenService.getAccessTokenCookie();
      const refreshToken = TokenService.getRefreshTokenCookie();
      if (token || refreshToken) {
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

  /**
   * Login a user
   * @param data
   * @returns
   */
  login(data: LoginUserData) {
    return new Promise((resolve, reject) => {
      axios
        .post("/auth/login", data)
        .then((response) => {
          const { data } = response;
          TokenService.setAccessTokenCookie(data.data.access_token);
          TokenService.setRefreshTokenCookie(data.data.refresh_token);
          resolve(data);
        })
        .catch((error) => {
          console.log(error);
          reject(error.response.data);
        });
    });
  }

  /**
   * Logout a user
   */
  logout() {
    Promise.all([
      TokenService.removeAccessTokenCookie(),
      TokenService.removeLocalAccessToken(),
      TokenService.removeRefreshTokenCookie(),
      TokenService.removeLocalRefreshToken(),
    ]);
  }

  /**
   * Register a new user
   * @param data
   * @returns
   */
  register(data: CreateUserData) {
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
