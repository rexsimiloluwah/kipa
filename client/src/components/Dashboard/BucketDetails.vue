<template>
  <main v-if="bucket">
    <Modal
      v-if="showEditBucketModal"
      title="Edit Bucket"
      @close="showEditBucketModal = !showEditBucketModal"
    >
      <UpdateBucketForm
        :bucket="bucket"
        @close-modal="showEditBucketModal = !showEditBucketModal"
      />
    </Modal>
    <Modal
      v-if="showDeleteBucketModal"
      @close="showDeleteBucketModal = !showDeleteBucketModal"
    >
      <DeleteBucketDialog
        :bucket="bucket"
        @close-modal="showDeleteBucketModal = !showDeleteBucketModal"
      />
    </Modal>
    <Modal
      v-if="showBucketItemDataModal"
      title="Bucket Item Data"
      @close="showBucketItemDataModal =! showBucketItemDataModal"
    >
    <div class="space-y-3 py-3">
      <div class="border-gray-400 border-[1px] p-2">
      {{ bucketItemStore.activeBucketItem?.data }}
    </div>
    <button
          type="button"
          class="simple__btn bg-blue-500"
          @click="handleCopyBucketItemData(bucketItemStore.activeBucketItem?.data)"
        >
          {{ copied ? "Copied!" : "Copy" }}
          <font-awesome-icon icon="clipboard" v-if="!copied"></font-awesome-icon>
      </button>
    </div>
    </Modal>
    <Modal
      v-if="showCreateBucketItemModal"
      title="Create Bucket Item"
      @close="showCreateBucketItemModal = !showCreateBucketItemModal"
    >
      <CreateBucketItemForm :bucketUID="bucket.uid"
        @close-modal="showCreateBucketItemModal = !showCreateBucketItemModal"
      />
    </Modal>
    <Modal
      v-if="showUpdateBucketItemModal"
      title="Edit Bucket Item"
      @close="showUpdateBucketItemModal = !showUpdateBucketItemModal"
    >
      <UpdateBucketItemForm :bucketUID="bucket.uid" @close-modal="showUpdateBucketItemModal = !showUpdateBucketItemModal"/>
    </Modal>
    <Modal
      v-if="showDeleteBucketItemModal"
      title="Delete Bucket Item"
      @close="showDeleteBucketItemModal = !showDeleteBucketItemModal"
    >
      <DeleteBucketItemDialog :bucketUID="bucket.uid" @close-modal="showDeleteBucketItemModal = !showDeleteBucketItemModal" />
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
            class="simple__btn bg-primarygreen"
            @click="showCreateBucketItemModal = true"
          >
            Add Item <font-awesome-icon icon="add" />
          </button>
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
      <p v-if="bucketItemStore.activeBucketItemsResponse?.data">
        {{ bucketItemStore.activeBucketItemsResponse.data.length }} Item(s)
      </p>
    </div>
    <BucketItemsTable
      v-if="bucketItemStore.activeBucketItemsResponse?.data"
      :bucket-items="bucketItemStore.activeBucketItemsResponse?.data"
      @open-bucket-item="showBucketItemData"
      @update-bucket-item="handleEditBucketItem"
      @delete-bucket-item="handleDeleteBucketItem"
    />
    <div class="my-4 space-x-2 flex">
      <button
        v-if="currentPage !== 1"
        :class="`simple__btn bg-blue-600 group`"
        @click="currentPage -= 1"
      >
        <font-awesome-icon
          icon="arrow-left"
          :class="`
            transition-all
            duration-300
            group-hover:translate-x-[-2px] group-hover:font-semibold
          `"
        />
        Previous
      </button>
      <div
        v-if="bucketItemStore.activeBucketItemsResponse?.page_info.total_pages !== 1"
        class="space-x-2"
      >
        <button
          v-for="page in bucketItemStore.activeBucketItemsResponse?.page_info.total_pages"
          :key="page"
          :class="`simple__btn hover:bg-primarygreen hover:text-white ${
            currentPage == page
              ? 'bg-primarygreen text-white'
              : 'bg-white text-black'
          }`"
          @click="currentPage = page"
        >
          {{ page }}
        </button>
      </div>
      <button
        v-if="bucketItemStore.activeBucketItemsResponse?.page_info.has_next_page"
        class="simple__btn bg-blue-600 group"
        @click="currentPage += 1"
      >
        Next
        <font-awesome-icon
          icon="arrow-right"
          class="transition-all duration-300 group-hover:translate-x-[2px] group-hover:font-semibold"
        />
      </button>
    </div>
  </main>
</template>

<script lang="ts">
import {
  defineComponent,
  computed,
  ref,
  onMounted,
  watch,
} from "vue";
import { Divider, Modal } from "../shared";
import { BucketItemsTable, UpdateBucketForm, DeleteBucketDialog } from "./ui";
import { useBucketStore } from "../../store/bucket";
import { useBucketItemStore } from "../../store/bucket_item";
import { useRoute } from "vue-router";
import { BucketItem } from "../../common/types/bucket_item";
import CreateBucketItemForm from "./ui/CreateBucketItemForm.vue";
import UpdateBucketItemForm from "./ui/UpdateBucketItemForm.vue"
import DeleteBucketItemDialog from "./ui/DeleteBucketItemDialog.vue"


export default defineComponent({
  name: "BucketDetails",
  components: {
    BucketItemsTable,
    UpdateBucketForm,
    CreateBucketItemForm,
    UpdateBucketItemForm,
    Modal,
    DeleteBucketDialog,
    DeleteBucketItemDialog,
    Divider,
  },
  setup() {
    const route = useRoute();
    const bucketStore = useBucketStore();
    const bucketItemStore = useBucketItemStore()
    const showEditBucketModal = ref<boolean>(false);
    const showDeleteBucketModal = ref<boolean>(false);
    const showBucketItemDataModal = ref<boolean>(false);
    const showCreateBucketItemModal = ref<boolean>(false);
    const showUpdateBucketItemModal = ref<boolean>(false);
    const showDeleteBucketItemModal = ref<boolean>(false);
    const copied = ref<boolean>(false)
    const currentPage = ref<number>(1);
    const bucketUID = String(route.params.uid);

    const bucket = computed(() => bucketStore.bucket(bucketUID));

    onMounted(async () => {
      await bucketItemStore.fetchBucketItems(bucketUID)
    });

    watch(currentPage, async (currentValue) => {
      await bucketItemStore.fetchBucketItems(bucketUID, {
        page: currentValue,
      })
    });

     // handle copying of the bucket item data
     const handleCopyBucketItemData = (data:string) => {
      let input = document.createElement("textarea");
      input.innerHTML = data as string;
      document.body.appendChild(input);
      input.select();
      document.execCommand("copy");
      document.body.removeChild(input);
      copied.value = true;
      // toggle copied state after 2s
      setTimeout(() => {
        copied.value = !copied.value;
      }, 2000);
    };

    const showBucketItemData = async (bucketItem:BucketItem) => {
      //console.log(bucketItem)
      await bucketItemStore.setActiveBucketItem(bucketItem)
      showBucketItemDataModal.value = true
    }

    const handleEditBucketItem = async (bucketItem:BucketItem) => {
      await bucketItemStore.setActiveBucketItem(bucketItem)
      showUpdateBucketItemModal.value = true
    }

    const handleDeleteBucketItem = async (bucketItem:BucketItem) => {
      await bucketItemStore.setActiveBucketItem(bucketItem)
      showDeleteBucketItemModal.value = true
    }

    return {
      currentPage,
      bucketStore,
      bucketItemStore,
      bucket,
      showEditBucketModal,
      showDeleteBucketModal,
      showBucketItemData,
      showDeleteBucketItemModal,
      showBucketItemDataModal,
      showCreateBucketItemModal,
      showUpdateBucketItemModal,
      copied,
      handleCopyBucketItemData,
      handleEditBucketItem,
      handleDeleteBucketItem
    };
  },
});
</script>
