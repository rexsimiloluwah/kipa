class TokenService {
  getLocalAccessToken() {
    return localStorage.getItem("accessToken");
  }

  getLocalRefreshToken() {
    return localStorage.getItem("refreshToken");
  }

  setLocalAccessToken(token: string) {
    localStorage.setItem("accessToken", token);
  }

  setLocalRefreshToken(token: string) {
    localStorage.setItem("refreshToken", token);
  }

  removeLocalAccessToken() {
    localStorage.removeItem("accessToken");
  }

  removeLocalRefreshToken() {
    localStorage.removeItem("refreshToken");
  }
}

export default new TokenService();
