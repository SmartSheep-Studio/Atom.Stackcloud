<template>
  <div class="container">
    <div class="pt-12 pb-4 px-10">
      <div class="text-2xl font-bold">Create a new app</div>
      <div class="text-lg">Don't forgot follow the community guidelines!</div>
    </div>

    <div class="px-10 pt-4">
      <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="create" class="max-w-[800px]">
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
import { reactive, ref } from "vue"
import { useRoute, useRouter } from "vue-router"

const $route = useRoute()
const $router = useRouter()
const $dialog = useDialog()
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
}

const payload = reactive({
  slug: "",
  name: "",
  description: "",
  tags: [],
  is_published: false,
})

function create() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return
    }

    try {
      submitting.value = true

      await http.post("/api/apps", payload)

      $dialog.success({
        title: "Successfully created an app",
        content: "Now back to console and start use our services!",
        positiveText: "OK",
        onPositiveClick: async () => {
          await $router.push(await parseRedirect($route.query))
        },
      })
    } catch (e: any) {
      $message.error(`Something went wrong... ${e}`)
    } finally {
      submitting.value = false
    }
  })
}
</script>
