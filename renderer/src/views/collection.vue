<template>
  <div>
    <n-spin :show="reverting">
      <n-tabs type="line" justify-content="space-around" animated>
        <n-tab-pane name="data" tab="Data">
          <div class="px-8 py-4" v-if="!reverting">
            <n-card>
              <records :data="collection" />
            </n-card>
          </div>
        </n-tab-pane>
        <n-tab-pane name="settings" tab="Settings">
          <div class="px-8 py-4 flex justify-center" v-if="!reverting">
            <n-grid item-responsive responsive="screen" x-gap="8" y-gap="8">
              <n-gi span="24 m:14 l:16">
                <n-card title="Update Collection">
                  <update-collection />
                </n-card>
              </n-gi>
              <n-gi span="24 m:10 l:8">
                <n-card title="Dangerous Zone">
                  <destroy-collection
                    :data="collection"
                    @done="
                      $router.push({ name: 'console.apps', params: { app: $route.params.app } }).then(() => reload())
                    "
                  />
                </n-card>
              </n-gi>
            </n-grid>
          </div>
        </n-tab-pane>
      </n-tabs>
    </n-spin>
  </div>
</template>

<script lang="ts" setup>
import Records from "@/views/parts/records.vue"
import UpdateCollection from "@/views/actions/update-collection.vue"
import DestroyCollection from "@/views/actions/destroy-collection.vue"
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

function reload() {
  window.location.reload()
}

onMounted(() => {
  fetch()
})
</script>
