<template>
  <div>
    <n-spin :show="requesting">
      <n-list bordered hoverable>
        <n-empty v-if="data.length <= 0" class="py-8" description="There's no data. Why not you create one?" />

        <n-list-item v-for="(item, i) in data" :key="i">
          <n-thing :title="item.name">
            <template #description>
              <n-space size="small">
                <n-tag :bordered="false" type="success" size="small">#{{ item.slug }}</n-tag>
                <n-tag v-for="(tag, i) in item.tags" :key="i" :bordered="false" type="primary" size="small">
                  {{ tag }}
                </n-tag>
              </n-space>
            </template>
            <div>{{ item.description }}</div>
            <n-space class="mt-2" size="small">
              <n-button
                type="info"
                size="small"
                @click="
                  $router.push({
                    name: 'console.apps.collections',
                    params: { app: props.data.slug, collection: item.slug },
                  })
                "
              >
                <template #icon>
                  <n-icon :component="TextSnippetRound" />
                </template>
                Inspect
              </n-button>
              <n-button
                type="warning"
                size="small"
                @click="
                  $router.push({
                    name: 'console.apps.collections.update',
                    params: { app: props.data.slug, collection: item.slug },
                  })
                "
              >
                <template #icon>
                  <n-icon :component="EditRound" />
                </template>
                Edit
              </n-button>
              <n-button type="error" size="small" @click="destroy(item)">
                <template #icon>
                  <n-icon :component="DeleteRound" />
                </template>
                Destroy
              </n-button>
            </n-space>
          </n-thing>
        </n-list-item>
      </n-list>
    </n-spin>

    <div class="flex justify-between mt-4">
      <n-button type="primary" size="small" @click="$router.push({ name: 'console.apps.collections.create' })">
        <template #icon>
          <n-icon :component="PlusRound" />
        </template>
        New Collection
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
import { PlusRound, TextSnippetRound, EditRound, DeleteRound } from "@vicons/material"
import { useMessage, useDialog } from "naive-ui"
import { http } from "@/utils/http"

const $dialog = useDialog()
const $message = useMessage()

const props = defineProps<{ data: any }>()

const rawData = ref<any[]>([])
const data = computed(() => {
  const start = (pagination.page - 1) * pagination.pageSize
  return rawData.value.slice(start, start + pagination.pageSize) ?? []
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
    rawData.value = (await http.get(`/api/apps/${props.data.slug}/records`)).data
  } catch (e: any) {
    $message.error(`Something went wrong... ${e}`)
  } finally {
    requesting.value = false
  }
}

function destroy(item: any) {
  $dialog.warning({
    title: "Warning",
    content: "This operation cannot be undo. Are you confirm?",
    positiveText: "Yes",
    negativeText: "Not really",
    onPositiveClick: async () => {
      try {
        requesting.value = true

        await http.delete(`/api/apps/${props.data.slug}/records/${item.slug}`)
        await fetch()

        $message.success("Successfully deleted the collection.")
      } catch (e: any) {
        $message.error(`Something went wrong... ${e}`)
      } finally {
        requesting.value = false
      }
    },
  })
}

onMounted(() => {
  fetch()
})
</script>
