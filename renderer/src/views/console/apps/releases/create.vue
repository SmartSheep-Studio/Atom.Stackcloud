<template>
  <v-form class="py-4" @submit.prevent="create">
    <v-row justify="center">
      <v-col :md="4" :sm="6" :cols="12">
        <v-text-field
          variant="outlined"
          density="comfortable"
          label="Slug"
          v-model="payload.slug"
          :hide-details="true"
        />
      </v-col>
      <v-col :md="4" :sm="6" :cols="12">
        <v-text-field
          variant="outlined"
          density="comfortable"
          label="Name"
          v-model="payload.name"
          :hide-details="true"
        />
      </v-col>
      <v-col :md="8" :sm="12" :cols="12">
        <v-select
          variant="outlined"
          density="comfortable"
          label="Type"
          item-title="name"
          item-value="value"
          v-model="payload.type"
          :items="types"
          :hide-details="true"
        />
      </v-col>
      <v-col :md="8" :sm="12" :cols="12">
        <v-text-field
          variant="outlined"
          density="comfortable"
          label="Tags"
          v-model="payload.tags"
          :hide-details="true"
        />
      </v-col>
      <v-col :md="8" :sm="12" :cols="12">
        <v-textarea
          variant="outlined"
          density="comfortable"
          label="Description"
          v-model="payload.description"
          :hide-details="true"
        />
      </v-col>
      <v-col :md="8" :sm="12" :cols="12">
        <v-md-editor v-model="payload.details" height="400px" />
      </v-col>
      <v-col :md="8" :sm="12" :cols="12">
        <v-card prepend-icon="mdi-play" title="Run & Play Options">
          <v-card-item class="pt-0 px-6">
            <div class="mt-2">
              <v-card v-for="(item, i) in payload.options.assets" :key="i" variant="outlined" class="mb-2">
                <v-card-item>
                  <div class="text-subtitle-2 px-2 pb-3">Assets #{{ i + 1 }}</div>
                  <v-text-field density="compact" label="URL" variant="outlined" v-model="item.url" />
                  <v-text-field density="compact" label="Decompressor" variant="outlined" v-model="item.decompressor" />
                  <v-select
                    density="compact"
                    label="Platform"
                    variant="outlined"
                    :items="platforms"
                    v-model="item.platform"
                  />
                  <v-btn
                    variant="text"
                    prepend-icon="mdi-close"
                    color="error"
                    @click="payload.options.assets.splice(i, 1)"
                  >
                    Remove
                  </v-btn>
                </v-card-item>
              </v-card>
              <v-card v-for="(item, i) in payload.options.preprocessing" :key="i" variant="outlined" class="mb-2">
                <v-card-item>
                  <div class="text-subtitle-2 px-2 pb-3">Preprocessing #{{ i + 1 }}</div>
                  <v-text-field density="compact" label="Script" variant="outlined" v-model="item.script" />
                  <v-select
                    density="compact"
                    label="Platform"
                    variant="outlined"
                    :items="platforms"
                    v-model="item.platform"
                  />
                  <v-btn
                    variant="text"
                    prepend-icon="mdi-close"
                    color="error"
                    @click="payload.options.preprocessing.splice(i, 1)"
                  >
                    Remove
                  </v-btn>
                </v-card-item>
              </v-card>
              <v-card v-for="(item, i) in payload.options.run_options" :key="i" variant="outlined" class="mb-2">
                <v-card-item>
                  <div class="text-subtitle-2 px-2 pb-3">Run Options #{{ i + 1 }}</div>
                  <v-text-field density="compact" label="Name" variant="outlined" v-model="item.name" />
                  <v-text-field density="compact" label="Script" variant="outlined" v-model="item.script" />
                  <v-select
                    density="compact"
                    label="Platform"
                    variant="outlined"
                    :items="platforms"
                    v-model="item.platform"
                  />
                  <v-btn
                    variant="text"
                    prepend-icon="mdi-close"
                    color="error"
                    @click="payload.options.run_options.splice(i, 1)"
                  >
                    Remove
                  </v-btn>
                </v-card-item>
              </v-card>
            </div>
          </v-card-item>
          <v-card-item class="px-6 pt-0">
            <v-btn
              variant="text"
              prepend-icon="mdi-script-text"
              color="primary"
              @click="payload.options.assets.push({ url: '', decompressor: '', platform: 'win32' })"
            >
              Add Assets
            </v-btn>
            <v-btn
              variant="text"
              prepend-icon="mdi-download-box"
              color="teal"
              @click="payload.options.preprocessing.push({ script: '', platform: 'win32' })"
            >
              Add Preprocessing
            </v-btn>
            <v-btn
              variant="text"
              prepend-icon="mdi-play"
              color="green"
              @click="payload.options.run_options.push({ name: '', script: '', platform: 'win32' })"
            >
              Add Run Option
            </v-btn>
          </v-card-item>
        </v-card>
      </v-col>
      <v-col :md="8" :sm="12" :cols="12">
        <v-checkbox-btn label="Published" v-model="payload.is_published" color="primary" />
      </v-col>
      <v-col :md="8" :sm="12" :cols="12">
        <v-btn color="primary" type="submit" variant="text" :loading="submitting">Submit</v-btn>
      </v-col>
    </v-row>
  </v-form>
</template>

<script lang="ts" setup>
import { ref } from "vue"
import { http } from "@/utils/http"
import { useSnackbar } from "@/stores/snackbar"
import { useRoute } from "vue-router"

const emits = defineEmits(["done"])

const $route = useRoute()
const $snackbar = useSnackbar()

const platforms = ["win32", "drawin", "linux", "freebsd"]
const types = [
  { name: "Minor Update", value: "minor-update" },
  { name: "Major Update", value: "major-update" },
  { name: "Hotfix Update", value: "hotfix-update" }
]

const submitting = ref(false)
const payload = ref<any>({
  id: 0,
  slug: "",
  name: "",
  tags: "",
  type: "minor-update",
  description: "",
  details: "",
  options: {
    assets: [],
    preprocessing: [],
    run_options: []
  },
  is_published: false
})

async function create() {
  const data: any = JSON.parse(JSON.stringify(payload.value))
  data.tags = payload.value.tags.split(",")

  try {
    submitting.value = true

    const res = await http.post(`/api/apps/${$route.params.app}/releases`, data)
    emits("done")

    $snackbar.show({ text: `Successfully created release ${res.data.name}.`, color: "success" })
    reset()
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  } finally {
    submitting.value = false
  }
}

function reset() {
  payload.value = {
    id: 0,
    slug: "",
    name: "",
    description: "",
    tags: "",
    url: "",
    details: "",
    is_published: false
  }
}
</script>
