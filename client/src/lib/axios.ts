import axios from "axios";
import { access } from "fs";
import { isTokenExpired } from "../common/utils/jwt";
import TokenService from "../services/token";

const BASE_URL = "http://localhost:5050/api/v1";
const axiosInstance = axios.create({
  baseURL: BASE_URL,
  timeout: 20000,
  headers: {
    Authorization: `Bearer ${TokenService.getLocalAccessToken()}`,
  },
});

axiosInstance.interceptors.request.use(async (req) => {
  const accessToken = TokenService.getLocalAccessToken();
  const refreshToken = TokenService.getLocalRefreshToken();

  // proceed with the request if there is no access and refresh token
  if (!accessToken && !refreshToken) {
    return req;
  }

  if (isTokenExpired(accessToken as string)) {
    const { data } = await axios.post(`${BASE_URL}/auth/refresh-token`, null, {
      headers: {
        "x-refresh-token": refreshToken as string,
      },
    });

    const accessToken = data.data.access_token;
    TokenService.setLocalAccessToken(accessToken);
    req.headers = { Authorization: `Bearer ${accessToken}` };
    return req;
  }

  return req;
});

export default axiosInstance;
