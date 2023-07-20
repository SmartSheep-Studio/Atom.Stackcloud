<template>
  <n-spin :show="reverting">
    <div>
      <n-card size="small" class="px-2 rounded-0">
        <div class="flex justify-between items-center">
          <div class="flex items-center">
            <div>
              <div class="flex items-center gap-2" v-if="isSaved">
                <div class="text-green-600">Saved</div>
                <n-badge type="success" dot />
              </div>
              <div class="flex items-center gap-2 cursor-pointer" @click="save" v-else>
                <div class="text-orange-400">Unsaved</div>
                <n-badge type="warning" dot processing />
              </div>
            </div>
          </div>
          <div>{{ func.name }}</div>
          <n-space size="small">
            <n-button size="small" type="info" tertiary circle @click="updating = true">
              <template #icon>
                <n-icon :component="SettingsRound" />
              </template>
            </n-button>
            <n-button size="small" type="primary" tertiary circle @click="calling = true">
              <template #icon>
                <n-icon :component="PlayArrowRound" />
              </template>
            </n-button>
          </n-space>
        </div>
      </n-card>
      <div>
        <vue-monaco-editor
          :height="isUnderShadow ? 'calc(100vh - 54px - 46px - 48px)' : 'calc(100vh - 54px - 46px)'"
          v-model:value="code"
          @update:value="isSaved = false"
          @keyup.ctrl.s.prevent="save"
          language="javascript"
          theme="vs"
        />
      </div>

      <n-modal v-model:show="updating">
        <n-card size="huge" class="w-[800px]">
          <div class="text-lg mb-5">Update a exists collection</div>
          <n-card>
            <update-function />
          </n-card>
          <div class="text-lg mt-7 mb-5">Dangerous Zone</div>
          <destroy-function
            :data="func"
            @done="$router.push({ name: 'console.apps', params: { app: $route.params.app } }).then(() => reload())"
          />
        </n-card>
      </n-modal>

      <n-modal v-model:show="calling">
        <n-card title="Call a function" size="huge" class="w-[800px]">
          <call-function :data="func" />
        </n-card>
      </n-modal>
    </div>
  </n-spin>
</template>

<script lang="ts" setup>
import CallFunction from "@/views/actions/call-function.vue"
import UpdateFunction from "@/views/actions/update-function.vue"
import DestroyFunction from "@/views/actions/destroy-function.vue"
import { computed, onMounted, ref } from "vue"
import { useMessage } from "naive-ui"
import { useRoute, useRouter } from "vue-router"
import { http } from "@/utils/http"
import { PlayArrowRound, SettingsRound } from "@vicons/material"

const $route = useRoute()
const $router = useRouter()
const $message = useMessage()

const reverting = ref(true)
const updating = ref(false)
const calling = ref(false)
const isSaved = ref(true)

const func = ref<any>({})
const code = ref<string>("")

async function fetch() {
  try {
    reverting.value = true
    func.value = (await http.get(`/api/apps/${$route.params.app}/functions/${$route.params.function}`)).data
    code.value = func.value.script
  } catch (e: any) {
    $message.error(`Something went wrong... ${e}`)
  } finally {
    reverting.value = false
  }
}

async function save() {
  if (isSaved.value) {
    return
  }
  try {
    reverting.value = true
    func.value.script = code.value
    await http.put(`/api/apps/${$route.params.app}/functions/${$route.params.function}`, func.value)
    isSaved.value = true
  } catch (e: any) {
    $message.error(`Something went wrong... ${e}`)
  } finally {
    reverting.value = false
  }
}

function reload() {
  $router.push({ name: "console.apps.functions", params: { app: $route.params.app, function: func.value.slug } })
  window.location.reload()
}

onMounted(() => {
  fetch()
  document.addEventListener(
    "keyup",
    (e: KeyboardEvent) => {
      if (e.key === "s" && (navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey)) {
        e.preventDefault()
        save()
      }
    },
    { passive: false }
  )
})

// Use for dynamic calculate height
const isUnderShadow = computed(() => {
  return (window as any).__POWERED_BY_WUJIE__ != null
})
</script>
