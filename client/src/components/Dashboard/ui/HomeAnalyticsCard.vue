<template>
  <div
    class="
      bg-white
      py-4
      px-8
      rounded-lg
      space-y-3
      offset-shadow
      border-2 border-gray-500
      relative
      hover:translate-y-[-3px]
    "
  >
    <div
      class="
        bg-green-50
        w-16
        h-16
        flex
        items-center
        justify-center
        rounded-full
        absolute
        top-2
        right-4
      "
    >
      <font-awesome-icon :icon="icon" class="text-4xl text-primarygreen" />
    </div>
    <h1 class="text-7xl font-serif font-extrabold text-primarygreen">
      {{ count }}
    </h1>
    <p>{{ title }}</p>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "HomeAnalyticsCard",
  props: {
    title: {
      type: String,
      required: true,
    },
    count: {
      type: Number,
      required: false,
      default: 0,
    },
    icon: {
      type: String,
      required: true,
    },
  },
  setup() {},
});
</script>
