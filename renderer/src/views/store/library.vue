<template>
  <v-container fluid class="pa-12">
    <div>
      <div class="text-h6">Library</div>
      <div class="text-subtitle-1">Your library</div>
    </div>

    <div class="mt-4" v-if=reverting>
      <v-progress-circular indeterminate />
    </div>
    <div class="mt-6" v-else>
      <div style="margin-left: -8px; margin-right: -8px">
        <v-data-iterator :items="library" :page="pagination.app">
          <template v-slot:default="{ items }">
            <v-row class="explore-queue">
              <template v-for="(item, i) in items" :key="i">
                <v-col :cols="12" :sm="6" :md="4">
                  <v-card class="explore-item pa-4">
                    <v-card-item>
                      <div class="text-h6">{{ apps[item.raw.app_id].name }}</div>
                      <div class="explore-tags" style="margin-left: -4px">
                        <v-badge v-for="tag in apps[item.raw.app_id].tags" inline color="primary" :content="tag" />
                      </div>
                    </v-card-item>
                    <v-card-item>
                      {{ apps[item.raw.app_id].description }}
                    </v-card-item>
                  </v-card>
                </v-col>
              </template>
            </v-row>
          </template>
        </v-data-iterator>
      </div>
    </div>
  </v-container>
</template>

<script lang="ts" setup>
import { useSnackbar } from "@/stores/snackbar"
import { http } from "@/utils/http"
import { onMounted } from "vue";
import { reactive, ref } from "vue"

const $snackbar = useSnackbar()

const reverting = ref(true)

const library = ref<any[]>([])
const apps = ref<{ [id: string]: any }>({})

const pagination = reactive({
  app: 1
})

async function fetch() {
  try {
    reverting.value = true
    library.value = (await http.get("/api/library")).data
    for (const item of library.value) {
      if (apps.value[item.app_id] == null) {
        apps.value[item.app_id] = (await http.get(`/api/explore/apps/${item.app_id}`)).data
      }
    }

    reverting.value = false
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "red" })
  }
}

onMounted(() => {
  fetch()
})
</script>
