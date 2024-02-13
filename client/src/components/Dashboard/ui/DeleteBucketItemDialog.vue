<template>
    <div class="space-y-3">
      <h1 class="text-md">
        Are you sure you want to delete bucket item `{{ bucketItemStore.activeBucketItem?.key }}` in bucket `{{ bucketItemStore.activeBucketItem?.bucket_uid }}`? 
      </h1>
      <div class="flex space-x-3">
        <button class="simple__btn bg-red-500 px-4" @click="handleDeleteBucketItem">
            Yes
        </button>
        <button class="simple__btn bg-blue-500 px-4" @click="handleCloseModal">
            No
        </button>
      </div>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, PropType, ref } from "vue";
  import { useBucketItemStore } from "../../../store/bucket_item";
  
  export default defineComponent({
    name: "DeleteBucketItemDialog",
    components: {},
    props: {
      bucketUID: {
        type: String as PropType<string>,
        required: true,
      },
    },
    emits: ["closeModal"],
    setup(props, { emit }) {
      const isLoading = ref<boolean>(false);
      const bucketItemStore = useBucketItemStore()

      const handleDeleteBucketItem = async () => {
        isLoading.value = true;
        await bucketItemStore.deleteBucketItem(props.bucketUID,bucketItemStore.activeBucketItem?.key as string).then(() => {
            isLoading.value = false
            emit("closeModal")
        })
      };

      const handleCloseModal = () => emit("closeModal")

      return {
        isLoading,
        bucketItemStore,
        handleCloseModal,
        handleDeleteBucketItem,
      };
    },
  });
  </script>
  