<template>
  <main>
    <NavBar />
    <Hero />
    <About />
    <HowKipaWorks />
    <StayUpdated />
    <FooterSection />
  </main>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { Hero, NavBar } from "../components/shared";
import About from "../components/About.vue";
import HowKipaWorks from "../components/HowKipaWorks.vue";
import StayUpdated from "../components/StayUpdated.vue";
import FooterSection from "../components/Footer.vue";

export default defineComponent({
  name: "Home",
  components: {
    NavBar,
    Hero,
    About,
    HowKipaWorks,
    StayUpdated,
    FooterSection,
  },
  setup() {},
});
</script>

<style scoped></style>
