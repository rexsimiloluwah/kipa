<template>
  <form
    class="flex flex-col space-y-3 z-[1000]"
    @submit.prevent="handleCreateAPIKey"
  >
    <div>
      <div class="relative space-y-1">
        <label for="name">API Key Name<span class="text-red-600">*</span></label>
        <input
          id="name"
          v-model="formState.name"
          type="text"
          name="name"
          placeholder="Enter Name i.e. 'test-api-key'"
          :class="`${v$.name.$errors.length && 'input--error'}`"
          @blur="handleBlur('name')"
        >
      </div>
      <div
        v-for="error of v$.name.$errors"
        :key="error.$uid"
        class="input-errors"
      >
        <p class="text-red-600">
          {{ parseErrorMessage(String(error.$message), "API Key Name") }}
        </p>
      </div>
    </div>

    <div>
      <div class="relative space-y-1">
        <label for="description">Role</label>
        <input
          id="role"
          v-model="formState.role"
          type="text"
          name="role"
          placeholder="Enter Role"
          :class="`${v$.role.$errors.length && 'input--error'}`"
          @blur="handleBlur('role')"
        >
      </div>
      <div
        v-for="error of v$.role.$errors"
        :key="error.$uid"
        class="input-errors"
      >
        <p class="text-red-600">
          {{ parseErrorMessage(String(error.$message), "Role") }}
        </p>
      </div>
    </div>

    <div>
      <label for="expires_at">Expires At<span class="text-red-600">*</span></label>
      <input
        id="expires_at"
        v-model="formState.expires_at"
        type="datetime-local"
        name="expires_at"
        @blur="handleBlur('expires_at')"
      >
    </div>

    <div class="space-y-1">
      <h1>
        Select API Key Permissions ({{
          selectedAPIKeyPermissions.filter((p) => p).length
        }})
      </h1>
      <div class="flex flex-wrap gap-1">
        <button
          v-for="permission in apiKeyPermissions.map((p, i) => ({
            name: p,
            id: i,
          }))"
          :key="permission.id"
          type="button"
          :class="`${
            selectedAPIKeyPermissions[permission.id]
              ? 'bg-primarygreen'
              : 'border-2 border-primarygreen'
          } p-2 rounded-lg`"
          @click="
            selectedAPIKeyPermissions[permission.id] =
              !selectedAPIKeyPermissions[permission.id]
          "
        >
          {{ permission.name }}
          <span
            v-if="selectedAPIKeyPermissions[permission.id]"
            class="text-md font-semibold"
          >x</span>
        </button>
      </div>
    </div>

    <custom-button
      title="Create API Key"
      type="submit"
      :disabled="v$.$invalid"
      :loading="isLoading"
    />
  </form>
</template>

<script lang="ts">
import { defineComponent, ref, reactive } from "vue";
import useVuelidate from "@vuelidate/core";
import { required, maxLength } from "@vuelidate/validators";
import { parseErrorMessage } from "../../../common/utils/form";
import { APIKEY_PERMISSIONS } from "../../../common/constants";
import { useAPIKeyStore } from "../../../store/apikey";

export default defineComponent({
  name: "CreateAPIKeyForm",
  emits: ["closeModal"],
  setup(props, { emit }) {
    const apiKeyStore = useAPIKeyStore();
    const isLoading = ref<boolean>(false);
    const apiKeyPermissions = reactive<Array<string>>(APIKEY_PERMISSIONS);
    const selectedAPIKeyPermissions = ref<Array<boolean>>(
      [...Array(apiKeyPermissions.length)].map(() => true)
    );
    const formState = reactive({
      name: "",
      role: "",
      expires_at: "",
    });
    const rules = {
      name: { required },
      role: {
        maxLength: maxLength(200),
      },
      expires_at: { required },
    };

    const v$ = useVuelidate(rules, formState);

    const handleBlur = (key: "name" | "role" | "expires_at") => {
      // @ts-ignore
      v$.value[key].$dirty = true;
    };

    const handleCreateAPIKey = async () => {
      const data = {
        ...formState,
        expires_at: new Date(formState.expires_at).toISOString(),
        permissions: apiKeyPermissions.filter(
          (_, id) => selectedAPIKeyPermissions.value[id]
        ),
      };
      isLoading.value = true;
      await apiKeyStore.createAPIKey(data).then(() => {
        isLoading.value = false;
        emit("closeModal");
      });
    };

    return {
      formState,
      v$,
      isLoading,
      handleBlur,
      handleCreateAPIKey,
      parseErrorMessage,
      apiKeyPermissions,
      selectedAPIKeyPermissions,
    };
  },
});
</script>
