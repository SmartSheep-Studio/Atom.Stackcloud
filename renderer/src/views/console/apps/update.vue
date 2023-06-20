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
        <v-text-field
          variant="outlined"
          density="comfortable"
          label="Tags"
          v-model="payload.tags"
          :hide-details="true"
        />
      </v-col>
      <v-col :md="8" :sm="12" :cols="12">
        <v-text-field variant="outlined" density="comfortable" label="URL" v-model="payload.url" :hide-details="true" />
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

const emits = defineEmits(["done"])
const props = defineProps<{
  data: {
    id: number
    slug: string
    name: string
    description: string
    tags: string
    url: string
    details: string
    is_published: boolean
  }
}>()

const $snackbar = useSnackbar()

const submitting = ref(false)
const payload = ref<any>({})

async function create() {
  payload.value.tags = payload.value.tags.split(",")

  try {
    submitting.value = true

    const res = await http.put(`/api/apps/${payload.value.slug}`, payload.value)
    emits("done")

    $snackbar.show({ text: `Successfully updated app ${res.data.name}.`, color: "success" })
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  payload.value = JSON.parse(JSON.stringify(props.data))
  payload.value.tags = payload.value.tags.join(",")
})
</script>
