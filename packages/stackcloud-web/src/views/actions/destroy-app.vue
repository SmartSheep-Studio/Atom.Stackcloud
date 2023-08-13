<template>
  <n-card embedded>
    <div class="text-lg">Destroy App</div>
    <div class="mt-1">
      This operation will destroy all resources belong this application. And cannot be undo. Think twice.
    </div>

    <div class="mt-3">
      <n-button class="w-full" type="error" :loading="submitting" @click="destroy">
        <template #icon>
          <n-icon :component="DeleteRound" />
        </template>
        Destroy
      </n-button>
    </div>
  </n-card>
</template>

<script lang="ts" setup>
import { http } from "@/utils/http"
import { DeleteRound } from "@vicons/material"
import { useDialog, useMessage } from "naive-ui"
import { ref } from "vue"

const $dialog = useDialog()
const $message = useMessage()

const props = defineProps<{ data: any }>()
const emits = defineEmits(["done"])

const submitting = ref(false)

function destroy() {
  $dialog.warning({
    title: "Warning",
    content: "This operation cannot be undo. Are you confirm?",
    positiveText: "Yes",
    negativeText: "Not really",
    onPositiveClick: async () => {
      try {
        submitting.value = true

        await http.delete(`/api/apps/${props.data.slug}`)

        emits("done")
        $message.success("Successfully deleted the app. Auto redirecting...")
      } catch (e: any) {
        $message.error(`Something went wrong... ${e}`)
      } finally {
        submitting.value = false
      }
    },
  })
}
</script>
