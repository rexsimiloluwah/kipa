<template>
  <div>
    <h1 class="text-lg font-semibold">
      Are you sure you want to delete API key "{{ apikey.name }}"
    </h1>
    <div class="flex flex-col space-y-3">
      <label for="apiKeyName"
        >Type:
        <span class="font-semibold text-primarygreen">{{ apikey.name }}</span>
        to delete</label
      >
      <input
        type="text"
        name="apiKeyName"
        id="apiKeyName"
        placeholder="Enter the API key name"
        v-model="apiKeyName"
      />
      <button
        :class="`simple__btn bg-red-700 ${
          apiKeyName.trim() !== apikey.name &&
          'cursor-not-allowed opacity-50 pointer-events-none'
        }`"
        @click="handleDeleteAPIKey"
      >
        <span>Delete API Key</span>{{ " " }}
        <font-awesome-icon icon="trash-can" v-if="!isLoading" />
        <font-awesome-icon
          v-if="isLoading"
          icon="spinner"
          class="animate-spin"
        />
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, ref } from "vue";
import Button from "../../shared/Button.vue";
import { useAPIKeyStore } from "../../../store/apikey";
import { APIKey } from "../../../common/types/apikey";

export default defineComponent({
  name: "DeleteAPIKeyDialog",
  emits: ["closeModal"],
  props: {
    apikey: {
      type: Object as PropType<APIKey>,
      required: true,
    },
  },
  components: {
    Button,
  },
  setup(props, { emit }) {
    const apiKeyStore = useAPIKeyStore();
    const apiKeyName = ref<string>("");
    const isLoading = ref<boolean>(false);
    const handleDeleteAPIKey = async () => {
      isLoading.value = true;
      await apiKeyStore.deleteAPIKey(props.apikey.id).then(() => {
        isLoading.value = false;
        emit("closeModal");
      });
    };
    return {
      isLoading,
      apiKeyName,
      handleDeleteAPIKey,
    };
  },
});
</script>
