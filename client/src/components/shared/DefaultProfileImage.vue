<template>
  <div
    class="bg-cyan-600 flex justify-center items-center font-semibold text-white rounded-full w-12 h-12 cursor-pointer opacity-90 hover:opacity-100"
  >
    {{ label }}
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "DefaultProfileImage",
  props: {
    firstname: {
      type: String,
      required: true,
    },
    lastname: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const label =
      props.firstname.charAt(0).toUpperCase() +
      props.lastname.charAt(0).toUpperCase();

    return {
      label,
    };
  },
});
</script>
