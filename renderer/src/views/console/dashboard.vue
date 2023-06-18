<template>
  <v-container fluid class="pa-12">
    <div>
      <div class="text-h6">Merchant Dashboard</div>
      <div class="text-subtitle-1">Welcome, {{ $account.profile?.nickname }}!</div>
    </div>

    <div class="mt-4">
      <v-card>
        <v-toolbar dark density="compact" color="blue-darken-3" class="px-2 mb-2 row">
          <v-toolbar-title style="font-size: 16px">Apps</v-toolbar-title>
          <v-col :offset="4" :cols="4" class="d-flex justify-end">
            <v-dialog v-model="popups.create" width="420px">
              <template #activator="{ props }">
                <v-btn color="white" v-bind="props">New</v-btn>
              </template>

              <v-card title="Create Shop" prepend-icon="mdi-plus">
                <v-card-item class="pt-0">
                  <v-form class="pt-4" @submit.prevent="create">
                    <v-text-field variant="outlined" density="comfortable" label="Slug" name="slug" />
                    <v-text-field variant="outlined" density="comfortable" label="Name" name="name" />
                    <v-text-field variant="outlined" density="comfortable" label="URL" name="url" />
                    <v-textarea variant="outlined" density="comfortable" label="Description" name="description" />
                    <v-text-field variant="outlined" density="comfortable" label="Tags" name="tags" />
                    <div class="mb-4">
                      <v-btn block color="primary" type="submit" variant="text" :loading="submitting">Submit</v-btn>
                    </div>
                  </v-form>
                </v-card-item>
              </v-card>
            </v-dialog>
          </v-col>
        </v-toolbar>

        <v-table density="comfortable" hover height="60vh">
          <tbody>
            <tr v-for="(item, i) in shops" :key="item.name">
              <td style="width: 64px">
                <v-avatar>
                  <v-icon icon="mdi-application-brackets-outline" />
                </v-avatar>
              </td>
              <td>
                <div class="pt-2">{{ item.name }}</div>
                <div class="pb-2 text-caption">{{ item.description }}</div>
              </td>
              <td style="text-align: right">
                <div class="me-2">
                  <v-dialog v-model="popups.update[i]" width="420px">
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

                    <v-card title="Edit Shop" prepend-icon="mdi-pencil">
                      <v-card-item class="pt-0">
                        <v-form class="pt-4" @submit.prevent="update">
                          <v-text-field density="comfortable" variant="outlined" v-model="payload.slug" label="Slug" />
                          <v-text-field density="comfortable" variant="outlined" v-model="payload.name" label="Name" />
                          <v-text-field variant="outlined" density="comfortable" v-model="payload.url" label="URL" />
                          <v-textarea
                            density="comfortable"
                            variant="outlined"
                            v-model="payload.description"
                            label="Description"
                          />
                          <v-text-field variant="outlined" density="comfortable" label="Tags" v-model="payload.tags" />
                          <div class="mb-4">
                            <v-btn block color="primary" type="submit" variant="text" :loading="submitting">
                              Submit
                            </v-btn>
                          </div>
                        </v-form>
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
                  <v-btn
                    icon="mdi-login-variant"
                    size="small"
                    color="primary"
                    variant="text"
                    :to="{ name: 'console.shop', params: { shop: item.slug } }"
                  />
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

const $account = useAccount()
const $snackbar = useSnackbar()

const shops = ref<any[]>([])
const methods = ref<any[]>([])

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
    shops.value = (await http.get("/api/shops")).data
    methods.value = (await http.get("/api/methods")).data
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  }
}

async function create(e: SubmitEvent) {
  const payload: any = Object.fromEntries(new FormData(e.target as HTMLFormElement).entries())
  payload.tags = payload.tags.split(",")

  try {
    submitting.value = true

    const res = await http.post("/api/apps", payload)
    await fetch()

    popups.create = false
    $snackbar.show({ text: `Successfully created app ${res.data.name}.`, color: "success" })
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  } finally {
    submitting.value = false
  }
}

async function update() {
  const data: any = payload.value
  data.tags = data.tags.split(",")

  try {
    submitting.value = true

    const res = await http.put(`/api/apps/${payload.value.slug}`, data)
    await fetch()

    popups.update = popups.update.map(() => false)
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
  payload.value.tags = (payload.value.tags as any).json(",")
}

onMounted(() => {
  fetch()
})
</script>
