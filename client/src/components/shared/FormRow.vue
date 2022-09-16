<template>
  <div className="w-full grid grid-cols-1 md:grid-cols-2 gap-3">
    <slot />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "FormRow",
  setup() {
    return;
  },
});
</script>
