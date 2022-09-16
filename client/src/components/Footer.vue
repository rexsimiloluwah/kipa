<template>
  <footer>
    <div class="bg-primarygreen">
      <Container>
        <div>
          <div
            class="
              flex
              justify-between
              md:flex-row
              space-y-4
              md:space-y-0
              flex-col
            "
          >
            <div>
              <div class="nav__brand font-bold text-2xl text-white">
                kipa
                <span
                  class="bg-green-500 w-2 h-2 rounded-full inline-block"
                ></span>
              </div>
              <p>A secure {{ "&" }} serverless key/value store.</p>
              <ul class="inline-flex space-x-5 items-center my-2">
                <li
                  class="
                    cursor-pointer
                    hover:text-white
                    transition-all
                    duration-300
                  "
                  v-for="item in socialIcons"
                  :key="socialIcons.indexOf(item)"
                >
                  <font-awesome-icon
                    :icon="['fab', item.icon]"
                    class="text-2xl"
                  ></font-awesome-icon>
                </li>
              </ul>
            </div>
            <div class="grid grid-cols-2 gap-12">
              <div>
                <h1 class="font-bold text-[18px]">Navigation</h1>
                <ul>
                  <li
                    class="cursor-pointer hover:underline text-gray-700"
                    v-for="icon in footerNavigation"
                    :key="footerNavigation.indexOf(icon)"
                  >
                    <router-link :to="icon.url">{{ icon.title }}</router-link>
                  </li>
                </ul>
              </div>
              <div>
                <h1 class="font-bold text-[18px]">Support</h1>
                <ul>
                  <li
                    class="cursor-pointer hover:underline text-gray-700"
                    v-for="icon in supportNavigation"
                    :key="supportNavigation.indexOf(icon)"
                  >
                    <router-link :to="icon.url">{{ icon.title }}</router-link>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
        <div
          class="flex justify-between border-t-[1px] border-gray-500 pt-3 my-4"
        >
          <div>
            <p>©️2022 Kipa</p>
          </div>
          <span class="max-w-md">
            Built with <span aria-label="img">❤️️</span> by THEBLKDV
          </span>
        </div>
      </Container>
    </div>
  </footer>
</template>

<script lang="ts">
import { defineComponent, reactive } from "vue";
import { Container } from "./shared";
import { useCounterStore } from "../store/counter";

export default defineComponent({
  name: "Footer",
  components: { Container },
  setup() {
    const store = useCounterStore();
    const socialIcons = reactive([
      {
        title: "Mail",
        url: "",
        icon: "google",
      },
      {
        title: "Whatsapp",
        url: "",
        icon: "whatsapp",
      },
      {
        title: "GitHub",
        url: "",
        icon: "github",
      },
      {
        title: "Twitter",
        url: "",
        icon: "twitter",
      },
    ]);

    const footerNavigation = reactive([
      {
        title: "Create Account",
        url: "/signup",
      },
      {
        title: "About",
        url: "",
      },
      {
        title: "Your Buckets",
        url: "/buckets",
      },
      {
        title: "Contact",
        url: "/contact",
      },
    ]);

    const supportNavigation = reactive([
      {
        title: "Report an Issue/Bug",
        url: "",
      },
      {
        title: "Stay Updated!",
        url: "",
      },
      {
        title: "Privacy Policy",
        url: "",
      },
      {
        title: "Contribute",
        url: "",
      },
    ]);

    return { store, supportNavigation, socialIcons, footerNavigation };
  },
});
</script>
