<template>
  <main class="relative">
    <button
      type="button"
      @click="$router.back()"
      className="shadow-lg absolute top-5 left-5 w-12 h-12 rounded-full bg-primarygreen hover:opacity-100"
    >
      <font-awesome-icon icon="arrow-left" class="text-2xl text-white" />
    </button>
    <div
      className="m-auto min-h-screen flex items-center justify-center flex-col font-uchen"
    >
      <div
        className="w-[90%] md:w-[30%] about__card py-5 px-8 rounded-xl flex flex-col justify-center"
      >
        <LoginForm />
      </div>
    </div>
  </main>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import LoginForm from "../components/LoginForm.vue";

export default defineComponent({
  name: "Login",
  components: {
    LoginForm,
  },
  setup() {},
});
</script>
