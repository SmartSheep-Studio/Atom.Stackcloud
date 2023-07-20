<template>
  <div class="px-8 py-4">
    <n-spin :show="reverting">
      <n-grid item-responsive responsive="screen" x-gap="8" y-gap="8" v-if="!reverting">
        <n-gi span="24 m:14 l:16">
          <n-card title="Update App">
            <update-app :data="app" />
          </n-card>
        </n-gi>
        <n-gi span="24 m:10 l:8">
          <n-card title="Dangerous Zone">
            <destroy-app
              :data="app"
              @done="$router.push({ name: 'console' })"
            />
          </n-card>
        </n-gi>
      </n-grid>
    </n-spin>
  </div>
</template>

<script lang="ts" setup>
import UpdateApp from "@/views/actions/update-app.vue"
import DestroyApp from "@/views/actions/destroy-app.vue"
import { onMounted, ref } from "vue"
import { http } from "@/utils/http"
import { useMessage } from "naive-ui"
import { useRoute } from "vue-router"

const $message = useMessage()
const $route = useRoute()

const app = ref<any>({})

const reverting = ref(true)

async function fetch() {
  try {
    reverting.value = true
    app.value = (await http.get(`/api/apps/${$route.params.app}`)).data
  } catch (e: any) {
    $message.error(`Something went wrong... ${e}`)
  } finally {
    reverting.value = false
  }
}

onMounted(() => {
  fetch()
})
</script>
