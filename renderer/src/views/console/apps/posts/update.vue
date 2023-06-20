<template>
  <v-form class="py-4" @submit.prevent="update">
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
import { onMounted } from "vue"
import { useRoute } from "vue-router"

const emits = defineEmits(["done"])
const props = defineProps<{
  data: {
    id: number
    title: string
    tags: string
    content: string
    is_published: boolean
  }
}>()

const $route = useRoute()
const $snackbar = useSnackbar()

const types = [
  { name: "Announcement", value: "announcement" },
  { name: "Blog", value: "blog" }
]

const submitting = ref(false)
const payload = ref<any>({})

async function update() {
  const data: any = JSON.parse(JSON.stringify(payload.value))
  data.tags = payload.value.tags.split(",")

  try {
    submitting.value = true

    const res = await http.put(`/api/apps/${$route.params.app}/posts/${payload.value.id}`, data)
    emits("done")

    $snackbar.show({ text: `Successfully updated post ${res.data.name}.`, color: "success" })
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
