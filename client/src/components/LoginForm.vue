<template>
  <div class="w-full">
    <div class="text-center mb-6">
      <h1 class="text-2xl font-bold">Log In</h1>
    </div>
    <form class="flex flex-col space-y-3" @submit.prevent="handleLogin">
      <div>
        <div class="relative">
          <label for="name">Email<span class="text-red-600">*</span></label>
          <input
            type="text"
            name="email"
            id="email"
            v-on:blur="handleBlur('email')"
            v-model="formState.email"
            placeholder="Enter E-mail"
            :class="`${v$.email.$errors.length && 'input--error'}`"
          />
        </div>
        <div
          class="input-errors"
          v-for="error of v$.email.$errors"
          :key="error.$uid"
        >
          <p class="text-red-600">
            {{ parseErrorMessage(String(error.$message), "Email") }}
          </p>
        </div>
      </div>

      <div>
        <div class="relative">
          <label for="name">Password<span class="text-red-600">*</span></label>
          <input
            :type="showPassword ? 'text' : 'password'"
            name="password"
            id="password"
            v-model="formState.password"
            v-on:blur="handleBlur('password')"
            placeholder="Enter Password"
            :class="`${v$.password.$errors.length && 'input--error'}`"
          />
          <font-awesome-icon
            :icon="showPassword ? 'eye-slash' : 'eye'"
            class="absolute top-10 right-4 text-gray-700 cursor-pointer"
            @click="toggleShowPassword"
          />
        </div>
        <div
          class="input-errors"
          v-for="error of v$.password.$errors"
          :key="error.$uid"
        >
          <p class="text-red-600">
            {{ parseErrorMessage(String(error.$message), "Password") }}
          </p>
        </div>
      </div>

      <div class="flex justify-between items-center pb-4">
        <router-link to="/forgot-password">Forgot Password</router-link>
        <router-link to="/signup">No Account? Signup</router-link>
      </div>

      <Button
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
import { Button, FormInput } from "./shared";
import { useUserStore } from "../store/user";
import { parseErrorMessage } from "../common/utils/form";

export default defineComponent({
  name: "LoginForm",
  components: { Button, FormInput },
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
    const showPassword = ref<Boolean>(false);
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
