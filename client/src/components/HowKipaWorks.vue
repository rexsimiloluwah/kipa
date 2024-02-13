<template>
  <section class="bg-green-50">
    <Container>
      <div class="text-center space-y-4">
        <h1 class="text-4xl font-bold text-primarygreen">
          How Kipa Works
        </h1>
        <p>Create API Key — Create Bucket — 'Kip' anything!</p>
      </div>
      <div class="my-5">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="flex items-center justify-center lg:mt-0">
            <img
              src="/images/guide_img.svg"
              alt="Hero Image"
              class="object-contain h-72 w-full sm:h-80 xl:h-112 2xl:h-128"
            >
          </div>
          <ul
            class="bg-white rounded-md border-[1px] border-gray-700 w-full md:w-[90%] flex flex-col justify-center m-auto divide-y-[1px] divide-gray-700 overflow-hidden relative"
          >
            <li
              v-for="guide in guides"
              :key="guides.indexOf(guide)"
              class="flex py-4 px-3 hover:bg-green-50"
            >
              <div
                class="w-12 h-12 bg-primarygreen rounded-full text-2xl flex justify-center items-center"
              >
                {{ guides.indexOf(guide) + 1 }}
              </div>
              <div class="ml-3 flex-1 overflow-hidden flex-grow">
                <h1 class="font-semibold">
                  {{ guide.title }}
                </h1>
                <p class="text-sm">
                  {{ guide.description }}
                </p>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </Container>
  </section>
</template>

<script lang="ts">
import { defineComponent, reactive } from "vue";
import { Container } from "./shared";

export default defineComponent({
  name: "HowKipaWorks",
  components: {
    Container,
  },
  setup() {
    const guides = reactive([
      {
        title: "Create Account",
        description:
          "Create an Account and API Key with permissions for authentication.",
      },
      {
        title: "Create Bucket",
        description: "Create a Bucket with optional bucket-level permissions.",
      },
      {
        title: "Store, Retrieve, and Update Items etc.",
        description:
          "Create items in a bucket. Items can also be updated (using direct or atomic operations), deleted, or expired etc.",
      },
    ]);

    return {
      guides,
    };
  },
});
</script>

<style scoped></style>
