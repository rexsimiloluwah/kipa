<template>
  <main>
    <transition name="fade">
      <div
        key="modal"
        @click.self="$emit('close')"
        class="
          modal
          fixed
          top-0
          left-0
          w-full
          grid
          overflow-x-hidden overflow-y-hidden
          h-screen
          z-[99]
          bg-black/[.80]
          place-items-center
        "
      >
        <transition name="drop-in" mode="out-in">
          <div
            class="modal__body rounded-lg bg-white w-[40%] overflow-hidden"
            key="modal__body"
          >
            <div
              :class="`
                modal__header
                flex
                items-center
                justify-between
                relative
                ${title !== '' ? ' bg-green-50 p-3' : 'px-3 pt-3'}
              `"
            >
              <span class="text-xl font-bold">{{ title }}</span>
              <button
                type="button"
                title=""
                @click="$emit('close')"
                class="
                  bg-primarygreen
                  rounded-full
                  w-8
                  h-8
                  border-[1px] border-black
                  text-white
                  opacity-80
                  hover:opacity-100
                "
              >
                <font-awesome-icon icon="x" />
              </button>
            </div>
            <div class="modal__content px-6 pb-4 max-h-[25rem] overflow-y-auto">
              <slot />
            </div>
          </div>
        </transition>
      </div>
    </transition>
  </main>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Button from "./Button.vue";

export default defineComponent({
  name: "Modal",
  components: {
    Button,
  },
  props: {
    title: {
      type: String,
      default: "",
      required: false,
    },
  },
  setup() {},
});
</script>

<style scoped>
</style>