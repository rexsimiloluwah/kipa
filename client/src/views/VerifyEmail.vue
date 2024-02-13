<template>
  <main class="container--center">
    <div
      v-if="!isLoading && error"
      class="w-[90%] shadow-lg bg-red-50 border-2 border-red-600 md:w-[30%] py-5 px-8 rounded-xl flex flex-col justify-center items-center space-y-4 text-center"
    >
      <font-awesome-icon
        icon="warning"
        class="text-5xl text-red-600"
      />
      <h3 class="text-lg md:text-xl font-semibold text-red-60">
        Email Verification Failed
      </h3>
      <p>Something went wrong while verifying your Email.</p>
      <button
        aria-label="resend_email_verification_link"
        class="py-2 px-4 bg-red-600 text-white rounded-lg"
      >
        Resend Verification Link
      </button>
      <router-link
        to="/"
        class="font-semibold underline"
      >
        Go Home
      </router-link>
    </div>
    <div
      v-if="!isLoading && !error"
      class="w-[90%] shadow-lg md:w-[30%] bg-green-50 border-2 border-primarygreen py-5 px-8 rounded-xl flex flex-col justify-center items-center space-y-4 text-center"
    >
      <font-awesome-icon
        :icon="['far', 'check-circle']"
        class="text-5xl text-primarygreen"
      />
      <h3 class="text-lg md:text-xl font-semibold text-primarygreen">
        Email Verification Successful
      </h3>
      <p>
        Your E-mail address has been successfully verified. Thank you for using
        Kipa!
      </p>
      <router-link
        to="/"
        class="font-semibold underline"
      >
        Go Home
      </router-link>
    </div>
  </main>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";

export default defineComponent({
  name: "VerifyEmail",
  setup() {
    const isLoading = ref<boolean>(false);
    const error = ref<string>("");

    return {
      isLoading,
      error,
    };
  },
});
</script>
