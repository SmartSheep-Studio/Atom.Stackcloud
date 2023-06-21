<template>
  <div>
    <v-card>
      <v-toolbar dark density="compact" color="blue-darken-3" class="px-2 mb-2 row">
        <v-toolbar-title style="font-size: 16px">Releases</v-toolbar-title>
        <v-col :offset="4" :cols="4" class="d-flex justify-end">
          <v-dialog fullscreen transition="dialog-bottom-transition" v-model="popups.create">
            <template #activator="{ props }">
              <v-btn color="white" v-bind="props">New</v-btn>
            </template>

            <v-card>
                <v-toolbar dark color="primary">
                  <v-btn icon dark @click="popups.create = false">
                    <v-icon>mdi-close</v-icon>
                  </v-btn>
                  <v-toolbar-title style="text-align: center; flex-grow: 1; padding-right: 74px">
                    Create Release
                  </v-toolbar-title>
                </v-toolbar>

                <v-card-item style="max-height: calc(100vh - 64px); overflow: auto">
                  <create-release @done="refresh" />
                </v-card-item>
              </v-card>
          </v-dialog>
        </v-col>
      </v-toolbar>

      <v-table density="comfortable" hover height="30vh">
        <tbody>
          <tr v-for="(item, i) in releases" :key="item.name">
            <td style="width: 64px">
              <v-avatar>
                <v-icon icon="mdi-tag" />
              </v-avatar>
            </td>
            <td>
              <div class="py-2">
                <div>{{ item.name }}</div>
                <div class="text-caption">{{ item.description }}</div>
              </div>
            </td>
            <td style="text-align: right">
              <div class="me-2">
                <v-dialog fullscreen transition="dialog-bottom-transition" v-model="popups.update[i]">
                  <template #activator="{ props }">
                    <v-btn
                      icon="mdi-pencil"
                      size="small"
                      color="warning"
                      variant="text"
                      v-bind="props"
                      @click="refer(item)"
                    />
                  </template>

                  <v-card>
                    <v-toolbar dark color="primary">
                      <v-btn icon dark @click="popups.update[i] = false">
                        <v-icon>mdi-close</v-icon>
                      </v-btn>
                      <v-toolbar-title style="text-align: center; flex-grow: 1; padding-right: 74px">
                        Edit Release
                      </v-toolbar-title>
                    </v-toolbar>

                    <v-card-item style="max-height: calc(100vh - 64px); overflow: auto">
                      <update-release @done="refresh" :data="item" />
                    </v-card-item>
                  </v-card>
                </v-dialog>
                <v-dialog v-model="popups.delete[i]" width="420px">
                  <template #activator="{ props }">
                    <v-btn
                      icon="mdi-delete"
                      size="small"
                      color="error"
                      variant="text"
                      v-bind="props"
                      @click="refer(item)"
                    />
                  </template>
                  <v-card title="Confirm">
                    <v-card-text>Do you sure you want to delete this release?</v-card-text>
                    <v-card-actions class="px-5 justify-end">
                      <v-btn @click="popups.delete[i] = false">Cancel</v-btn>
                      <v-btn :loading="submitting" color="red" @click="dispose">Confirm</v-btn>
                    </v-card-actions>
                  </v-card>
                </v-dialog>
              </div>
            </td>
          </tr>
        </tbody>
      </v-table>
    </v-card>
  </div>
</template>

<script lang="ts" setup>
import { useSnackbar } from "@/stores/snackbar"
import { http } from "@/utils/http"
import { onMounted, reactive } from "vue"
import { ref } from "vue"
import { useRoute } from "vue-router"
import CreateRelease from "@/views/console/apps/releases/create.vue"
import UpdateRelease from "@/views/console/apps/releases/update.vue"

const $route = useRoute()
const $snackbar = useSnackbar()

const releases = ref<any[]>([])

const submitting = ref(false)
const popups = reactive({ create: false, update: [false], delete: [false] })

const payload = ref({
  id: 0,
  slug: "",
  name: "",
  tags: "",
  type: "minor-update",
  description: "",
  details: "",
  is_published: false
})

async function fetch() {
  try {
    releases.value = (await http.get(`/api/apps/${$route.params.app}/releases`)).data
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  }
}

async function refresh() {
  popups.create = false
  popups.update = popups.update.map(() => false)
  await fetch()
}

async function dispose() {
  try {
    submitting.value = true

    await http.delete(`/api/apps/${$route.params.app}/releases/${payload.value.id}`)
    await fetch()

    popups.delete = popups.delete.map(() => false)
    $snackbar.show({ text: `Successfully deleted release ${payload.value.name}.`, color: "success" })
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  } finally {
    submitting.value = false
  }
}

function refer(v: any) {
  payload.value = JSON.parse(JSON.stringify(v))
}

onMounted(() => {
  fetch()
})
</script>
