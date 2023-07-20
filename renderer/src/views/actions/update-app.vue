<template>
  <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="update">
    <n-form-item label="Slug" path="slug">
      <n-input
        placeholder="Use for the link to your application. Only accepts url safe characters."
        v-model:value="payload.slug"
      />
    </n-form-item>
    <n-form-item label="Name" path="name">
      <n-input placeholder="Use for pointing out topics. Accepts anything you want." v-model:value="payload.name" />
    </n-form-item>
    <n-form-item label="Tags" path="tags">
      <n-dynamic-tags v-model:value="payload.tags" />
    </n-form-item>
    <n-form-item label="Description" path="description">
      <n-input
        type="textarea"
        placeholder="Use for describe main content. Accepts anything you want."
        v-model:value="payload.description"
      />
    </n-form-item>

    <n-space size="small">
      <n-button type="primary" attr-type="submit" :loading="submitting">Submit</n-button>
    </n-space>
  </n-form>
</template>

<script lang="ts" setup>
import { http } from "@/utils/http"
import { useMessage, type FormRules, type FormInst } from "naive-ui"
import { onMounted, ref } from "vue"

const $message = useMessage()

const submitting = ref(false)

const props = defineProps<{data: any}>()

const form = ref<FormInst | null>(null)
const rules: FormRules = {
  slug: {
    required: true,
    validator: (_, v) => new RegExp(/^[A-Za-z0-9-_]+$/).test(v),
    message: "Only accepts letters, underscore and numbers without space",
    trigger: ["blur", "input"],
  },
  name: {
    required: true,
    validator: (_, v) => v.length >= 4,
    message: "Need least four characters",
    trigger: ["blur", "input"],
  },
  description: {
    required: true,
    validator: (_, v) => v.length >= 6,
    message: "Need least six characters",
    trigger: ["blur", "input"],
  },
}

const payload = ref({
  slug: "",
  name: "",
  description: "",
  tags: [],
  is_published: false,
})

function update() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return
    }

    try {
      submitting.value = true

      await http.put(`/api/apps/${props.data.slug}`, payload.value)

      $message.success("Successfully updated the app")
      window.location.reload()
    } catch (e: any) {
      $message.error(`Something went wrong... ${e}`)
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  payload.value = props.data
})
</script>
