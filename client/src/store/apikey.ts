import { defineStore } from "pinia";
import { APIKey, CreateAPIKeyData } from "../common/types/apikey";
import APIKeyService from "../services/apikey";
import { useToast } from "vue-toastification";

const toast = useToast();

export const useAPIKeyStore = defineStore("apikey", {
  state: () => ({
    apikeys: [] as Array<APIKey>,
    isLoadingAPIKeys: false as boolean,
    activeKeyID: "" as string, // for storing the state of a newly created API key ID
    activeKey: "" as string, // for storing the state of the newly created API key
  }),
  getters: {
    apikey: (state) => {
      return (id: string) => state.apikeys.filter((key) => key.id === id)[0];
    },
  },
  actions: {
    async fetchUserAPIKeys() {
      try {
        this.isLoadingAPIKeys = true;
        const apiKeys = await APIKeyService.getUserAPIKeys();
        console.log(apiKeys);
        this.apikeys = apiKeys;
        this.isLoadingAPIKeys = false;
      } catch (error: any) {
        toast.error(error.error || error.message);
        this.apikeys = [];
        this.isLoadingAPIKeys = false;
      }
    },
    async createAPIKey(data: CreateAPIKeyData) {
      try {
        const newAPIKey = await APIKeyService.createAPIKey(data);
        toast.success("Successfully created API key.");
        this.fetchUserAPIKeys();
        this.activeKeyID = newAPIKey.id;
        this.activeKey = newAPIKey.key;
      } catch (error: any) {
        toast.error(error.error || error.message);
      }
    },
    async updateAPIKey(id: string, data: CreateAPIKeyData) {
      try {
        await APIKeyService.updateAPIKey(id, data);
        toast.success("Successfully updated API key.");
        this.fetchUserAPIKeys();
      } catch (error: any) {
        toast.error(error.error || error.message);
      }
    },
    async deleteAPIKey(id: string) {
      try {
        await APIKeyService.deleteAPIKey(id);
        toast.success("API Key Deleted Successfully! 🎉");
        this.fetchUserAPIKeys();
      } catch (error: any) {
        toast.error(error.error || error.message);
      }
    },
    async revokeAPIKey(id: string) {
      try {
        await APIKeyService.revokeAPIKey(id);
        toast.success("Successfully revoked API key.");
        this.fetchUserAPIKeys();
      } catch (error: any) {
        toast.error(error.error || error.message);
      }
    },
    async deleteAPIKeys() {},
    async revokeAPIKeys() {},
  },
});
