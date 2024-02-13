<template>
  <div class="w-full">
    <div class="text-center mb-6 space-y-3">
      <h1 class="text-2xl font-bold">
        Forgot Password
      </h1>
      <p class="text-gray-700">
        Enter your Kipa E-mail address. A reset password link will be sent your
        mail in if it exists in our database.
      </p>
    </div>
    <form
      class="flex flex-col space-y-3"
      @submit.prevent="handleForgotPasswordSubmit"
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

      <CustomButton
        title="Submit E-mail"
        type="submit"
        :disabled="v$.$invalid"
        :loading="isLoading"
      />
    </form>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from "vue";
import { required, email } from "@vuelidate/validators";
import { useVuelidate } from "@vuelidate/core";
import { CustomButton } from "./shared";
import { parseErrorMessage } from "../common/utils/form";

export default defineComponent({
  name: "ForgotPasswordForm",
  components: { CustomButton },
  setup() {
    const isLoading = ref<boolean>(false);
    const formState = reactive({
      email: "",
    });
    const rules = {
      email: { required, email },
    };

    const v$ = useVuelidate(rules, formState);

    const handleBlur = (key: "email") => {
      // @ts-ignore
      v$.value[key].$dirty = true;
    };

    const handleForgotPasswordSubmit = () => {
      console.log("handling forgot password submission");
    };

    return {
      formState,
      v$,
      isLoading,
      handleBlur,
      handleForgotPasswordSubmit,
      parseErrorMessage,
    };
  },
});
</script>
