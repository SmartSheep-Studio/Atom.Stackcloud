<template>
  <div>
    <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="update" class="max-w-[800px]">
      <n-form-item label="Slug" path="slug">
        <n-input
          placeholder="Use for the link to your collection. Only accepts url safe characters."
          v-model:value="payload.slug"
        />
      </n-form-item>
      <n-form-item label="Name" path="name">
        <n-input placeholder="Used to hint the developer what the collection is for." v-model:value="payload.name" />
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
</template>

<script lang="ts" setup>
import { parseRedirect } from "@/utils/callback"
import { http } from "@/utils/http"
import { useMessage, type FormRules, type FormInst } from "naive-ui"
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
  script: {
    required: true,
    message: "Need least one character",
    trigger: ["blur", "input"],
  },
}

const payload = ref({
  slug: "",
  name: "",
  description: "",
  tags: [],
  script:
    "// Your code goes there.\n// Support ES5 and most ES6 syntax\n// Learn more from our official documentation!",
})

async function fetch() {
  try {
    payload.value = (await http.get(`/api/apps/${$route.params.app}/functions/${$route.params.function}`)).data
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

      await http.put(`/api/apps/${$route.params.app}/functions/${$route.params.function}`, payload.value)

      $message.success("Successfully updated a function.")
      window.location.reload()
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
