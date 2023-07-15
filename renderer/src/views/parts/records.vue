<template>
  <div>
    <div class="flex justify-between mb-4">
      <div class="text-lg">Collection Records</div>
      <n-button type="primary" size="small" @click="creating = true">
        <template #icon>
          <n-icon :component="PlusRound" />
        </template>
        New Record
      </n-button>
    </div>

    <n-spin :show="requesting">
      <div>
        <n-data-table :columns="tableColumns" :data="data" :pagination="{ pageSize: 20 }" />
      </div>
    </n-spin>

    <n-modal v-model:show="creating">
      <n-card size="huge" title="Create a new record" class="w-[800px]">
        <create-record @close="creating = false" @refresh="fetch" />
      </n-card>
    </n-modal>

    <n-modal v-model:show="updating">
      <n-card size="huge" title="Update a exists record" class="w-[800px]">
        <update-record :data="focus" @close="updating = false" @refresh="fetch" />
      </n-card>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import CreateRecord from "@/views/actions/create-record.vue"
import UpdateRecord from "@/views/actions/update-record.vue"
import hljs from "highlight.js/lib/core"
import json from "highlight.js/lib/languages/json"
import { computed, onMounted, ref, h } from "vue"
import { PlusRound } from "@vicons/material"
import { NCode, NSpace, NButton, useMessage, useDialog, type DataTableColumns } from "naive-ui"
import { http } from "@/utils/http"
import { useRoute } from "vue-router"

hljs.registerLanguage("json", json)

const $route = useRoute()
const $dialog = useDialog()
const $message = useMessage()

const props = defineProps<{ data: any }>()

const rawData = ref<any[]>([])
const data = computed(() => rawData.value)

const focus = ref<null | any>(null)
const creating = ref(false)
const updating = ref(false)
const requesting = ref(true)

const tableColumns: DataTableColumns<any> = [
  { title: "ID", key: "id" },
  {
    title: "Payload",
    key: "payload",
    render(row) {
      return h(NCode, {
        hljs,
        code: JSON.stringify(row.payload),
        language: "json",
        wordWrap: true,
      })
    },
  },
  { title: "Created At", key: "created_at", sorter: "default" },
  { title: "Updated At", key: "updated_at", sorter: "default" },
  {
    title: "Actions",
    key: "actions",
    render(row) {
      return h(
        NSpace,
        {
          size: "small",
        },
        {
          default: () => [
            h(
              NButton,
              {
                size: "small",
                type: "warning",
                onClick: () => {
                  focus.value = row
                  updating.value = true
                },
              },
              { default: () => "Update" }
            ),
            h(
              NButton,
              {
                size: "small",
                type: "error",
                onClick: () => destroy(row),
              },
              { default: () => "Destroy" }
            ),
          ],
        }
      )
    },
  },
]

async function fetch() {
  try {
    requesting.value = true
    rawData.value = (await http.get(`/api/apps/${$route.params.app}/records/${props.data.slug}/data`)).data
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

        await http.delete(`/api/apps/${$route.params.app}/records/${props.data.slug}/data/${item.id}`)
        await fetch()

        $message.success("Successfully deleted the record.")
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
