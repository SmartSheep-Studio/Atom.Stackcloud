<template>
  <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="update">
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
import { reactive, ref, onMounted } from "vue"
import { useRoute } from "vue-router"

const $route = useRoute()
const $message = useMessage()

const submitting = ref(false)

const props = defineProps<{ data: any }>()
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

function update() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return
    }

    try {
      submitting.value = true

      await http.put(`/api/apps/${$route.params.app}/records/${$route.params.collection}/data/${props.data.id}`, {
        payload: JSON.parse(payload.payload),
      })

      emits("close")
      emits("refresh")
      $message.success("Successfully updated a record.")
    } catch (e: any) {
      $message.error(`Something went wrong... ${e}`)
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  payload.payload = JSON.stringify(props.data.payload, null, 4)
})
</script>
