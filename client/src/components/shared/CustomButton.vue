<template>
  <button
    :type="type"
    :class="`
      py-3
      px-4
      rounded-md
      hover:opacity-100
      group
      ${
      primary
        ? `bg-primarygreen text-white border-[1px] border-black
      hover:border-2`
        : `bg-white text-primarygreen  hover:border-black
      group
      border-2 border-primarygreen`
    }
      ${
      disabled || loading
        ? `cursor-not-allowed opacity-50 pointer-events-none`
        : `cursor-pointer opacity-85`
    }
    `"
  >
    {{ title }}
    <span v-if="!loading">
      <font-awesome-icon
        v-if="primary"
        icon="arrow-right"
        class="transition-all duration-300 group-hover:translate-x-2 group-hover:font-semibold"
      />
      <slot />
    </span>
    <font-awesome-icon
      v-if="loading"
      icon="spinner"
      class="animate-spin"
    />
  </button>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";

export default defineComponent({
  name: "CustomButton",
  props: {
    title: {
      required: true,
      type: String,
    },
    type: {
      type: String as PropType<"button" | "submit" | "reset">,
      default: "button",
    },
    loading: {
      type: Boolean,
      default: false,
    },
    primary: {
      type: Boolean,
      default: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup() {},
});
</script>

<style scoped></style>
