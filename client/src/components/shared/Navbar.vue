<template>
  <nav
    class="
      flex
      py-4
      px-12
      justify-between
      items-center
      border-[1px] border-b-zinc-200
    "
  >
    <div class="nav__brand font-bold text-2xl">
      kipa
      <span class="bg-primarygreen w-2 h-2 rounded-full inline-block"></span>
    </div>
    <ul>
      <li>
        <a href="/" target="_blank" rel="noreferrer">API Docs</a>
      </li>
      <li>
        <router-link to="/signup">Sign Up</router-link>
      </li>
      <li>
        <router-link to="/login"><Button title="Log In" /></router-link>
      </li>
      <li><router-link to="/dashboard">Dashboard</router-link></li>
      <li>
        <div class="relative">
          <DefaultProfileImage
            firstname="Similoluwa"
            lastname="Okunowo"
            @click="showAccountDropdown = !showAccountDropdown"
            @mouseenter="showAccountDropdown = true"
            @mouseleave.once="showAccountDropdown = false"
          />
          <div
            class="
              absolute
              bg-white
              top-14
              right-0
              z-20
              rounded-md
              shadow-xl
              w-44
              p-3
              border-[1px]
              duration-300
            "
            v-if="showAccountDropdown"
          >
            <ul>
              <li>Logout</li>
            </ul>
          </div>
        </div>
      </li>
    </ul>
  </nav>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import Button from "./Button.vue";
import DefaultProfileImage from "./DefaultProfileImage.vue";
import { useUserStore } from "../../store/user";

export default defineComponent({
  name: "Navbar",
  components: {
    Button,
    DefaultProfileImage,
  },
  setup() {
    const userStore = useUserStore();
    const showAccountDropdown = ref<boolean>(false);

    return {
      userStore,
      showAccountDropdown,
    };
  },
});
</script>


<style scoped>
</style>