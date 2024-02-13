import axios from "axios";
import { isTokenExpired } from "../common/utils/jwt";
import TokenService from "../services/token";

const BASE_URL = "http://localhost:5050/api/v1";
const axiosInstance = axios.create({
  baseURL: BASE_URL,
  timeout: 20000,
  headers: {
    Authorization: `Bearer ${TokenService.getAccessTokenCookie()}`,
  },
});

axiosInstance.interceptors.request.use(async (req) => {
  const accessToken = TokenService.getAccessTokenCookie();
  const refreshToken = TokenService.getRefreshTokenCookie();

  // proceed with the request if there is no access and refresh token
  if (!accessToken && !refreshToken) {
    return req;
  }

  // attach the accessToken to the request if it does not have it
  const reqToken = (req.headers as { Authorization: string })[
    "Authorization"
  ].split(" ")[1];

  if (accessToken && (!reqToken || reqToken === "undefined")) {
    req.headers = { Authorization: `Bearer ${accessToken}` };
    return req;
  }

  if (!accessToken || isTokenExpired(accessToken as string)) {
    const { data } = await axios.post(`${BASE_URL}/auth/refresh-token`, null, {
      headers: {
        "x-refresh-token": refreshToken as string,
      },
    });

    const accessToken = data.data.access_token;
    TokenService.setAccessTokenCookie(accessToken);
    req.headers = { Authorization: `Bearer ${accessToken}` };
    return req;
  }

  return req;
});

export default axiosInstance;
