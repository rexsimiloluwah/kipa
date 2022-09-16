<template>
  <section>
    <Container>
      <div class="text-center space-y-4">
        <h1 class="text-4xl font-bold text-primarygreen">Why Kipa?</h1>
        <p>Turbocharge your development workflows with kipa!</p>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8 my-8">
        <div class="about__card">
          <div class="flex flex-col">
            <font-awesome-icon
              icon="database"
              class="text-5xl text-primarygreen"
            />
            <div class="my-5 space-y-4 text-center">
              <h1 class="text-2xl font-bold">Serverless Store</h1>
              <p class="text-gray-700">
                Kipa provides a serverless key/value store with a simple
                interface for storing and retrieving data.
              </p>
            </div>
          </div>
        </div>
        <div class="about__card">
          <div class="flex flex-col">
            <font-awesome-icon
              icon="file-code"
              class="text-5xl text-primarygreen"
            />
            <div class="my-5 space-y-4 text-center">
              <h1 class="text-2xl font-bold">API Documentation</h1>
              <p class="text-gray-700">
                Kipa provides a simple yet robust API documentation for
                developers to use the Kipa APIs.
              </p>
            </div>
          </div>
        </div>
        <div class="about__card">
          <div class="flex flex-col">
            <font-awesome-icon icon="lock" class="text-5xl text-primarygreen" />
            <div class="my-5 space-y-4 text-center">
              <h1 class="text-2xl font-bold">Security</h1>
              <p class="text-gray-700">
                Kipa is focusing on improving the security of your data using
                API keys, and bucket-level permissions etc.
              </p>
            </div>
          </div>
        </div>
      </div>
    </Container>
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Container from "./shared/Container.vue";

export default defineComponent({
  name: "About",
  components: {
    Container,
  },
  setup() {},
});
</script>
