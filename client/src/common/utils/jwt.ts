import jwtDecode from "jwt-decode";

export const isTokenExpired = (token: string) => {
  const user = jwtDecode(token);
  // @ts-ignore
  if (Date.now() >= user.exp * 1000) {
    return true;
  }
  return false;
};
