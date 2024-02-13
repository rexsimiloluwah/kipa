<template>
  <main class="relative">
    <button
      type="button"
      className="navback__btn group"
      @click="$router.back()"
    >
      <font-awesome-icon
        icon="arrow-left"
        class="text-2xl text-white group-hover:-translate-x-1 duration-200"
      />
    </button>
    <div className="container--center">
      <div
        className="w-[90%] md:w-[30%] about__card py-5 px-8 rounded-xl flex flex-col justify-center"
      >
        <ResetPasswordForm />
      </div>
    </div>
  </main>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import ResetPasswordForm from "../components/ResetPasswordForm.vue";

export default defineComponent({
  name: "ResetPassword",
  components: {
    ResetPasswordForm,
  },
  setup() {},
});
</script>
