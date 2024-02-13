<template>
  <form
    class="flex flex-col space-y-3 z-[1000]"
    @submit.prevent="handleUpdateBucket"
  >
    <div>
      <div class="relative space-y-1">
        <label for="name">Bucket Name<span class="text-red-600">*</span></label>
        <input
          id="name"
          v-model="formState.name"
          type="text"
          name="name"
          placeholder="Enter Name i.e. 'test-pipeline-bucket'"
          :class="`${v$.name.$errors.length && 'input--error'}`"
          @blur="handleBlur('name')"
        >
      </div>
      <div
        v-for="error of v$.name.$errors"
        :key="error.$uid"
        class="input-errors"
      >
        <p class="text-red-600">
          {{ parseErrorMessage(String(error.$message), "Bucket Name") }}
        </p>
      </div>
    </div>

    <div>
      <div class="relative space-y-1">
        <label for="description">Bucket Description</label>
        <input
          id="description"
          v-model="formState.description"
          type="text"
          name="description"
          placeholder="Enter Description"
          :class="`${v$.description.$errors.length && 'input--error'}`"
          @blur="handleBlur('description')"
        >
      </div>
      <div
        v-for="error of v$.description.$errors"
        :key="error.$uid"
        class="input-errors"
      >
        <p class="text-red-600">
          {{ parseErrorMessage(String(error.$message), "Description") }}
        </p>
      </div>
    </div>

    <div class="space-y-1">
      <h1>
        Select Bucket Permissions ({{
          selectedBucketPermissions.filter((p) => p).length
        }})
      </h1>
      <div class="flex flex-wrap gap-1">
        <button
          v-for="permission in bucketPermissions.map((p, i) => ({
            name: p,
            id: i,
          }))"
          :key="permission.id"
          type="button"
          :class="`${
            selectedBucketPermissions[permission.id]
              ? 'bg-primarygreen'
              : 'border-2 border-primarygreen'
          } p-2 rounded-lg`"
          @click="
            selectedBucketPermissions[permission.id] =
              !selectedBucketPermissions[permission.id]
          "
        >
          {{ permission.name }}
          <span
            v-if="selectedBucketPermissions[permission.id]"
            class="text-md font-semibold"
          >x</span>
        </button>
      </div>
    </div>

    <custom-button
      title="Edit Bucket"
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
import { BUCKET_PERMISSIONS } from "../../../common/constants";
import { useBucketStore } from "../../../store/bucket";
import { BucketDetails } from "../../../common/types/bucket";

export default defineComponent({
  props: {
    bucket: {
      type: Object as PropType<BucketDetails>,
      required: true,
    },
  },
  emits: ["closeModal"],
  setup(props, { emit }) {
    const bucketStore = useBucketStore();
    const isLoading = ref<boolean>(false);
    const bucketPermissions = reactive<Array<string>>(BUCKET_PERMISSIONS);
    const selectedBucketPermissions = ref<Array<boolean>>(
      bucketPermissions.map((p) => (props.bucket.permissions || []).includes(p))
    );

    const formState = reactive({
      name: props.bucket.name,
      description: props.bucket.description,
    });
    const rules = {
      name: { required },
      description: {
        maxLength: maxLength(200),
      },
    };

    const v$ = useVuelidate(rules, formState);

    const handleBlur = (key: "name" | "description") => {
      // @ts-ignore
      v$.value[key].$dirty = true;
    };

    // For updating the bucket
    const handleUpdateBucket = async () => {
      isLoading.value = true;
      const updateBucketData = {
        ...formState,
        permissions: bucketPermissions.filter(
          (_, id) => selectedBucketPermissions.value[id]
        ),
      };
      await bucketStore
        .updateBucket(props.bucket.uid, updateBucketData)
        .then(() => {
          isLoading.value = false;
          emit("closeModal");
        });
    };

    return {
      bucketStore,
      formState,
      v$,
      isLoading,
      handleBlur,
      handleUpdateBucket,
      parseErrorMessage,
      bucketPermissions,
      selectedBucketPermissions,
    };
  },
});
</script>
