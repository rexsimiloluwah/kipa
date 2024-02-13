<template>
  <main>
    <Modal
      v-if="showDeleteAPIKeyModal"
      @close="showDeleteAPIKeyModal = !showDeleteAPIKeyModal"
    >
      <DeleteAPIKeyDialog
        :apikey="apikey"
        @close-modal="showDeleteAPIKeyModal = !showDeleteAPIKeyModal"
      />
    </Modal>
    <div
      :class="`
        relative
        p-4
        bg-white
        shadow-lg
        border-2 ${
        apiKeyStore.activeKeyID === apikey.id
          ? 'border-green-600 bg-green-100 rounded-lg'
          : 'border-gray-500'
      }
        hover:bg-gray-50
        cursor-pointer
        hover:-translate-y-1
      `"
    >
      <div class="absolute top-2 right-4 space-x-3">
        <button
          v-if="apikey.id === apiKeyStore.activeKeyID"
          type="button"
          class="bg-blue-500 px-2 rounded-md hover:scale-105 duration-100 relative"
          @click="handleCopyKey"
        >
          {{ copied ? "Copied!" : "Copy" }}
        </button>
        <button
          type="button"
          class="bg-primarygreen px-2 rounded-md hover:scale-105 duration-100 relative"
          @click="$router.push(`apikey/${apikey.id}`)"
        >
          Edit
        </button>
        <button
          class="text-red-600 text-xl hover:scale-105 duration-100"
          type="button"
          @click="showDeleteAPIKeyModal = true"
        >
          <font-awesome-icon icon="trash-can" />
        </button>
      </div>

      <div
        v-if="apikey.id !== apiKeyStore.activeKeyID"
        @click="$router.push(`apikey/${apikey.id}`)"
      >
        <h1 class="text-xl font-semibold">
          {{ apikey.name }}
        </h1>
        <p>{{ apikey.role || "No Description." }}</p>
      </div>
      <div v-if="apikey.id === apiKeyStore.activeKeyID">
        <h1 class="text-xl font-semibold break-words max-w-xl text-green-600">
          {{ apikey.name }}
        </h1>
        <h1 class="text-lg font-semibold break-words max-w-xl">
          {{ apiKeyStore.activeKey }}
        </h1>
        <div>
          <span class="text-red-600">*</span>
          Make sure you copy your API key token now, you won't be able to see it
          again!
        </div>
      </div>
    </div>
  </main>
</template>

<script lang="ts">
import { defineComponent, PropType, ref } from "vue";
import { APIKey } from "../../../common/types/apikey";
import DeleteAPIKeyDialog from "./DeleteAPIKeyDialog.vue";
import Modal from "../../shared/Modal.vue";
import { useAPIKeyStore } from "../../../store/apikey";

export default defineComponent({
  name: "APIKeyCard",
  components: {
    Modal,
    DeleteAPIKeyDialog,
  },
  props: {
    apikey: {
      type: Object as PropType<APIKey>,
      required: true,
    },
  },
  setup(props) {
    const apiKeyStore = useAPIKeyStore();
    const showDeleteAPIKeyModal = ref<boolean>(false);
    const copied = ref<boolean>(false);

    // handle copying of the API key
    const handleCopyKey = () => {
      let input = document.createElement("textarea");
      input.innerHTML = props.apikey.key as string;
      document.body.appendChild(input);
      input.select();
      document.execCommand("copy");
      document.body.removeChild(input);
      copied.value = true;
      // toggle copied state after 2s
      setTimeout(() => {
        copied.value = !copied.value;
      }, 2000);
    };

    return {
      copied,
      apiKeyStore,
      handleCopyKey,
      showDeleteAPIKeyModal,
    };
  },
});
</script>
