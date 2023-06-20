<template>
  <v-container fluid class="pa-12">
    <div>
      <div class="text-h6">{{ app?.name ?? "Loading..." }}</div>
      <div class="text-subtitle-1">App configuration</div>
    </div>

    <div class="mt-4">
      <releases />
    </div>

    <div class="mt-4">
      <posts />
    </div>

    <v-footer class="border-t" app>
      <v-breadcrumbs density="compact">
        <v-breadcrumbs-item exact :to="{ name: 'console.dashboard' }">Dashboard</v-breadcrumbs-item>
        <v-breadcrumbs-divider />
        <v-breadcrumbs-item disabled>{{ app.name }}</v-breadcrumbs-item>
      </v-breadcrumbs>
    </v-footer>
  </v-container>
</template>

<script lang="ts" setup>
import { ref } from "vue"
import { http } from "@/utils/http";
import { useRoute } from "vue-router";
import { useSnackbar } from "@/stores/snackbar";
import { onMounted } from "vue";
import Releases from "@/views/console/apps/releases.vue"
import Posts from "@/views/console/apps/posts.vue"

const $route = useRoute()
const $snackbar = useSnackbar()

const app = ref<any>({})

async function fetch() {
  try {
    app.value = (await http.get(`/api/apps/${$route.params.app}`)).data
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  }
}

onMounted(() => {
  fetch()
})
</script>
