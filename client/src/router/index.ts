import { createRouter, createWebHistory } from "vue-router";
import Home from "../views/Home.vue";
import Login from "../views/Login.vue";
import Signup from "../views/Signup.vue";
import Dashboard from "../views/Dashboard.vue";
import ForgotPassword from "../views/ForgotPassword.vue";
import ResetPassword from "../views/ResetPassword.vue";
import VerifyEmail from "../views/VerifyEmail.vue";
import {
  Buckets,
  UserSettings,
  APIKeys,
  BucketDetails,
  APIKeyDetails,
} from "../components/Dashboard";
import TokenService from "../services/token";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
    meta: {
      guest: true,
    },
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
  },
  {
    path: "/signup",
    name: "Signup",
    component: Signup,
    meta: {
      guest: true,
    },
  },
  {
    path: "/verify-email",
    name: "VerifyEmail",
    component: VerifyEmail,
  },
  {
    path: "/forgot-password",
    name: "ForgotPassword",
    component: ForgotPassword,
  },
  {
    path: "/reset-password",
    name: "ResetPassword",
    component: ResetPassword,
  },
  {
    path: "/dashboard",
    component: Dashboard,
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: "",
        component: Buckets,
      },
      {
        path: "apikeys",
        component: APIKeys,
      },
      {
        path: "user-settings",
        component: UserSettings,
      },
      {
        path: "bucket/:uid",
        component: BucketDetails,
      },
      {
        path: "apikey/:id",
        component: APIKeyDetails,
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Using Navigation guards to protect authenticated routes by redirecting to specific routes
router.beforeEach((to, _, next) => {
  // Check whether a certain route requires authentication
  const accessToken = TokenService.getRefreshTokenCookie();
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
  const isGuestRoute = to.matched.some((record) => record.meta.guest);

  // If it is a guest route and an access token exists
  if (isGuestRoute && accessToken) {
    next("/dashboard");
  }
  // Go to login page if un-authenticated
  if (requiresAuth && !accessToken) {
    next("/login");
  } else next();
});

export default router;
