<template>
  <div>
    <h1 class="text-lg font-semibold">
      Are you sure you want to delete bucket "{{ bucket.name }}"? This action deletes all the bucket items under this bucket. 
    </h1>
    <div class="flex flex-col space-y-3">
      <label for="bucketName">Type:
        <span class="font-semibold text-primarygreen">{{ bucket.name }}</span>
        to delete</label>
      <input
        id="bucketName"
        v-model="bucketName"
        type="text"
        name="bucketName"
        placeholder="Enter the bucket name"
      >
      <button
        :class="`simple__btn bg-red-700 ${
          bucketName.trim() !== bucket.name &&
          'cursor-not-allowed opacity-50 pointer-events-none'
        }`"
        @click="handleDeleteBucket"
      >
        <span>Delete Bucket</span>{{ " " }}
        <font-awesome-icon
          v-if="!isLoading"
          icon="trash-can"
        />
        <font-awesome-icon
          v-if="isLoading"
          icon="spinner"
          class="animate-spin"
        />
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, ref } from "vue";
import { BucketDetails } from "../../../common/types/bucket";
import { useBucketStore } from "../../../store/bucket";

export default defineComponent({
  name: "DeleteBucketDialog",
  components: {},
  props: {
    bucket: {
      type: Object as PropType<BucketDetails>,
      required: true,
    },
  },
  emits: ["closeModal"],
  setup(props, { emit }) {
    const bucketStore = useBucketStore();
    const bucketName = ref<string>("");
    const isLoading = ref<boolean>(false);
    const handleDeleteBucket = async () => {
      isLoading.value = true;
      await bucketStore.deleteBucket(props.bucket.uid).then(() => {
        isLoading.value = false;
        emit("closeModal");
      });
    };
    return {
      isLoading,
      bucketName,
      handleDeleteBucket,
    };
  },
});
</script>
