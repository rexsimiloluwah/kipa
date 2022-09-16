<template>
  <router-view></router-view>
</template>

<script lang="ts">
import { defineComponent, onMounted } from "vue";
import { useUserStore } from "./store/user";

export default defineComponent({
  name: "App",
  setup() {
    const userStore = useUserStore();
    onMounted(async () => {
      // fetch the user on mount
      await userStore.fetchUser();
    });
  },
});
</script>