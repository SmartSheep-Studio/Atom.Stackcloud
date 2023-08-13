<template>
  <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="create">
    <n-form-item label="Payload" path="payload">
      <vue-monaco-editor v-model:value="payload.payload" language="json" height="400px" theme="vs-dark" />
    </n-form-item>

    <n-space size="small">
      <n-button type="primary" attr-type="submit" :loading="submitting">Submit</n-button>
    </n-space>
  </n-form>
</template>

<script lang="ts" setup>
import { http } from "@/utils/http"
import { useMessage, type FormRules, type FormInst } from "naive-ui"
import { reactive, ref } from "vue"
import { useRoute } from "vue-router"

const $route = useRoute()
const $message = useMessage()

const submitting = ref(false)

const emits = defineEmits(["refresh", "close"])

const form = ref<FormInst | null>(null)
const rules: FormRules = {
  payload: {
    required: true,
    validator: (_, v) => {
      try {
        JSON.parse(v)
        return true
      } catch {
        return false
      }
    },
    message: "Only accepts valid json",
    trigger: ["blur", "input"],
  },
}

const payload = reactive({
  payload: "{}",
})

function create() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return
    }

    try {
      submitting.value = true

      await http.post(`/api/apps/${$route.params.app}/records/${$route.params.collection}/data`, {
        payload: JSON.parse(payload.payload),
      })

      emits("close")
      emits("refresh")
      $message.success("Successfully created a new record.")
    } catch (e: any) {
      $message.error(`Something went wrong... ${e}`)
    } finally {
      submitting.value = false
    }
  })
}
</script>
