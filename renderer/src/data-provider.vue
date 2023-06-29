<template>
  <n-layout-content class="content-placeholder w-full h-max" v-if="reverting">
    <div class="w-full h-max flex justify-center items-center">
      <n-spin show>
        <template #description>Connecting...</template>
      </n-spin>
    </div>
  </n-layout-content>

  <n-layout-content class="content" v-else>
    <slot />
  </n-layout-content>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { usePrincipal } from "@/stores/principal";
import { useEndpoint } from "@/stores/endpoint";
import { useDialog } from "naive-ui";

const $dialog = useDialog();

const $endpoint = useEndpoint();
const $principal = usePrincipal();

const reverting = ref(true);

async function fetch() {
  reverting.value = true;

  try {
    await Promise.all([$endpoint.fetch(), $principal.fetch()]);
  } catch (e: any) {
    $dialog.error({
      closable: false,
      closeOnEsc: false,
      title: "An error occurred during the connection.",
      content: e.toString(),
      positiveText: "Retry",
      onPositiveClick: () => {
        fetch();
      },
    });
  } finally {
    reverting.value = false;
  }
}

onMounted(() => {
  fetch();
});
</script>
