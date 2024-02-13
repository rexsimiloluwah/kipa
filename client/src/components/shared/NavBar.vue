<template>
  <div>
    <nav
      class="flex py-4 px-8 md:px-12 justify-between items-center border-[1px] border-b-zinc-200"
    >
      <div class="nav__brand font-bold text-2xl">
        kipa
        <span class="bg-primarygreen w-2 h-2 rounded-full inline-block" />
      </div>
      <ul>
        <li v-if="!userStore.user">
          <a
            href="/"
            target="_blank"
            rel="noreferrer"
          >API Docs</a>
        </li>
        <li v-if="!userStore.user">
          <router-link to="/signup">
            Sign Up
          </router-link>
        </li>
        <li v-if="!userStore.user">
          <router-link to="/login">
            <CustomButton title="Log In" />
          </router-link>
        </li>
        <li>
          <router-link
            v-if="userStore.user"
            to="/dashboard"
          >
            Dashboard
          </router-link>
        </li>
        <li v-if="userStore.user">
          <div class="relative">
            <DefaultProfileImage
              :firstname="userStore.user.firstname"
              :lastname="userStore.user.lastname"
              @click="showAccountDropdown = !showAccountDropdown"
              @mouseenter="showAccountDropdown = true"
              @mouseleave.once="showAccountDropdown = false"
            />
            <div
              v-if="showAccountDropdown"
              class="absolute bg-white top-14 right-0 z-20 rounded-md shadow-xl w-44 p-3 border-[1px] duration-300"
            >
              <ul>
                <li>Logout</li>
              </ul>
            </div>
          </div>
        </li>
      </ul>
      <font-awesome-icon
        :icon="openMobileNavMenu ? 'x' : 'bars'"
        class="md:hidden block text-lg cursor-pointer hover:text-primarygreen"
        @click="openMobileNavMenu = !openMobileNavMenu"
      />
    </nav>
    <div
      v-if="openMobileNavMenu"
      class="w-full z-[1000] py-4 px-8 md:px-12 transition-all duration-300 fixed bg-white md:hidden"
    >
      <ul class="divide-y-[1px] divide-gray-300 md:divide-y-0 mobile">
        <li v-if="!userStore.user">
          <a
            href="/"
            target="_blank"
            rel="noreferrer"
          >API Docs</a>
        </li>
        <li v-if="!userStore.user">
          <router-link to="/signup">
            Sign Up
          </router-link>
        </li>
        <li v-if="!userStore.user">
          <router-link to="/login">
            <CustomCustomButton title="Log In" />
          </router-link>
        </li>
        <li>
          <router-link
            v-if="userStore.user"
            to="/dashboard"
          >
            Dashboard
          </router-link>
        </li>
        <li v-if="userStore.user">
          <div class="relative">
            <DefaultProfileImage
              :firstname="userStore.user.firstname"
              :lastname="userStore.user.lastname"
              @click="showAccountDropdown = !showAccountDropdown"
              @mouseenter="showAccountDropdown = true"
              @mouseleave.once="showAccountDropdown = false"
            />
            <div
              v-if="showAccountDropdown"
              class="absolute bg-white top-0 right-4 z-20 rounded-md shadow-xl w-44 md:p-3 p-1 border-[1px] duration-300"
            >
              <ul>
                <li>Logout</li>
              </ul>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import CustomButton from "./CustomButton.vue";
import DefaultProfileImage from "./DefaultProfileImage.vue";
import { useUserStore } from "../../store/user";

export default defineComponent({
  name: "NavBar",
  components: {
    CustomButton,
    DefaultProfileImage,
  },
  setup() {
    const userStore = useUserStore();
    const showAccountDropdown = ref<boolean>(false);
    const openMobileNavMenu = ref<boolean>(false);

    return {
      userStore,
      showAccountDropdown,
      openMobileNavMenu,
    };
  },
});
</script>

<style scoped></style>
