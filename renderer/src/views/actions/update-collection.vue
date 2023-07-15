<template>
  <div class="container">
    <div class="pt-12 pb-4 px-10">
      <div class="text-2xl font-bold">Update a exists collection</div>
      <div class="text-lg">A place to store a lot of serializable data.</div>
    </div>

    <div class="px-10 pt-4">
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
  </div>
</template>

<script lang="ts" setup>
import { http } from "@/utils/http"
import { useMessage, type FormRules, type FormInst } from "naive-ui"
import { onMounted, reactive, ref } from "vue"
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
}

const payload = ref({
  slug: "",
  name: "",
  description: "",
  tags: [],
  is_published: false,
})

async function fetch() {
  try {
    payload.value = (await http.get(`/api/apps/${$route.params.app}/records/${$route.params.collection}`)).data
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

      await http.put(`/api/apps/${$route.params.app}/records/${$route.params.collection}`, payload.value)

      $message.success("Successfully updated a collection.")
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
