import Cookies from "js-cookie";

class TokenService {
  /**
   * Returns the access token in the cookie
   * @returns
   */
  getAccessTokenCookie() {
    return Cookies.get("accessToken");
  }

  /**
   * Returns the refresh token in the cookie
   * @returns
   */
  getRefreshTokenCookie() {
    return Cookies.get("refreshToken");
  }

  /**
   * Saves the access token in the cookie
   * @param token
   */
  setAccessTokenCookie(token: string) {
    Cookies.set("accessToken", token, {
      sameSite: "strict",
      secure: true,
      expires: new Date(new Date().getTime() + 15 * 60 * 1000), // expires in 15 minutes
    });
  }

  /**
   * Saves the refresh token in the cookie
   * @param token
   */
  setRefreshTokenCookie(token: string) {
    Cookies.set("refreshToken", token, {
      sameSite: "strict",
      secure: true,
      expires: 365 * 2, // 2 years
    });
  }

  /**
   * Removes the access token from the cookie
   */
  removeAccessTokenCookie() {
    Cookies.remove("accessToken", {
      sameSite: "strict",
      secure: true,
      expires: new Date(new Date().getTime() + 15 * 60 * 1000), // expires in 15 minutes
    });
  }

  /**
   * Removes the refresh token from the cookie
   */
  removeRefreshTokenCookie() {
    Cookies.remove("refreshToken", {
      sameSite: "strict",
      secure: true,
      expires: 365,
    });
  }

  /**
   * Returns the access token in the local storage
   * @returns
   */
  getLocalAccessToken() {
    return localStorage.getItem("accessToken");
  }

  /**
   * Returns the refresh token in the local storage
   * @returns
   */
  getLocalRefreshToken() {
    return localStorage.getItem("refreshToken");
  }

  /**
   * Sets the access token in the local storage
   * @param token
   */
  setLocalAccessToken(token: string) {
    localStorage.setItem("accessToken", token);
  }

  /**
   * Sets the refresh token in the local storage
   * @param token
   */
  setLocalRefreshToken(token: string) {
    localStorage.setItem("refreshToken", token);
  }

  /**
   * Removes the access token from the local storage
   */
  removeLocalAccessToken() {
    localStorage.removeItem("accessToken");
  }

  /**
   * Removes the refresh token from the local storage
   */
  removeLocalRefreshToken() {
    localStorage.removeItem("refreshToken");
  }
}

export default new TokenService();
