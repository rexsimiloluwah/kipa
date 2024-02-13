<template>
  <table class="items__table">
    <thead>
      <tr>
        <th>Key</th>
        <th>Type</th>
        <th>Data</th>
        <th>Expires At</th>
        <th>Created At</th>
        <th><font-awesome-icon icon="pencil"></font-awesome-icon></th>
        <th><font-awesome-icon icon="trash-can"></font-awesome-icon></th>
      </tr>
    </thead>
    <tbody v-if="bucketItems">
      <tr
        v-for="item in bucketItems"
        :key="item.id"
      >
        <td>{{ item.key }}</td>
        <td>{{ item.type }}</td>
        <td
          class="underline underline-offset-4 cursor-pointer hover:font-semibold"
          @click="$emit('open-bucket-item', item)"
        >
          View
        </td>
        <td>
          {{
            item.ttl === 0
              ? "no expiry"
              : isExpired(new Date(item.created_at), item.ttl)
                ? "expired"
                : getExpiredAt(new Date(item.created_at), item.ttl)
          }}
        </td>
        <td>
          {{ new Date(item.created_at).toISOString() }}
        </td>
        <td>
          <button name="edit_item" role="edit_item" aria-label="edit_item" @click="$emit('update-bucket-item', item)">
            <font-awesome-icon icon="pencil" class="text-blue-500"></font-awesome-icon>
          </button>
        </td>
        <td>
          <button name="delete_item" role="delete_item" aria-label="delete_item" @click="$emit('delete-bucket-item', item)">
            <font-awesome-icon icon="trash-can" class="text-red-500"></font-awesome-icon>
          </button>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { BucketItem } from "../../../common/types/bucket_item";

export default defineComponent({
  name: "BucketItemsTable",
  props: {
    bucketItems: {
      type: Object as PropType<BucketItem[]>,
      required: true,
    },
  },
  emits: ["open-bucket-item","update-bucket-item","delete-bucket-item"],
  setup() {
    // computes the expired status of the bucket item
    const isExpired = (created_at: Date, ttl: number) => {
      if (ttl == 0) {
        return false;
      }
      const now = new Date().getTime();
      return now > new Date(created_at.getTime() + ttl * 1000).getTime();
    };

    // computes the expired at datetime value as an ISO string
    const getExpiredAt = (created_at: Date, ttl: number) => {
      return new Date(created_at.getTime() + ttl * 1000).toISOString();
    };

    return {
      isExpired,
      getExpiredAt,
    };
  },
});
</script>
