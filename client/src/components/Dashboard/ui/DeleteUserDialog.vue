<template>
  <div class="space-y-3">
    <h1 class="text-lg font-semibold">
      Are you sure you want to close your account?
    </h1>
    <p>
      Closing your account will delete your store buckets, items, and delete
      your account entirely.
    </p>
    <div class="flex flex-col space-y-3">
      <label for="deleteMessage">Type:
        <span class="font-semibold text-primarygreen">{{
          DELETE_MESSAGE_TEXT
        }}</span>
        to delete</label>
      <input
        id="deleteMessage"
        v-model="deleteMessage"
        type="text"
        name="deleteMessage"
        :placeholder="`Enter '${DELETE_MESSAGE_TEXT}'`"
      >
      <button
        :class="`simple__btn bg-red-700 ${
          deleteMessage.trim() !== DELETE_MESSAGE_TEXT &&
          'cursor-not-allowed opacity-50 pointer-events-none'
        }`"
        @click="handleDeleteUser"
      >
        <span>Close Account</span>{{ " " }}
        <font-awesome-icon
          v-if="!isLoading"
          icon="trash-can"
        />
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
import { defineComponent, ref } from "vue";
import { useUserStore } from "../../../store/user";

export default defineComponent({
  name: "DeleteUserDialog",
  components: {},
  emits: ["closeModal"],
  setup(props, { emit }) {
    const userStore = useUserStore();
    const DELETE_MESSAGE_TEXT = ref<string>("delete account");
    const deleteMessage = ref<string>("");
    const isLoading = ref<boolean>(false);
    const handleDeleteUser = async () => {
      isLoading.value = true;
      await userStore.deleteUser().then(() => {
        isLoading.value = false;
        emit("closeModal");
      });
    };
    return {
      isLoading,
      deleteMessage,
      handleDeleteUser,
      DELETE_MESSAGE_TEXT,
    };
  },
});
</script>
