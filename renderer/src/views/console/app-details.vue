<template>
  <n-spin :show="reverting">
    <div class="container h-full">
      <div class="pt-8 pb-4 lg:px-10">
        <n-page-header :title="app.name" :subtitle="app.description" @back="$router.push({ name: 'console' })">
          <template #header>
            <n-breadcrumb>
              <n-breadcrumb-item @click="$router.push({ name: 'console' })">Console</n-breadcrumb-item>
              <n-breadcrumb-item>{{ app.name }}</n-breadcrumb-item>
            </n-breadcrumb>
          </template>
        </n-page-header>
      </div>

      <n-grid v-if="!reverting" item-responsive responsive="screen" x-gap="8" y-gap="8" class="lg:px-10 pt-4">
        <n-gi span="24 m:14">
          <n-card title="Update App">
            <update-app :data="app" @refresh="fetch()" />
          </n-card>
        </n-gi>
        <n-gi span="24 m:6 l:8">
          <n-card title="Dangerous Zone">
            <publish-app :data="app" @refresh="fetch()" />
            <destroy-app class="mt-2" :data="app" @done="$router.push({ name: 'console' })" />
          </n-card>
        </n-gi>
      </n-grid>
    </div>
  </n-spin>
</template>

<script lang="ts" setup>
import UpdateApp from "@/views/console/actions/update-app.vue"
import DestroyApp from "@/views/console/actions/destroy-app.vue"
import PublishApp from "@/views/console/actions/publish-app.vue"
import { useMessage } from "naive-ui"
import { useRoute } from "vue-router"
import { onMounted, ref } from "vue"
import { http } from "@/utils/http"

const $route = useRoute()
const $message = useMessage()

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
