<template>
  <div>
    <n-spin :show="requesting">
      <n-list bordered hoverable>
        <n-empty v-if="data.length <= 0" class="py-8" description="There's no data. Why not you create one?" />

        <n-list-item v-for="(item, i) in data" :key="i">
          <n-thing
            class="cursor-pointer"
            :title="item.name"
            @click="$router.push({ name: 'console.apps', params: { app: item.slug } })"
          >
            <template #description>
              <n-space size="small">
                <n-tag :bordered="false" type="success" size="small">#{{ item.slug }}</n-tag>
                <n-tag v-for="(tag, i) in item.tags" :key="i" :bordered="false" type="primary" size="small">
                  {{ tag }}
                </n-tag>
              </n-space>
            </template>
            <div>{{ item.description }}</div>
          </n-thing>
        </n-list-item>
      </n-list>
    </n-spin>

    <div class="flex justify-between mt-4">
      <n-button type="primary" size="small" @click="$router.push({ name: 'console.apps.create' })">
        <template #icon>
          <n-icon :component="PlusRound" />
        </template>
        New App
      </n-button>

      <n-pagination
        v-model:page="pagination.page"
        :page-count="Math.ceil((rawData.length ?? 0) / pagination.pageSize)"
        :page-slot="pagination.slot"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from "vue"
import { PlusRound } from "@vicons/material"
import { http } from "@/utils/http"
import { useMessage } from "naive-ui"

const $message = useMessage()

const rawData = ref<any[]>([])
const data = computed(() => {
  const start = (pagination.page - 1) * pagination.pageSize
  return rawData.value.reverse().slice(start, start + pagination.pageSize) ?? []
})

const requesting = ref(true)

const pagination = reactive({
  page: 1,
  pageSize: 5,
  slot: 5,
})

async function fetch() {
  try {
    requesting.value = true
    rawData.value = (await http.get("/api/apps")).data
  } catch (e: any) {
    $message.error(`Something went wrong... ${e}`)
  } finally {
    requesting.value = false
  }
}

onMounted(() => {
  fetch()
})
</script>
