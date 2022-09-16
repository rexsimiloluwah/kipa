<template>
  <section>
    <Modal @close="showModal = !showModal" v-if="showModal">
      Lorem ipsum dolor sit amet consectetur adipisicing elit. Quod esse
      perspiciatis minus, repellat nihil voluptatum, laborum assumenda nobis
      fuga, quae ut inventore. Adipisci esse, quasi quaerat nulla ipsum earum?
      Nemo.
    </Modal>
    <div class="relative">
      <div class="bg-green-100 overflow-hidden hero pb-8">
        <Container>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div
              class="
                flex
                justify-center
                items-center
                md:justify-items-start md:items-start
                m-0
                space-y-6
                flex-col
              "
            >
              <h1 class="font-bold space-x-3">
                <span class="text-6xl">KIPA</span>
                <span class="text-gray-600 text-xl">(/ˈki·pər/)</span>
              </h1>
              <p>
                A secure {{ "&" }} serverless key/value store for modern
                applications
              </p>
              <div class="space-x-4">
                <Button title="Get Started" />
                <Button
                  title="View API Docs"
                  :primary="false"
                  @click="showModal = true"
                >
                  <font-awesome-icon icon="gear" />
                </Button>
              </div>
            </div>
            <div class="flex items-center justify-center lg:mt-0">
              <img
                src="/images/hero_img.svg"
                alt="Hero Image"
                class="object-contain h-72 w-full sm:h-80 xl:h-112 2xl:h-128"
              />
            </div>
          </div>
        </Container>
      </div>
    </div>
  </section>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import Container from "./Container.vue";
import Button from "./Button.vue";
import Modal from "./Modal.vue";

export default defineComponent({
  name: "Hero",
  components: {
    Container,
    Button,
    Modal,
  },
  setup() {
    const showModal = ref<Boolean>(false);
    return {
      showModal,
    };
  },
});
</script>

<style scoped>
</style>