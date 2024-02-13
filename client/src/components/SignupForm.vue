<template>
  <div>
    <div class="w-full">
      <div class="text-center mb-6">
        <h1 class="text-2xl font-bold">
          Sign Up
        </h1>
        <p>Create a New Account</p>
      </div>
      <form
        class="flex flex-col space-y-3"
        @submit.prevent="handleSignup"
      >
        <FormRow>
          <div>
            <div class="relative">
              <label for="firstname">Firstname<span class="text-red-600">*</span></label>
              <input
                id="firstname"
                v-model="formState.firstname"
                type="text"
                name="firstname"
                placeholder="Enter Firstname"
                :class="`${v$.firstname.$errors.length && 'input--error'}`"
                @blur="handleBlur('firstname')"
              >
            </div>
            <div
              v-for="error of v$.firstname.$errors"
              :key="error.$uid"
              class="input-errors"
            >
              <p class="text-red-600">
                {{ parseErrorMessage(String(error.$message), "Firstname") }}
              </p>
            </div>
          </div>

          <div>
            <div class="relative">
              <label for="lastname">Lastname<span class="text-red-600">*</span></label>
              <input
                id="lastname"
                v-model="formState.lastname"
                type="text"
                name="lastname"
                placeholder="Enter Lastname"
                :class="`${v$.lastname.$errors.length && 'input--error'}`"
                @blur="handleBlur('lastname')"
              >
            </div>
            <div
              v-for="error of v$.lastname.$errors"
              :key="error.$uid"
              class="input-errors"
            >
              <p class="text-red-600">
                {{ parseErrorMessage(String(error.$message), "Lastname") }}
              </p>
            </div>
          </div>
        </FormRow>
        <div>
          <div class="relative">
            <label for="email">Email<span class="text-red-600">*</span></label>
            <input
              id="email"
              v-model="formState.email"
              type="text"
              name="email"
              placeholder="Enter Email"
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
          <span />
          <router-link to="/login">
            Already registered? Log In
          </router-link>
        </div>

        <CustomButton
          title="Sign Up"
          type="submit"
          :disabled="v$.$invalid"
          :loading="isLoading"
        />
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from "vue";
import { required, email, minLength } from "@vuelidate/validators";
import { useVuelidate } from "@vuelidate/core";
import { CustomButton, FormRow } from "./shared";
import { useUserStore } from "../store/user";
import { parseErrorMessage } from "../common/utils/form";

export default defineComponent({
  name: "SignupForm",
  components: { CustomButton, FormRow },
  setup() {
    const userStore = useUserStore();
    const isLoading = ref<boolean>(false);
    const formState = reactive({
      firstname: "",
      lastname: "",
      email: "",
      password: "",
    });
    const rules = {
      firstname: { required, minLength: minLength(2) },
      lastname: { required, minLength: minLength(2) },
      email: { required, email },
      password: {
        required,
        minLength: minLength(8),
      },
    };

    const v$ = useVuelidate(rules, formState);

    const handleBlur = (
      key: "email" | "password" | "firstname" | "lastname"
    ) => {
      // @ts-ignore
      v$.value[key].$dirty = true;
    };

    const handleSignup = async () => {
      isLoading.value = true;
      await userStore
        .registerUser({ ...formState, username: "theblackdove" })
        .then(() => {
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
      handleSignup,
      parseErrorMessage,
      showPassword,
      toggleShowPassword,
    };
  },
});
</script>
