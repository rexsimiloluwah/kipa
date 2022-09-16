<template>
  <section>
    <div class="tabs">
      <div
        v-for="tab in tabTitles"
        :key="tab"
        :class="`tab ${tab === selectedTab && 'active'}`"
        @click="selectedTab = tab"
      >
        <span>{{ tab }}</span>
      </div>
    </div>
    <slot />
  </section>
</template>

<script lang="ts">
import { defineComponent, useSlots, ref, provide } from "vue";

export default defineComponent({
  name: "SettingsTabWrapper",
  setup() {
    const slots = useSlots();

    // @ts-ignore
    const tabTitles = ref(slots.default()?.map((tab) => tab.props.title));
    const selectedTab = ref(tabTitles.value[0]);
    provide("selectedTab", selectedTab);

    return {
      tabTitles,
      selectedTab,
    };
  },
});
</script>

