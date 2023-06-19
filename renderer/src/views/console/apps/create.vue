<template>
  <v-form class="py-4" @submit.prevent="create">
    <v-row justify="center">
      <v-col :md="4" :sm="6">
        <v-text-field
          variant="outlined"
          density="comfortable"
          label="Slug"
          v-model="payload.slug"
          :hide-details="true"
        />
      </v-col>
      <v-col :md="4" :sm="6">
        <v-text-field
          variant="outlined"
          density="comfortable"
          label="Name"
          v-model="payload.name"
          :hide-details="true"
        />
      </v-col>
      <v-col :md="8" :sm="12">
        <v-text-field
          variant="outlined"
          density="comfortable"
          label="Tags"
          v-model="payload.tags"
          :hide-details="true"
        />
      </v-col>
      <v-col :md="8" :sm="12">
        <v-text-field variant="outlined" density="comfortable" label="URL" v-model="payload.url" :hide-details="true" />
      </v-col>
      <v-col :md="8" :sm="12">
        <v-textarea
          variant="outlined"
          density="comfortable"
          label="Description"
          v-model="payload.description"
          :hide-details="true"
        />
      </v-col>
      <v-col :md="8" :sm="12">
        <v-md-editor v-model="payload.details" height="400px" />
      </v-col>
      <v-col :md="8" :sm="12">
        <v-checkbox-btn label="Published" v-model="payload.is_published" color="primary" />
      </v-col>
      <v-col :md="8" :sm="12">
        <v-btn color="primary" type="submit" variant="text" :loading="submitting">Submit</v-btn>
      </v-col>
    </v-row>
  </v-form>
</template>

<script lang="ts" setup>
import { ref } from "vue"
import { http } from "@/utils/http"
import { useSnackbar } from "@/stores/snackbar"

const emits = defineEmits(["done"])

const $snackbar = useSnackbar()

const submitting = ref(false)
const payload = ref<any>({
  id: 0,
  slug: "",
  name: "",
  description: "",
  tags: "",
  url: "",
  details: "",
  is_published: false
})

async function create() {
  payload.value.tags = payload.value.tags.split(",")

  try {
    submitting.value = true

    const res = await http.post("/api/apps", payload.value)
    emits("done")

    $snackbar.show({ text: `Successfully created app ${res.data.name}.`, color: "success" })
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
