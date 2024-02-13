<template>
    <form
      class="flex flex-col space-y-3 z-[1000]"
      @submit.prevent="handleUpdateBucketItem"
    >
      <div>
        <div class="relative space-y-1">
          <label for="key">Key</label>
          <input
            id="key"
            v-model="formState.key"
            type="text"
            name="key"
            placeholder="Enter Key Name i.e. test-key"
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
          <label for="data">Data</label>
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
        title="Update Bucket Item"
        type="submit"
        :disabled="v$.$invalid"
        :loading="isLoading"
      />
    </form>
  </template>
  
  <script lang="ts">
  import { defineComponent, ref, reactive, PropType } from "vue";
  import useVuelidate from "@vuelidate/core";
  import { maxLength } from "@vuelidate/validators";
  import { parseErrorMessage } from "../../../common/utils/form";
  import BucketItemService from "../../../services/bucket_item"
  import { useBucketItemStore } from "../../../store/bucket_item";
  
  export default defineComponent({
    name: "UpdateBucketItemForm",
    props: {
      bucketUID: {
        type: String as PropType<string>,
        required: true,
      }
    },
    emits: ["closeModal"],
    setup(props, { emit }) {
      const isLoading = ref<boolean>(false);
      const bucketItemStore = useBucketItemStore()
      
      const formState = reactive({
        key: bucketItemStore.activeBucketItem?.key,
        data: bucketItemStore.activeBucketItem?.data,
        ttl: bucketItemStore.activeBucketItem?.ttl,
      });
      const rules = {
        key: { maxLength: maxLength(200) },
        data: {},
        ttl: {min:0}
      };

    const v$ = useVuelidate(rules, formState);
  
      const handleBlur = (key: "key"|"data"|"ttl") => {
        // @ts-ignore
        v$.value[key].$dirty = true;
      };
  
      const handleUpdateBucketItem = async () => {
        delete formState["key"]
        const bucketItemData = {
          ...formState,
        };
        isLoading.value = true;
        await bucketItemStore.updateBucketItem(props.bucketUID, bucketItemStore.activeBucketItem?.key as string,bucketItemData).then(() => {
          isLoading.value = false
          emit("closeModal")
        })
      };
  
      return {
        formState,
        v$,
        isLoading,
        handleBlur,
        handleUpdateBucketItem,
        parseErrorMessage,
      };
    },
  });
  </script>
  