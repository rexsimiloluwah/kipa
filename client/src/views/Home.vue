<template>
  <main>
    <Navbar />
    <Hero />
    <About />
    <HowKipaWorks />
    <StayUpdated />
    <Footer />
  </main>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Navbar from "../components/shared/Navbar.vue";
import Button from "../components/shared/Button.vue";
import Hero from "../components/shared/Hero.vue";
import About from "../components/About.vue";
import HowKipaWorks from "../components/HowKipaWorks.vue";
import StayUpdated from "../components/StayUpdated.vue";
import Footer from "../components/Footer.vue";

export default defineComponent({
  name: "Home",
  components: {
    Navbar,
    Button,
    Hero,
    About,
    HowKipaWorks,
    StayUpdated,
    Footer,
  },
  setup() {},
});
</script>

<style scoped>
</style>