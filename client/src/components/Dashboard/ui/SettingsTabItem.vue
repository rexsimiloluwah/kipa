<template>
  <div
    v-if="selectedTab === title"
    class="bg-white p-5 shadow-lg"
  >
    <slot />
  </div>
</template>

<script lang="ts">
import { defineComponent, inject } from "vue";

export default defineComponent({
  name: "SettingsTabItem",
  props: {
    title: {
      type: String,
      default: "",
    },
  },
  setup() {
    const selectedTab = inject("selectedTab");

    return {
      selectedTab,
    };
  },
});
</script>
