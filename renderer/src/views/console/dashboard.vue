<template>
  <v-container fluid class="pa-12">
    <div>
      <div class="text-h6">Dashboard</div>
      <div class="text-subtitle-1">Welcome, {{ $account.profile?.nickname }}!</div>
    </div>

    <div class="mt-4">
      <v-card>
        <v-toolbar dark density="compact" color="blue-darken-3" class="px-2 mb-2 row">
          <v-toolbar-title style="font-size: 16px">Apps</v-toolbar-title>
          <v-col :offset="4" :cols="4" class="d-flex justify-end">
            <v-dialog fullscreen transition="dialog-bottom-transition" v-model="popups.create">
              <template #activator="{ props }">
                <v-btn color="white" v-bind="props">New</v-btn>
              </template>

              <v-card>
                <v-toolbar dark color="primary" absolute>
                  <v-btn icon dark @click="popups.create = false">
                    <v-icon>mdi-close</v-icon>
                  </v-btn>
                  <v-toolbar-title style="text-align: center; flex-grow: 1; padding-right: 74px">
                    Create App
                  </v-toolbar-title>
                </v-toolbar>

                <v-card-item style="padding-top: 64px; max-height: 100vh; overflow: auto">
                  <create-app @done="refresh" />
                </v-card-item>
              </v-card>
            </v-dialog>
          </v-col>
        </v-toolbar>

        <v-table density="comfortable" hover height="60vh">
          <tbody>
            <tr v-for="(item, i) in apps" :key="item.name">
              <td style="width: 64px">
                <v-avatar>
                  <v-icon icon="mdi-application-brackets-outline" />
                </v-avatar>
              </td>
              <td>
                <div class="py-2">
                  <div>{{ item.name }}</div>
                  <div class="text-caption">{{ item.description }}</div>
                  <div class="text-caption mt-1" style="margin-left: -4px">
                    <v-badge v-for="tag in item.tags" inline color="primary" :content="tag" />
                  </div>
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
                      <v-toolbar dark color="primary" absolute>
                        <v-btn icon dark @click="popups.update[i] = false">
                          <v-icon>mdi-close</v-icon>
                        </v-btn>
                        <v-toolbar-title style="text-align: center; flex-grow: 1; padding-right: 74px">
                          Edit App
                        </v-toolbar-title>
                      </v-toolbar>

                      <v-card-item style="padding-top: 64px; max-height: 100vh; overflow: auto">
                        <update-app @done="refresh" :data="item" />
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
                      <v-card-text>Do you sure you want to delete this app?</v-card-text>
                      <v-card-actions class="px-5 justify-end">
                        <v-btn @click="popups.delete[i] = false">Cancel</v-btn>
                        <v-btn :loading="submitting" color="red" @click="dispose">Confirm</v-btn>
                      </v-card-actions>
                    </v-card>
                  </v-dialog>
                  <v-btn icon="mdi-login-variant" size="small" color="primary" variant="text" />
                </div>
              </td>
            </tr>
          </tbody>
        </v-table>
      </v-card>
    </div>

    <v-footer class="border-t" app>
      <v-breadcrumbs density="compact">
        <v-breadcrumbs-item disabled>Dashboard</v-breadcrumbs-item>
      </v-breadcrumbs>
    </v-footer>
  </v-container>
</template>

<script lang="ts" setup>
import { useAccount } from "@/stores/account"
import { useSnackbar } from "@/stores/snackbar"
import { http } from "@/utils/http"
import { reactive } from "vue"
import { onMounted } from "vue"
import { ref } from "vue"
import CreateApp from "@/views/console/apps/create.vue"
import UpdateApp from "@/views/console/apps/update.vue"

const $account = useAccount()
const $snackbar = useSnackbar()

const apps = ref<any[]>([])

const payload = ref({
  id: 0,
  slug: "",
  name: "",
  description: "",
  tags: "",
  url: "",
  details: ""
})

const submitting = ref(false)
const popups = reactive({ create: false, update: [false], delete: [false] })

async function fetch() {
  try {
    apps.value = (await http.get("/api/apps")).data
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  }
}

async function refresh() {
  popups.create = false
  popups.update = popups.update.map(() => false)
  await fetch()
}

async function update() {
  const data: any = payload.value
  data.tags = data.tags.split(",")

  try {
    submitting.value = true

    const res = await http.put(`/api/apps/${payload.value.slug}`, data)
    await fetch()

    popups.update = popups.update.map(() => false)
    payload.value.details = ""
    $snackbar.show({ text: `Successfully updated app ${res.data.name}.`, color: "success" })
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  } finally {
    submitting.value = false
  }
}

async function dispose() {
  try {
    submitting.value = true

    await http.delete(`/api/apps/${payload.value.slug}`)
    await fetch()

    popups.delete = popups.delete.map(() => false)
    $snackbar.show({ text: `Successfully deleted app ${payload.value.name}.`, color: "success" })
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  } finally {
    submitting.value = false
  }
}

function refer(v: any) {
  payload.value = JSON.parse(JSON.stringify(v))
  payload.value.tags = (payload.value.tags as any).join(",")
}

onMounted(() => {
  fetch()
})
</script>
