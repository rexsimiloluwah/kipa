<template>
  <div class="w-full">
    <div class="text-center mb-6">
      <h1 class="text-2xl font-bold">
        Log In
      </h1>
    </div>
    <form
      class="flex flex-col space-y-3"
      @submit.prevent="handleLogin"
    >
      <div>
        <div class="relative">
          <label for="email">Email<span class="text-red-600">*</span></label>
          <input
            id="email"
            v-model="formState.email"
            type="text"
            name="email"
            placeholder="Enter E-mail"
            :class="`${v$.email.$errors.length && 'input--error'}`"
            @blur="handleBlur('email')"
          >
        </div>
        <div
          v-for="error of v$.email.$errors"
          :key="error.$uid"
          class="input-errors"
        >
          <p class="text-red-600">
            {{ parseErrorMessage(String(error.$message), "Email") }}
          </p>
        </div>
      </div>

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

      <div class="flex justify-between items-center pb-4">
        <router-link to="/forgot-password">
          Forgot Password
        </router-link>
        <router-link to="/signup">
          No Account? Signup
        </router-link>
      </div>

      <CustomButton
        title="Log In"
        type="submit"
        :disabled="v$.$invalid"
        :loading="isLoading"
      />
    </form>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from "vue";
import { required, email, minLength } from "@vuelidate/validators";
import { useVuelidate } from "@vuelidate/core";
import { CustomButton } from "./shared";
import { useUserStore } from "../store/user";
import { parseErrorMessage } from "../common/utils/form";

export default defineComponent({
  name: "LoginForm",
  components: { CustomButton },
  setup() {
    const userStore = useUserStore();
    const isLoading = ref<boolean>(false);
    const formState = reactive({
      email: "",
      password: "",
    });
    const rules = {
      email: { required, email },
      password: {
        required,
        minLength: minLength(8),
      },
    };

    const v$ = useVuelidate(rules, formState);

    const handleBlur = (key: "email" | "password") => {
      // @ts-ignore
      v$.value[key].$dirty = true;
    };

    const handleLogin = async () => {
      isLoading.value = true;
      await userStore.loginUser(formState).then(() => {
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
      handleLogin,
      parseErrorMessage,
      showPassword,
      toggleShowPassword,
    };
  },
});
</script>
