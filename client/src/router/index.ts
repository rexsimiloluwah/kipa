import { createRouter, createWebHistory } from "vue-router";
import Home from "../views/Home.vue";
import Login from "../views/Login.vue";
import Signup from "../views/Signup.vue";
import Dashboard from "../views/Dashboard.vue";
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
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
    meta: {
      loginPath: false,
    },
  },
  {
    path: "/signup",
    name: "Signup",
    component: Signup,
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
  const accessToken = TokenService.getLocalAccessToken();
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
  const loginPath = to.matched.some((record) => record.meta.loginPath);

  if (loginPath && accessToken) {
    next("/dashboard");
  }
  // Go to login page if un-authenticated
  if (requiresAuth && !accessToken) {
    next("/login");
  } else next();
});

export default router;
