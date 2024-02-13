<template>
  <div class="w-full">
    <div class="text-center mb-6">
      <h1 class="text-2xl font-bold">
        Reset Password
      </h1>
      <p class="text-gray-700">
        Set a new secure password for your Kipa account
      </p>
    </div>
    <form
      class="flex flex-col space-y-3"
      @submit.prevent="handleResetPassword"
    >
      <div>
        <div class="relative">
          <label for="password">Password<span class="text-red-600">*</span></label>
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
            :icon="showPassword ? 'eye' : 'eye-slash'"
            class="absolute top-10 right-4 text-gray-700 cursor-pointer"
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

      <div>
        <div class="relative">
          <label for="confirmPassword">Confirm Password<span class="text-red-600">*</span></label>
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
            :icon="showPassword ? 'eye' : 'eye-slash'"
            class="absolute top-10 right-4 text-gray-700 cursor-pointer"
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
        title="Reset Password"
        type="submit"
        :disabled="v$.$invalid"
        :loading="isLoading"
      />
    </form>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from "vue";
import { required, minLength } from "@vuelidate/validators";
import { useVuelidate } from "@vuelidate/core";
import { CustomButton } from "./shared";
import { parseErrorMessage } from "../common/utils/form";

export default defineComponent({
  name: "ResetPasswordForm",
  components: { CustomButton },
  setup() {
    const isLoading = ref<boolean>(false);
    const formState = reactive({
      password: "",
      confirmPassword: "",
    });
    const rules = {
      password: {
        required,
        minLength: minLength(8),
      },
      confirmPassword: {
        required,
        minLength: minLength(8),
      },
    };

    const v$ = useVuelidate(rules, formState);

    const handleBlur = (key: "password" | "confirmPassword") => {
      // @ts-ignore
      v$.value[key].$dirty = true;
    };

    const handleResetPassword = () => {
      console.log("handling reset password");
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
      handleResetPassword,
      parseErrorMessage,
      showPassword,
      toggleShowPassword,
    };
  },
});
</script>
