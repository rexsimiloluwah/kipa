<template>
  <div class="p-4 bg-white shadow-lg border-2 border-gray-500 relative">
    <div class="absolute right-4 space-x-2">
      <button
        class="border-[1px] border-gray-600 p-2 rounded-lg hover:border-primarygreen hover:text-primarygreen hover:border-2 hover:font-semibold group"
      >
        Create Item
        <font-awesome-icon icon="add" />
      </button>
      <router-link :to="`/dashboard/bucket/${bucket.uid}`">
        <button
          class="border-[1px] border-gray-600 p-2 rounded-lg bg-primarygreen text-white opacity-80 hover:opacity-100"
        >
          Show Bucket <font-awesome-icon icon="external-link" />
        </button>
      </router-link>
    </div>
    <h1 class="text-xl font-semibold">
      {{ name }}
    </h1>
    <p>{{ description || "No Description." }}</p>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { Bucket } from "../../../common/types/bucket";

export default defineComponent({
  name: "BucketCard",
  components: {},
  props: {
    name: {
      type: String,
      required: true,
    },
    description: {
      type: String,
      required: false,
      default: "",
    },
    itemsCount: {
      type: Number,
      required: false,
      default: 0,
    },
    bucket: {
      type: Object as PropType<Bucket>,
      required: true,
    },
  },
  setup() {},
});
</script>
