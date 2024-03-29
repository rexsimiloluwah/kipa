<template>
  <div class="w-full py-4">
    <h1 class="text-2xl font-bold mb-3">
      Update Password
    </h1>
    <form
      class="space-y-3"
      @submit.prevent="handleUpdatePassword"
    >
      <div>
        <div class="relative">
          <input
            id="password"
            v-model="formState.password"
            :type="showPassword ? 'text' : 'password'"
            name="password"
            placeholder="Enter Password"
            :class="`${v$.password.$errors.length && 'input--error'}`"
            @blur="handleBlur('password')"
          >
          <font-awesome-icon
            :icon="showPassword ? 'eye-slash' : 'eye'"
            class="absolute top-5 right-4 text-gray-700 cursor-pointer"
            @click="toggleShowPassword"
          />
        </div>
        <div
          v-for="error of v$.password.$errors"
          :key="error.$uid"
          class="input-errors"
        >
          <p class="text-red-600">
            {{ parseErrorMessage(String(error.$message), "Password") }}
          </p>
        </div>
      </div>
      <CustomButton
        title="Update Password"
        type="submit"
        :disabled="v$.$invalid"
        :loading="isLoading"
        class="max-w-sm"
      />
    </form>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, reactive } from "vue";
import { required, minLength } from "@vuelidate/validators";
import useVuelidate from "@vuelidate/core";
import { parseErrorMessage } from "../../../common/utils/form";
import { CustomButton } from "../../shared";
import { useUserStore } from "../../../store/user";

export default defineComponent({
  name: "UpdatePasswordForm",
  components: {
    CustomButton,
  },
  setup() {
    const userStore = useUserStore();
    const isLoading = ref<boolean>(false);
    const formState = reactive({
      password: "",
    });
    const rules = {
      password: {
        required,
        minLength: minLength(8),
      },
    };

    const v$ = useVuelidate(rules, formState);

    const handleBlur = (key: "password") => {
      // @ts-ignore
      v$.value[key].$dirty = true;
    };

    const handleUpdatePassword = async () => {
      const data = { ...formState };
      isLoading.value = true;
      await userStore.updateUserPassword(data).then(() => {
        isLoading.value = false;
      });
    };

    // show password control
    const showPassword = ref<boolean>(false);
    const toggleShowPassword = () => {
      showPassword.value = !showPassword.value;
    };

    return {
      formState,
      v$,
      isLoading,
      handleBlur,
      handleUpdatePassword,
      parseErrorMessage,
      showPassword,
      toggleShowPassword,
    };
  },
});
</script>
