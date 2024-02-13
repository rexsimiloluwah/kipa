<template>
  <section>
    <Container>
      <div
        class="bg-green-50 space-y-8 p-8 rounded-md border-[1px] border-gray-500"
      >
        <div class="text-center space-y-3">
          <h1 class="text-4xl text-primarygreen font-bold">
            Stay Updated
          </h1>
          <p class="text-gray-500">
            Be the first to when new features are released!
          </p>
        </div>
        <div
          class="m-auto flex flex-col justify-center w-[95%] md:w-[80%] rounded-xl border-[1px] border-black overflow-hidden"
        >
          <input
            v-model="formState.email"
            type="text"
            class="flex-1 py-3 px-2 rounded-xl border-0"
            placeholder="Enter E-mail Address"
          >
          <CustomButton
            title="Subscribe"
            class="border-none rounded-none"
            :disabled="v$.$invalid"
          />
        </div>
      </div>
    </Container>
  </section>
</template>

<script lang="ts">
import { defineComponent, reactive } from "vue";
import { CustomButton, Container } from "./shared";
import { required, email } from "@vuelidate/validators";
import { useVuelidate } from "@vuelidate/core";

export default defineComponent({
  name: "StayUpdated",
  components: {
    CustomButton,
    Container,
  },
  setup() {
    const formState = reactive({
      email: "",
    });
    const rules = {
      email: { required, email },
    };

    const v$ = useVuelidate(rules, formState);

    return {
      formState,
      v$,
    };
  },
});
</script>
