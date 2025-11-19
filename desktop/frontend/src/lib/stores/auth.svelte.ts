let isLoggedIn = $state(false);

export const authState = {
  get isLoggedIn() {
    return isLoggedIn;
  },
  login() {
    isLoggedIn = true;
  },
  logout() {
    isLoggedIn = false;
  },
};
