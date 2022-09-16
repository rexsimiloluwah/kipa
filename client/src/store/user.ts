import { defineStore } from "pinia";
import { CreateUserData, UpdateUserData, User } from "../common/types/user";
import AuthService from "../services/auth";
import UserService from "../services/user";
import { useToast } from "vue-toastification";
import router from "../router";

const toast = useToast();

export const useUserStore = defineStore("user", {
  state: () => ({
    user: null as User | null,
  }),
  getters: {},
  actions: {
    async fetchUser() {
      try {
        const user = await AuthService.getAuthUser();
        this.user = user;
      } catch {
        this.user = null;
      }
    },

    async registerUser(data: CreateUserData) {
      try {
        await AuthService.register(data);
        toast.success("Sign up Successful! 🎉");
        router.push("/login");
      } catch (error: any) {
        // console.log(error);
        toast.error(error.error || error.message);
      }
    },

    async loginUser(data: { email: string; password: string }) {
      try {
        await AuthService.login(data);
        toast.success("Log In Successful! 🎉");
        await this.fetchUser();
        setTimeout(() => router.push("/dashboard"), 1000);
      } catch (error: any) {
        //console.log(error);
        toast.error(error.error || error.message);
      }
    },
    async updateUser(data: UpdateUserData) {
      try {
        await UserService.updateUser(data);
        toast.success("User Updated Successfully! 🎉");
        await this.fetchUser();
      } catch (error: any) {
        toast.error(error.error || error.message);
      }
    },
    async updateUserPassword(data: { password: string }) {
      try {
        await UserService.updateUserPassword(data);
        toast.success("Password Updated Successfully! 🎉");
      } catch (error: any) {
        toast.error(error.error || error.message);
      }
    },
    async deleteUser() {
      try {
        await UserService.deleteUser();
        toast.success("User Account Deleted Successfully! 🎉");
        router.push("/login");
      } catch (error: any) {
        toast.error(error.error || error.message);
      }
    },
  },
});
