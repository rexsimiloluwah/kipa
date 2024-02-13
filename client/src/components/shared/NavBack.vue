<template>
  <div className="border-b-[1px] sticky top-0 bg-red-600 z-50">
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
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "NavBack",
  setup() {},
});
</script>
