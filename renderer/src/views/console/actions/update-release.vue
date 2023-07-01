<template>
  <div class="container">
    <div class="pt-12 pb-4 px-10">
      <div class="text-2xl font-bold">Update a exists release</div>
      <div class="text-lg">Don't forgot follow the community guidelines!</div>
    </div>

    <div class="px-10 pt-4">
      <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="update" class="max-w-[800px]">
        <n-form-item label="Slug" path="slug">
          <n-input
            placeholder="Use for the link to your release. Only accepts url safe characters."
            v-model:value="payload.slug"
          />
        </n-form-item>
        <n-form-item label="Type" path="type">
          <n-select :options="types" v-model:value="payload.type" />
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
            placeholder="Use for describe your release main content. Accepts anything you want."
            v-model:value="payload.description"
          />
        </n-form-item>
        <n-form-item label="Details" path="details">
          <v-md-editor v-model="payload.details" height="400px" />
        </n-form-item>
        <n-form-item label="Is Published" path="is_published">
          <n-switch v-model:value="payload.is_published" />
        </n-form-item>

        <n-space size="small">
          <n-button type="primary" attr-type="submit" :loading="submitting">Submit</n-button>
          <n-button @click="$router.back()">Cancel</n-button>
        </n-space>
      </n-form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { parseRedirect } from "@/utils/callback"
import { http } from "@/utils/http"
import { useMessage, type FormRules, type FormInst, useDialog } from "naive-ui"
import { onMounted, ref } from "vue"
import { useRoute, useRouter } from "vue-router"

const $route = useRoute()
const $router = useRouter()
const $message = useMessage()

const submitting = ref(false)

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
  details: {
    required: true,
    validator: (_, v) => v.length >= 6,
    message: "Need least six characters",
    trigger: ["blur", "input"],
  },
}

const types = [
  { label: "Minor Update", value: "minor-update" },
  { label: "Major Update", value: "major-update" },
  { label: "Hotfix Update", value: "hotfix-update" },
]

const payload = ref({
  slug: "",
  type: "minor-update",
  name: "",
  description: "",
  details: "",
  tags: [],
  is_published: false,
})

async function fetch() {
  try {
    const res = await http.get(`/api/apps/${$route.params.app}/releases/${$route.params.release}`)
    payload.value = res.data
    payload.value.type = res.data.post.type
    payload.value.tags = res.data.post.tags
    payload.value.details = res.data.post.content
  } catch (e: any) {
    $message.error(`Something went wrong... ${e}`)
  }
}

function update() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return
    }

    try {
      submitting.value = true

      await http.put(`/api/apps/${$route.params.app}/releases/${$route.params.release}`, payload.value)

      $message.success("Successfully updated a release")
      await $router.push(
        await parseRedirect($route.query, { name: "console.apps", params: { app: $route.params.app } })
      )
    } catch (e: any) {
      $message.error(`Something went wrong... ${e}`)
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  fetch()
})
</script>
