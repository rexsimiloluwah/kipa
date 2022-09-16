<template>
  <form class="space-y-3 z-[1000]" @submit.prevent="handleUpdateAPIKey">
    <div>
      <div class="relative space-y-1">
        <label for="name"
          >API Key Name<span class="text-red-600">*</span></label
        >
        <input
          type="text"
          name="name"
          id="name"
          v-on:blur="handleBlur('name')"
          v-model="formState.name"
          placeholder="Enter Name i.e. 'test-api-key'"
          :class="`${v$.name.$errors.length && 'input--error'}`"
        />
      </div>
      <div
        class="input-errors"
        v-for="error of v$.name.$errors"
        :key="error.$uid"
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
          type="text"
          name="role"
          id="role"
          v-model="formState.role"
          v-on:blur="handleBlur('role')"
          placeholder="Enter Role"
          :class="`${v$.role.$errors.length && 'input--error'}`"
        />
      </div>
      <div
        class="input-errors"
        v-for="error of v$.role.$errors"
        :key="error.$uid"
      >
        <p class="text-red-600">
          {{ parseErrorMessage(String(error.$message), "Role") }}
        </p>
      </div>
    </div>

    <div>
      <label for="expires_at"
        >Expires At<span class="text-red-600">*</span></label
      >
      <input
        type="datetime-local"
        name="expires_at"
        id="expires_at"
        v-model="formState.expires_at"
        step="1"
        v-on:blur="handleBlur('expires_at')"
      />
    </div>

    <div class="space-y-1">
      <h1>
        Select API Key Permissions ({{
          selectedAPIKeyPermissions.filter((p) => p).length
        }})
      </h1>
      <div class="flex flex-wrap gap-1">
        <button
          type="button"
          v-for="permission in apiKeyPermissions.map((p, i) => ({
            name: p,
            id: i,
          }))"
          :key="permission.id"
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
            >x</span
          >
        </button>
      </div>
    </div>

    <custom-button
      title="Update API Key"
      type="submit"
      :disabled="v$.$invalid"
      :loading="isLoading"
    />
  </form>
</template>

<script lang="ts">
import { defineComponent, ref, reactive, PropType } from "vue";
import useVuelidate from "@vuelidate/core";
import { required, maxLength } from "@vuelidate/validators";
import { parseErrorMessage } from "../../../common/utils/form";
import { APIKEY_PERMISSIONS } from "../../../common/constants";
import { useAPIKeyStore } from "../../../store/apikey";
import { APIKey } from "../../../common/types/apikey";

export default defineComponent({
  name: "CreateAPIKeyForm",
  emits: ["closeModal"],
  props: {
    apikey: {
      type: Object as PropType<APIKey>,
      required: true,
    },
  },
  setup(props, { emit }) {
    const apiKeyStore = useAPIKeyStore();
    const isLoading = ref<boolean>(false);
    const apiKeyPermissions = reactive<Array<string>>(APIKEY_PERMISSIONS);
    const selectedAPIKeyPermissions = ref<Array<boolean>>(
      apiKeyPermissions.map((p) => (props.apikey.permissions || []).includes(p))
    );

    const formState = reactive({
      name: props.apikey.name,
      role: props.apikey.role,
      expires_at: new Date(props.apikey.expires_at)
        .toISOString()
        .replace(".000Z", ""),
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

    const handleUpdateAPIKey = async () => {
      const data = {
        ...formState,
        expires_at: new Date(formState.expires_at).toISOString(),
        permissions: apiKeyPermissions.filter(
          (_, id) => selectedAPIKeyPermissions.value[id]
        ),
      };
      isLoading.value = true;
      await apiKeyStore.updateAPIKey(props.apikey.id, data).then(() => {
        isLoading.value = false;
      });
    };

    return {
      formState,
      v$,
      isLoading,
      handleBlur,
      handleUpdateAPIKey,
      parseErrorMessage,
      apiKeyPermissions,
      selectedAPIKeyPermissions,
    };
  },
});
</script>
