<template>
  <n-spin :show="reverting">
    <div class="container h-full">
      <div class="pt-8 px-2 pb-6 md:px-4 lg:px-10">
        <n-page-header
          :title="app.name"
          :subtitle="app.description"
          @back="$router.push({ name: 'console.apps', params: { app: $route.params.app } })"
        >
          <template #header>
            <n-breadcrumb>
              <n-breadcrumb-item @click="$router.push({ name: 'console' })">Console</n-breadcrumb-item>
              <n-breadcrumb-item @click="$router.push({ name: 'console.apps', params: { app: $route.params.app } })">
                {{ app.name }}
              </n-breadcrumb-item>
              <n-breadcrumb-item>{{ collection.name }}</n-breadcrumb-item>
            </n-breadcrumb>
          </template>
        </n-page-header>
      </div>

      <div class="px-2 md:px-4 lg:px-10" v-if="!reverting">
        <n-card>
          <records :data="collection" />
        </n-card>
      </div>
    </div>
  </n-spin>
</template>

<script lang="ts" setup>
import Records from "@/views/parts/records.vue"
import { useMessage } from "naive-ui"
import { useRoute } from "vue-router"
import { onMounted, ref } from "vue"
import { http } from "@/utils/http"

const $route = useRoute()
const $message = useMessage()

const app = ref<any>({})
const collection = ref<any>({})

const reverting = ref(true)

async function fetch() {
  try {
    reverting.value = true
    app.value = (await http.get(`/api/apps/${$route.params.app}`)).data
    collection.value = (await http.get(`/api/apps/${$route.params.app}/records/${$route.params.collection}`)).data
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
