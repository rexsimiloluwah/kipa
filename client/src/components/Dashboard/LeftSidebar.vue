<template>
  <div
    class="min-w-[250px] rounded-lg border-[1px] border-gray-500 bg-white shadow-lg overflow-hidden"
  >
    <div class="nav__brand font-bold text-2xl p-3 bg-gray-50">
      <Logo />
    </div>
    <ul>
      <li
        v-for="link in links"
        :key="links.indexOf(link)"
        :class="`
              text-xl
              ${$route.path === link.path && 'dashboard_nav--active'}
              p-4
              cursor-pointer
              hover:bg-green-50
            `"
        @click="$router.push(link.path)"
      >
        {{ link.title }}
      </li>
      <li
        :class="`
              text-xl
              p-4
              cursor-pointer
              hover:bg-green-50
            `"
        @click="userStore.logoutUser()"
      >
        Logout
      </li>
    </ul>
    <div
      class="p-3 border-t-[1px] border-gray-500 hover:text-primarygreen cursor-pointer"
      @click="$router.back()"
    >
      <font-awesome-icon icon="arrow-left" /> Go Back
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive } from "vue";
import { Logo } from "../shared";
import { useUserStore } from "../../store/user";

export default defineComponent({
  name: "LeftSidebar",
  components: {
    Logo,
  },
  setup() {
    const userStore = useUserStore();
    const links = reactive([
      {
        title: "Home",
        path: "/dashboard",
      },
      {
        title: "Settings",
        path: "/dashboard/user-settings",
      },
      {
        title: "API Keys",
        path: "/dashboard/apikeys",
      },
    ]);

    return {
      links,
      userStore,
    };
  },
});
</script>
