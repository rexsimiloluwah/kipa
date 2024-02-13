<template>
    <form
      class="flex flex-col space-y-3 z-[1000]"
      @submit.prevent="handleCreateBucketItem"
    >
      <div>
        <div class="relative space-y-1">
          <label for="key">Key<span class="text-red-600">*</span></label>
          <input
            id="key"
            v-model="formState.key"
            type="text"
            name="key"
            placeholder="Enter key name i.e. test-key"
            :class="`p-3 rounded-md border-primarygreen
            border-2 w-full ${v$.key.$errors.length && 'input--error'}`"
            @blur="handleBlur('key')"
          >
        </div>
        <div
          v-for="error of v$.key.$errors"
          :key="error.$uid"
          class="input-errors"
        >
          <p class="text-red-600">
            {{ parseErrorMessage(String(error.$message), "Key Name") }}
          </p>
        </div>
      </div>
  
      <div>
        <div class="relative space-y-1">
          <label for="data">Data<span class="text-red-600">*</span></label>
          <input
            id="data"
            v-model="formState.data"
            type="text"
            name="data"
            placeholder="Enter Data"
            :class="`p-3 rounded-md border-primarygreen
            border-2 w-full ${v$.data.$errors.length && 'input--error'}`"
            @blur="handleBlur('data')"
          >
        </div>
        <div
          v-for="error of v$.data.$errors"
          :key="error.$uid"
          class="input-errors"
        >
          <p class="text-red-600">
            {{ parseErrorMessage(String(error.$message), "Data") }}
          </p>
        </div>
      </div>

      <div>
        <div class="relative space-y-1">
          <label for="ttl">TTL</label>
          <input
            id="ttl"
            v-model="formState.ttl"
            type="number"
            name="ttl"
            placeholder="Enter TTL (time to live in seconds)"
            :class="`p-3 rounded-md border-primarygreen
            border-2 w-full ${v$.ttl.$errors.length && 'input--error'}`"
            @blur="handleBlur('ttl')"
          >
        </div>
        <div
          v-for="error of v$.ttl.$errors"
          :key="error.$uid"
          class="input-errors"
        >
          <p class="text-red-600">
            {{ parseErrorMessage(String(error.$message), "TTL") }}
          </p>
        </div>
      </div>
  
      <custom-button
        title="Create Bucket"
        type="submit"
        :disabled="v$.$invalid"
        :loading="isLoading"
      />
    </form>
  </template>
  
  <script lang="ts">
  import { defineComponent, ref, reactive, PropType } from "vue";
  import useVuelidate from "@vuelidate/core";
  import { required, maxLength } from "@vuelidate/validators";
  import { parseErrorMessage } from "../../../common/utils/form";
  import { useBucketItemStore } from "../../../store/bucket_item";
  
  export default defineComponent({
    name: "CreateBucketItemForm",
    props: {
      bucketUID: {
        type: String as PropType<string>,
        required: true,
      }
    },
    emits: ["closeModal"],
    setup(props, { emit }) {
      const isLoading = ref<boolean>(false);
      const bucketItemStore = useBucketItemStore();
      
      const formState = reactive({
        key: "",
        data: "",
        ttl: 0,
      });
      const rules = {
        key: { required, maxLength: maxLength(200) },
        data: {required},
        ttl: {min:0}
      };

    const v$ = useVuelidate(rules, formState);
  
      const handleBlur = (key: "key"|"data"|"ttl") => {
        // @ts-ignore
        v$.value[key].$dirty = true;
      };
  
      const handleCreateBucketItem = async () => {
        const bucketItemData = {
          ...formState,
        };
        isLoading.value = true;
        await bucketItemStore.createBucketItem(props.bucketUID, bucketItemData).then(async () => {
          isLoading.value = false
          emit("closeModal")
        })
      };
  
      return {
        formState,
        v$,
        isLoading,
        handleBlur,
        handleCreateBucketItem,
        parseErrorMessage,
      };
    },
  });
  </script>
  