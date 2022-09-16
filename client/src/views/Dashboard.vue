<template>
  <main class="block md:flex justify-between bg-green-50 h-screen">
    <div class="md:hidden flex justify-center items-center h-screen">
      <h1>Not Yet Optimized for Mobile Responsiveness! 😐</h1>
    </div>
    <div class="hidden md:block h-screen p-12 border-r-[1px] border-gray-200">
      <LeftSidebar />
    </div>
    <section class="hidden md:block flex-grow overflow-y-scroll relative p-12">
      <router-view></router-view>
    </section>
  </main>
</template>

<script lang="ts">
import { defineComponent, onMounted } from "vue";
import LeftSidebar from "../components/Dashboard/LeftSidebar.vue";
import { useBucketStore } from "../store/bucket";

export default defineComponent({
  setup() {
    const bucketStore = useBucketStore();

    onMounted(async () => {
      // fetch the user's buckets when this component mounts
      await bucketStore.fetchBuckets();
    });
  },
  components: { LeftSidebar },
});
</script>
