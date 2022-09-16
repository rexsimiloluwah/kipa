<template>
  <section>
    <Modal
      @close="showCreateAPIKeyModal = !showCreateAPIKeyModal"
      v-if="showCreateAPIKeyModal"
      title="Create API Key"
    >
      <CreateAPIKeyForm
        @close-modal="showCreateAPIKeyModal = !showCreateAPIKeyModal"
      />
    </Modal>
    <div class="flex justify-between items-center">
      <h1 class="md:text-4xl text-2xl font-serif font-extrabold text-gray-800">
        Your API Keys
      </h1>
      <button
        class="
          border-[1px] border-gray-600
          p-2
          rounded-lg
          bg-primarygreen
          text-white
          opacity-80
          hover:opacity-100
        "
        @click="showCreateAPIKeyModal = true"
      >
        Create New API Key <font-awesome-icon icon="add" />
      </button>
    </div>
    <Divider />
    <div class="my-8">
      <div v-if="apiKeyStore.apikeys" class="space-y-4">
        <APIKeyCard
          name="test-key"
          v-for="apikey in apiKeyStore.apikeys"
          :key="apikey.id"
          :apikey="apikey"
        />
      </div>
      <div v-if="apiKeyStore.isLoadingAPIKeys && !apiKeyStore.apikeys.length">
        Loading API keys...
      </div>
    </div>
  </section>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from "vue";
import { Divider, Modal } from "../shared";
import { APIKeyCard, CreateAPIKeyForm } from "./ui";
import { useAPIKeyStore } from "../../store/apikey";

export default defineComponent({
  name: "APIKeys",
  components: {
    Divider,
    APIKeyCard,
    CreateAPIKeyForm,
    Modal,
  },
  setup() {
    const apiKeyStore = useAPIKeyStore();
    const showCreateAPIKeyModal = ref<boolean>(false);

    onMounted(async () => {
      // fetch the user's API keys when this component mounts
      await apiKeyStore.fetchUserAPIKeys();
    });

    return {
      apiKeyStore,
      showCreateAPIKeyModal,
    };
  },
});
</script>