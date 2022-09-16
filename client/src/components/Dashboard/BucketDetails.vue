<template>
  <main v-if="bucket">
    <Modal
      @close="showEditBucketModal = !showEditBucketModal"
      v-if="showEditBucketModal"
      title="Edit Bucket"
    >
      <UpdateBucketForm
        :bucket="bucket"
        @close-modal="showEditBucketModal = !showEditBucketModal"
      />
    </Modal>
    <Modal
      @close="showDeleteBucketModal = !showDeleteBucketModal"
      v-if="showDeleteBucketModal"
    >
      <DeleteBucketDialog
        :bucket="bucket"
        @close-modal="showDeleteBucketModal = !showDeleteBucketModal"
      />
    </Modal>
    <div class="mb-6">
      <div class="flex items-center justify-between">
        <h1 class="text-2xl font-serif font-extrabold text-gray-800">
          Bucket: {{ bucket.name }}
        </h1>
        <Divider />
        <div class="flex justify-end space-x-3">
          <button
            type="button"
            class="simple__btn bg-blue-600"
            @click="showEditBucketModal = true"
          >
            Edit <font-awesome-icon icon="pencil-alt" />
          </button>
          <button
            class="simple__btn bg-red-600"
            @click="showDeleteBucketModal = true"
          >
            Delete
          </button>
        </div>
      </div>
      <p v-if="bucket">{{ bucket.bucket_items.length }} Item(s)</p>
    </div>
    <BucketItemsTable :bucket="bucket" v-if="bucket" />
  </main>
</template>

<script lang="ts">
import { defineComponent, computed, ref } from "vue";
import { NavBack, Divider, Modal } from "../shared";
import { BucketItemsTable, UpdateBucketForm, DeleteBucketDialog } from "./ui";
import { useBucketStore } from "../../store/bucket";
import { useRoute } from "vue-router";

export default defineComponent({
  name: "BucketDetails",
  components: {
    NavBack,
    BucketItemsTable,
    UpdateBucketForm,
    Modal,
    DeleteBucketDialog,
    Divider,
  },
  setup() {
    const route = useRoute();
    const bucketStore = useBucketStore();
    const showEditBucketModal = ref<boolean>(false);
    const showDeleteBucketModal = ref<boolean>(false);

    const bucket = computed(() => bucketStore.bucket(String(route.params.uid)));

    return {
      bucketStore,
      bucket,
      showEditBucketModal,
      showDeleteBucketModal,
    };
  },
});
</script>
