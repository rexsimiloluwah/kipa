<template>
  <section>
    <Container>
      <div
        class="
          bg-green-50
          space-y-8
          p-8
          rounded-md
          border-[1px] border-gray-500
        "
      >
        <div class="text-center space-y-3">
          <h1 class="text-4xl text-primarygreen font-bold">Stay Updated</h1>
          <p class="text-gray-500">
            Be the first to when new features are released!
          </p>
        </div>
        <div
          class="
            m-auto
            flex flex-col
            justify-center
            w-[95%]
            md:w-[80%]
            rounded-xl
            border-[1px] border-black
            overflow-hidden
          "
        >
          <input
            type="text"
            class="flex-1 py-3 px-2 rounded-xl border-0"
            placeholder="Enter E-mail Address"
          />
          <Button title="Subscribe" class="border-none rounded-none" />
        </div>
      </div>
    </Container>
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { Button, Container } from "./shared";

export default defineComponent({
  name: "StayUpdated",
  components: {
    Button,
    Container,
  },
  setup() {},
});
</script>
