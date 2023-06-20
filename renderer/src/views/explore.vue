<template>
  <v-container fluid class="pa-12">
    <div>
      <div class="text-h6">Matrix Marketplace</div>
      <div class="text-subtitle-1">A open source apps & video games shop, but friendly for everyone.</div>
    </div>

    <div class="mt-6">
      <div class="text-subtitle-1">Explore</div>
      <div class="mt-2" style="margin-left: -8px; margin-right: -8px">
        <v-data-iterator :items="apps" :page="pagination.app">
          <template v-slot:default="{ items }">
            <v-row class="explore-queue">
              <template v-for="(item, i) in items" :key="i">
                <v-col :cols="12" :sm="6" :md="4">
                  <v-card class="explore-item pa-4">
                    <v-card-item>
                      <div class="text-h6">{{ item.raw.name }}</div>
                      <div style="margin-left: -4px">
                        <v-badge v-for="tag in item.raw.tags" inline color="primary" :content="tag" />
                      </div>
                    </v-card-item>
                    <v-card-item>
                      {{ item.raw.description }}
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
import { reactive } from "vue"
import { onMounted } from "vue"
import { ref } from "vue"

const apps = ref<any[]>([])

const pagination = reactive({
  app: 1
})

const $snackbar = useSnackbar()

async function fetch() {
  try {
    apps.value = (await http.get("/api/explore/apps")).data
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  }
}

onMounted(() => {
  fetch()
})
</script>

<style scoped>
.explore-item {
  cursor: pointer;
  box-shadow: 0px 3px 1px -2px var(--v-shadow-key-umbra-opacity, rgba(0, 0, 0, 0.2)),
    0px 2px 2px 0px var(--v-shadow-key-penumbra-opacity, rgba(0, 0, 0, 0.14)),
    0px 1px 5px 0px var(--v-shadow-key-penumbra-opacity, rgba(0, 0, 0, 0.12)) !important;
}

.explore-item:hover {
  box-shadow: 0px 2px 4px -1px var(--v-shadow-key-umbra-opacity, rgba(0, 0, 0, 0.2)),
    0px 4px 5px 0px var(--v-shadow-key-penumbra-opacity, rgba(0, 0, 0, 0.14)),
    0px 1px 10px 0px var(--v-shadow-key-penumbra-opacity, rgba(0, 0, 0, 0.12)) !important;
}
</style>
