<template>
  <main v-if="apiKey">
    <div class="">
      <h1 class="text-2xl font-serif font-extrabold text-gray-800">
        API Key Details
      </h1>
      <Divider />
    </div>
    <div class="space-x-3">
      <button
        v-if="apiKey.revoked"
        class="bg-orange-400 rounded-lg p-1"
      >
        status:revoked
      </button>
      <button
        v-if="isAPIKeyExpired"
        class="bg-red-400 rounded-lg p-1"
      >
        status:expired
      </button>
      <button
        v-if="isAPIKeyActive"
        class="bg-green-400 rounded-lg p-1"
      >
        status:active
      </button>
    </div>
    <div class="my-4 space-y-4">
      <div class="border-2 px-4 py-2 rounded-md">
        <UpdateAPIKeyForm :apikey="apiKey" />
      </div>
      <div class="border-2 px-4 py-2 rounded-md space-y-3">
        <h1 class="text-xl font-serif font-extrabold text-gray-800">
          Revoke API Key
        </h1>
        <p>
          Note: If you revoke this key, be aware that the key will become
          unusable, and any scripts or applications using it will also need to
          be udpated. You can revoke this key and create a new API key if you
          have lost, forgotten, or exposed this key (for optimal security).
        </p>
        <custom-button
          title="Revoke API Key"
          type="button"
          class="bg-red-500"
          :loading="isRevoking"
          @click="handleRevokeAPIKey"
        />
      </div>
    </div>
  </main>
</template>

<script lang="ts">
import { computed, defineComponent, ref } from "vue";
import { Divider } from "../shared";
import { UpdateAPIKeyForm } from "./ui";
import { useAPIKeyStore } from "../../store/apikey";
import { useRoute } from "vue-router";

export default defineComponent({
  name: "APIKeyDetails",
  components: {
    Divider,
    UpdateAPIKeyForm,
  },
  setup() {
    const route = useRoute();
    const apiKeyStore = useAPIKeyStore();
    const apiKey = apiKeyStore.apikey(String(route.params.id));
    const isRevoking = ref<boolean>(false);

    // compute if the api key has expired
    const isAPIKeyExpired = computed(
      () => new Date(apiKey.expires_at).getTime() < new Date().getTime()
    );

    // compute if the api key is active --> not revoked and not expired
    const isAPIKeyActive = computed(
      () => !apiKey.revoked && !isAPIKeyExpired.value
    );

    const handleRevokeAPIKey = async () => {
      isRevoking.value = true;
      await apiKeyStore.revokeAPIKey(apiKey.id).then(() => {
        isRevoking.value = false;
      });
    };

    return {
      apiKey,
      isRevoking,
      apiKeyStore,
      isAPIKeyExpired,
      isAPIKeyActive,
      handleRevokeAPIKey,
    };
  },
});
</script>
