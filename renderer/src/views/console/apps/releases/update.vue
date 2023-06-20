<template>
  <v-form class="py-4" @submit.prevent="update">
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
import { onMounted } from "vue"
import { useRoute } from "vue-router"

const emits = defineEmits(["done"])
const props = defineProps<{
  data: {
    id: number
    slug: string
    name: string
    tags: string
    description: string
    details: string
    is_published: boolean
  }
}>()

const $route = useRoute()
const $snackbar = useSnackbar()

const types = [
  { name: "Minor Update", value: "minor-update" },
  { name: "Major Update", value: "major-update" },
  { name: "Hotfix Update", value: "hotfix-update" }
]

const submitting = ref(false)
const payload = ref<any>({})

async function update() {
  const data: any = JSON.parse(JSON.stringify(payload.value))
  data.tags = payload.value.tags.split(",")

  try {
    submitting.value = true

    const res = await http.put(`/api/apps/${$route.params.app}/releases/${payload.value.id}`, data)
    emits("done")

    $snackbar.show({ text: `Successfully updated release ${res.data.name}.`, color: "success" })
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  payload.value = JSON.parse(JSON.stringify(props.data))
  payload.value.details = payload.value.post.content
  payload.value.type = payload.value.post.type
  payload.value.tags = payload.value.post.tags.join(",")
})
</script>
