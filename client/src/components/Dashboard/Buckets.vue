<template>
  <section>
    <Modal
      @close="showCreateBucketModal = !showCreateBucketModal"
      v-if="showCreateBucketModal"
      title="Create Bucket"
    >
      <CreateBucketForm
        @close-modal="showCreateBucketModal = !showCreateBucketModal"
      />
    </Modal>
    <div class="mb-8 space-y-4">
      <h1 class="md:text-4xl text-2xl font-serif font-extrabold text-gray-800">
        <span aria-label="robot" role="img">🤖</span> Howdy,
        <span
          >{{ userStore.user?.firstname }} {{ userStore.user?.lastname }}</span
        >
      </h1>
      <p class="text-md">Welcome to your Kipa Dasboard</p>
    </div>

    <div class="grid grid-cols-2 gap-8">
      <HomeAnalyticsCard
        title="Buckets Created"
        icon="database"
        :count="bucketStore.bucketsCount"
      />
      <HomeAnalyticsCard
        title="Items Kept"
        icon="database"
        :count="bucketStore.bucketItemsCount"
      />
    </div>

    <div class="my-8">
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-serif font-extrabold text-gray-800">
          Your Buckets
        </h1>
        <button
          class="
            border-[1px] border-gray-600
            p-2
            rounded-lg
            bg-primarygreen
            text-white
            opacity-80
            hover:opacity-100
          "
          @click="showCreateBucketModal = true"
        >
          Create Bucket <font-awesome-icon icon="add" />
        </button>
      </div>
      <Divider />

      <div class="buckets-list__container flex flex-col space-y-3 my-4">
        <div v-if="bucketStore.isLoadingBuckets">Loading Buckets...</div>
        <BucketCard
          v-for="bucket in bucketStore.buckets"
          :key="bucket.uid"
          :name="bucket.name"
          :bucket="bucket"
        />
      </div>
    </div>
  </section>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import { useUserStore } from "../../store/user";
import { useBucketStore } from "../../store/bucket";
import { HomeAnalyticsCard, BucketCard, CreateBucketForm } from "./ui";
import { Divider, Modal } from "../shared";

export default defineComponent({
  name: "Buckets",
  components: {
    HomeAnalyticsCard,
    Divider,
    BucketCard,
    Modal,
    CreateBucketForm,
  },
  setup() {
    const userStore = useUserStore();
    const bucketStore = useBucketStore();
    const showCreateBucketModal = ref<boolean>(false);

    return {
      userStore,
      bucketStore,
      showCreateBucketModal,
    };
  },
});
</script>
