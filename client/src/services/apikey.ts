import axios from "../lib/axios";
import { APIKey, CreateAPIKeyData } from "../common/types/apikey";
import TokenService from "../services/token";

class APIKeyService {
  /**
   * Create a new API Key
   * @param data // accepts the API key data
   * @returns
   */
  createAPIKey(data: CreateAPIKeyData) {
    return new Promise((resolve, reject) => {
      axios
        .post("/api_key", data)
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Find an API key by ID for the authenticated user
   * @param id // accepts the API Key ID
   * @returns
   */
  findAPIKey(id: string): Promise<APIKey> {
    return new Promise((resolve, reject) => {
      axios
        .get(`/api_key/${id}`)
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Returns the authenticated user's API keys
   * @returns
   */
  getUserAPIKeys(): Promise<Array<APIKey>> {
    return new Promise((resolve, reject) => {
      const token = TokenService.getLocalAccessToken();
      const authHeader = { Authorization: `Bearer ${token}` };
      axios
        .get(`/api_keys`, { headers: authHeader })
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Update a user's API key
   * @param id // accepts the API key ID
   * @param data // accepts the updated API key data
   * @returns
   */
  updateAPIKey(id: string, data: CreateAPIKeyData) {
    return new Promise((resolve, reject) => {
      axios
        .put(`/api_key/${id}`, data)
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Delete an API key
   * @param id // accepts the API key ID
   * @returns
   */
  deleteAPIKey(id: string) {
    return new Promise((resolve, reject) => {
      axios
        .delete(`/api_key/${id}`)
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Revoke an API key
   * @param id // accepts the API key ID
   * @returns
   */
  revokeAPIKey(id: string) {
    return new Promise((resolve, reject) => {
      axios
        .put(`/api_key/${id}/revoke`)
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Delete multiple API keys
   * @param ids // accepts the list of API key IDs
   * @returns
   */
  deleteAPIKeys(ids: Array<string>) {
    return new Promise((resolve, reject) => {
      axios
        .delete(`/api_keys`, {
          data: {
            ids: ids,
          },
        })
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }

  /**
   * Revoke multiple API keys
   * @param ids // accepts the list of API key IDs
   * @returns
   */
  revokeAPIKeys(ids: Array<string>) {
    return new Promise((resolve, reject) => {
      axios
        .put(`api_keys/revoke`, { ids: ids })
        .then((response) => {
          const { data } = response.data;
          resolve(data);
        })
        .catch((error) => {
          reject(error.response.data);
        });
    });
  }
}

export default new APIKeyService();
