<template>
  <div>
    <div class="w-full">
      <div class="text-center mb-6">
        <h1 class="text-2xl font-bold">Sign Up</h1>
        <p>Create a New Account</p>
      </div>
      <form class="flex flex-col space-y-3" @submit.prevent="handleSignup">
        <FormRow>
          <div>
            <div class="relative">
              <label for="name"
                >Firstname<span class="text-red-600">*</span></label
              >
              <input
                type="text"
                name="firstname"
                id="firstname"
                v-on:blur="handleBlur('firstname')"
                v-model="formState.firstname"
                placeholder="Enter Firstname"
                :class="`${v$.firstname.$errors.length && 'input--error'}`"
              />
            </div>
            <div
              class="input-errors"
              v-for="error of v$.firstname.$errors"
              :key="error.$uid"
            >
              <p class="text-red-600">
                {{ parseErrorMessage(String(error.$message), "Firstname") }}
              </p>
            </div>
          </div>

          <div>
            <div class="relative">
              <label for="name"
                >Lastname<span class="text-red-600">*</span></label
              >
              <input
                type="text"
                name="lastname"
                id="lastname"
                v-on:blur="handleBlur('lastname')"
                v-model="formState.lastname"
                placeholder="Enter Lastname"
                :class="`${v$.lastname.$errors.length && 'input--error'}`"
              />
            </div>
            <div
              class="input-errors"
              v-for="error of v$.lastname.$errors"
              :key="error.$uid"
            >
              <p class="text-red-600">
                {{ parseErrorMessage(String(error.$message), "Lastname") }}
              </p>
            </div>
          </div>
        </FormRow>
        <div>
          <div class="relative">
            <label for="name">Email<span class="text-red-600">*</span></label>
            <input
              type="text"
              name="email"
              id="email"
              v-on:blur="handleBlur('email')"
              v-model="formState.email"
              placeholder="Enter Email"
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
            <label for="name"
              >Password<span class="text-red-600">*</span></label
            >
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
          <span></span>
          <router-link to="/login">Already registered? Log In</router-link>
        </div>

        <Button
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
import { Button, FormRow, FormInput } from "./shared";
import { useUserStore } from "../store/user";
import { parseErrorMessage } from "../common/utils/form";

export default defineComponent({
  name: "SignupForm",
  components: { Button, FormInput, FormRow },
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
      await userStore.registerUser(formState).then(() => {
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
      handleSignup,
      parseErrorMessage,
      showPassword,
      toggleShowPassword,
    };
  },
});
</script>