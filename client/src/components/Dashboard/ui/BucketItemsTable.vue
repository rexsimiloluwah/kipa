<template>
  <table class="items__table">
    <thead>
      <tr>
        <th>Key</th>
        <th>Type</th>
        <th>Data</th>
        <th>Expires At</th>
        <th>Created At</th>
      </tr>
    </thead>
    <tbody v-if="bucket">
      <tr v-for="item in bucket.bucket_items" :key="item.id">
        <td>{{ item.key }}</td>
        <td>string</td>
        <td class="underline underline-offset-4">View</td>
        <td>
          {{ new Date(item.created_at).toISOString() }}
        </td>
        <td>
          {{ new Date(item.ttl).toISOString() }}
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { BucketDetails } from "../../../common/types/bucket";

export default defineComponent({
  name: "BucketItemsTable",
  props: {
    bucket: {
      type: Object as PropType<BucketDetails>,
      required: true,
    },
  },
  setup() {},
});
</script>
