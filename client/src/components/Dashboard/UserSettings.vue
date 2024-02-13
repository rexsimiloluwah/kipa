<template>
  <main>
    <Modal
      v-if="showDeleteAccountModal"
      title="Delete Account"
      @close="showDeleteAccountModal = !showDeleteAccountModal"
    >
      <DeleteUserDialog />
    </Modal>
    <h1 class="md:text-4xl text-2xl font-serif font-extrabold text-gray-800">
      Settings
    </h1>
    <SettingsTabWrapper>
      <SettingsTabItem title="Account">
        <div class="divide-gray-100 divide-y-[1px] space-y-4">
          <UpdateAccountForm />
          <UpdatePasswordForm />
          <div class="py-3 space-y-3">
            <h1 class="text-2xl font-bold">
              Close Account
            </h1>
            <p>
              NOTE: Closing your account implies deleting all your stored
              buckets, items, API keys, and user account data.
            </p>
            <CustomButton
              title="Close Account"
              :primary="false"
              class="bg-red-500 hover:bg-red-700 text-white border-none"
              @click="showDeleteAccountModal = true"
            >
              <font-awesome-icon icon="trash-can" />
            </CustomButton>
          </div>
        </div>
      </SettingsTabItem>
      <SettingsTabItem title="Subscriptions">
        <div>Still Cooking!</div>
      </SettingsTabItem>
    </SettingsTabWrapper>
  </main>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import {
  SettingsTabWrapper,
  SettingsTabItem,
  UpdateAccountForm,
  UpdatePasswordForm,
  DeleteUserDialog,
} from "./ui";
import { CustomButton, Modal } from "../shared";

export default defineComponent({
  name: "UserSettings",
  components: {
    Modal,
    CustomButton,
    SettingsTabWrapper,
    SettingsTabItem,
    UpdateAccountForm,
    UpdatePasswordForm,
    DeleteUserDialog,
  },
  setup() {
    const showDeleteAccountModal = ref<boolean>(false);

    return {
      showDeleteAccountModal,
    };
  },
});
</script>
