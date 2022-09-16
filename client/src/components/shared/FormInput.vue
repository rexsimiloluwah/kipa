<template>
  <div>
    <div class="relative">
      <input
        type="text"
        :name="name"
        :id="name"
        :value="value"
        :placeholder="placeholder"
        class="p-3 rounded-md border-primarygreen border-2 w-full"
      />
    </div>
    <div class="input-errors" v-for="error of errors" :key="error.$uid">
      <p class="text-red-600">
        {{
          parseErrorMessage(
            String(error.$message),
            name.charAt(0).toUpperCase() + name.slice(1)
          )
        }}
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import { ErrorObject } from "@vuelidate/core";
import { defineComponent, PropType } from "vue";

export default defineComponent({
  name: "FormInput",
  props: {
    name: {
      type: String,
      required: true,
    },
    errors: {
      type: Array as PropType<ErrorObject[]>,
      required: false,
    },
    type: {
      type: String as PropType<"text" | "password">,
      required: false,
      default: "text",
    },
    value: {
      type: String as PropType<string>,
      required: true,
    },
    placeholder: {
      type: String,
      required: true,
    },
  },
  setup() {
    // utility function for parsing the error message
    const parseErrorMessage = (message: string, key: string) => {
      const firstWord = message.split(" ")[0];
      const firstTwoWords = message.split(" ").slice(0, 2).join(" ");
      if (firstWord.toLowerCase() === "value") {
        return message.replace("Value", key);
      }
      if (firstTwoWords.toLowerCase() === "this field") {
        return message.replace("This field", key);
      }
      return message;
    };

    return {
      parseErrorMessage,
    };
  },
});
</script>
