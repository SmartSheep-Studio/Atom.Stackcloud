<template>
  <n-form ref="form" :rules="rules" :model="payload" @submit.prevent="call">
    <n-form-item label="Payload" path="payload">
      <n-card content-style="padding: 0">
        <n-tabs type="line" animated :tabs-padding="20">
          <n-tab-pane name="Body" style="padding: 0">
            <vue-monaco-editor v-model:value="payload.payload" language="json" height="400px" theme="vs" />
          </n-tab-pane>
          <n-tab-pane name="Headers" style="padding: 0">
            <div class="px-5 py-5">
              <n-dynamic-input
                preset="pair"
                v-model:value="payload.headers"
                key-placeholder="Key"
                value-placeholder="Value"
              />
            </div>
          </n-tab-pane>
        </n-tabs>
      </n-card>
    </n-form-item>

    <n-form-item label="Response">
      <n-spin :show="submitting" class="w-full">
        <n-card size="small">
          <n-card size="small" embedded v-if="!finished">
            <n-empty description="Request isn't sent/completed" />
          </n-card>
  
          <n-space vertical size="small" v-else>
            <n-card size="small" title="Header" embedded>
              <div class="flex justify-between">
                <div>status</div>
                <div>{{ response.status }}</div>
              </div>
              <div v-for="item in response.headers" class="flex justify-between">
                <div>{{ item[0] }}</div>
                <div>{{ item[1] }}</div>
              </div>
            </n-card>
            <n-card size="small" title="Body" embedded>
              <div v-if="response.data == null || response.data.length <= 0"><n-code code="<null>" /></div>
              <n-code :hljs="hljs" :code="response.data" language="json" word-wrap v-else />
            </n-card>
          </n-space>
        </n-card>
      </n-spin>
    </n-form-item>

    <n-space size="small">
      <n-button type="primary" attr-type="submit" :loading="submitting">Call</n-button>
      <n-button @click="finished = false">Reset</n-button>
    </n-space>
  </n-form>
</template>

<script lang="ts" setup>
import hljs from "highlight.js/lib/core"
import json from "highlight.js/lib/languages/json"
import { http } from "@/utils/http"
import { useMessage, type FormInst, type FormRules } from "naive-ui"
import { reactive, ref } from "vue"
import { useRoute } from "vue-router"

hljs.registerLanguage("json", json)

const $route = useRoute()
const $message = useMessage()

const props = defineProps<{ data: any }>()

const submitting = ref(false)
const finished = ref(false)

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

const payload = reactive<{ headers: { key: string; value: string }[]; payload: string }>({
  headers: [],
  payload: "{}",
})

const response = reactive<{ status: number; headers: any[]; data: any }>({
  status: 200,
  headers: [],
  data: null,
})

function call() {
  form.value?.validate(async (errors) => {
    if (errors) {
      return
    }

    try {
      submitting.value = true

      const headers: { [key: string]: string } = {}
      payload.headers.forEach((v) => {
        headers[v.key] = v.value
      })

      const res = await http.post(
        `/api/apps/${$route.params.app}/functions/${props.data.slug}/call`,
        JSON.parse(payload.payload),
        {
          headers: headers,
          validateStatus: () => true,
        }
      )

      response.status = res.status
      response.headers = Object.entries(res.headers)
      response.data = JSON.stringify(res.data, null, 4)
      finished.value = true

      $message.success("Successfully called the function.")
    } catch (e: any) {
      $message.error(`Something went wrong... ${e}`)
    } finally {
      submitting.value = false
    }
  })
}
</script>
