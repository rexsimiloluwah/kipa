<template>
  <div class="w-full py-4">
    <h1 class="text-2xl font-bold mb-3">Update Profile</h1>
    <form class="space-y-3" @submit.prevent="handleUpdateAccount">
      <FormRow>
        <div>
          <div class="relative">
            <input
              type="text"
              name="firstname"
              id="firstname"
              v-on:blur="handleBlur('firstname')"
              v-model="formState.firstname"
              placeholder="Enter Firstname"
              :class="`${v$.firstname.$errors.length && 'input--error'}`"
            />
          </div>
          <div
            class="input-errors"
            v-for="error of v$.firstname.$errors"
            :key="error.$uid"
          >
            <p class="text-red-600">
              {{ parseErrorMessage(String(error.$message), "Firstname") }}
            </p>
          </div>
        </div>

        <div>
          <div class="relative">
            <input
              type="text"
              name="lastname"
              id="lastname"
              v-on:blur="handleBlur('lastname')"
              v-model="formState.lastname"
              placeholder="Enter Lastname"
              :class="`${v$.lastname.$errors.length && 'input--error'}`"
            />
          </div>
          <div
            class="input-errors"
            v-for="error of v$.lastname.$errors"
            :key="error.$uid"
          >
            <p class="text-red-600">
              {{ parseErrorMessage(String(error.$message), "Lastname") }}
            </p>
          </div>
        </div>
      </FormRow>
      <div>
        <div class="relative">
          <input
            type="text"
            name="email"
            id="email"
            v-on:blur="handleBlur('email')"
            v-model="formState.email"
            placeholder="Enter Email"
            :class="`${v$.email.$errors.length && 'input--error'}`"
          />
        </div>
        <div
          class="input-errors"
          v-for="error of v$.email.$errors"
          :key="error.$uid"
        >
          <p class="text-red-600">
            {{ parseErrorMessage(String(error.$message), "Email") }}
          </p>
        </div>
      </div>

      <Button
        title="Update Profile"
        type="submit"
        :disabled="v$.$invalid"
        :loading="isLoading"
      />
    </form>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from "vue";
import { required, email, minLength } from "@vuelidate/validators";
import { useVuelidate } from "@vuelidate/core";
import Button from "../../shared/Button.vue";
import FormRow from "../../shared/FormRow.vue";
import FormInput from "../../shared/FormInput.vue";
import { parseErrorMessage } from "../../../common/utils/form";
import { useUserStore } from "../../../store/user";

export default defineComponent({
  name: "UpdateAccountForm",
  components: { Button, FormInput, FormRow },
  setup() {
    const userStore = useUserStore();
    const isLoading = ref<boolean>(false);
    const formState = reactive({
      firstname: userStore.user?.firstname as string,
      lastname: userStore.user?.lastname as string,
      email: userStore.user?.email as string,
    });
    const rules = {
      firstname: { required, minLength: minLength(2) },
      lastname: { required, minLength: minLength(2) },
      email: { required, email },
    };

    const v$ = useVuelidate(rules, formState);

    const handleBlur = (key: "email" | "firstname" | "lastname") => {
      // @ts-ignore
      v$.value[key].$dirty = true;
    };

    const handleUpdateAccount = async () => {
      isLoading.value = true;
      const data = { ...formState };
      await userStore.updateUser(data).then(() => {
        isLoading.value = false;
      });
    };

    return {
      formState,
      v$,
      isLoading,
      handleBlur,
      parseErrorMessage,
      handleUpdateAccount,
    };
  },
});
</script>