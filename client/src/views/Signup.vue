<template>
  <main class="relative">
    <button
      type="button"
      className="shadow-lg absolute top-5 left-5 w-12 h-12 rounded-full bg-primarygreen hover:opacity-100"
      @click="$router.back()"
    >
      <font-awesome-icon
        icon="arrow-left"
        class="text-2xl text-white"
      />
    </button>
    <div className="container--center">
      <div
        className="w-[90%] md:w-[40%] about__card py-5 px-8 rounded-xl flex flex-col justify-center"
      >
        <SignupForm />
      </div>
    </div>
  </main>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import SignupForm from "../components/SignupForm.vue";

export default defineComponent({
  name: "Signup",
  components: {
    SignupForm,
  },
  setup() {},
});
</script>
