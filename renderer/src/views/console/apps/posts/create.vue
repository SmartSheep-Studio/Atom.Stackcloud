<template>
  <v-form class="py-4" @submit.prevent="create">
    <v-row justify="center">
      <v-col :md="8" :sm="12" :cols="12">
        <v-text-field
          variant="outlined"
          density="comfortable"
          label="Title"
          v-model="payload.title"
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
        <v-md-editor v-model="payload.content" height="400px" />
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

const types = [
  { name: "Announcement", value: "announcement" },
  { name: "Blog", value: "blog" }
]

const submitting = ref(false)
const payload = ref<any>({
  id: 0,
  title: "",
  tags: "",
  type: "blog",
  content: "",
  is_published: false
})

async function create() {
  const data: any = JSON.parse(JSON.stringify(payload.value))
  data.tags = payload.value.tags.split(",")

  try {
    submitting.value = true

    const res = await http.post(`/api/apps/${$route.params.app}/posts`, data)
    emits("done")

    $snackbar.show({ text: `Successfully created post ${res.data.name}.`, color: "success" })
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
    title: "",
    tags: "",
    type: "blog",
    content: "",
    is_published: false
  }
}
</script>
