<template>
  <v-container fluid class="pa-md-12">
    <div>
      <div class="text-h6">{{ app.name }}</div>
      <div class="text-subtitle-1">{{ app.description }}</div>
    </div>

    <div class="mt-4">
      <v-row dense>
        <v-col :md="8" :cols="12">
          <v-card>
            <v-md-editor :model-value="app.details" mode="preview" />
          </v-card>
        </v-col>
        <v-col :md="4" :cols="12">
          <v-card title="Buy Now" prepend-icon="mdi-shopping">
            <v-card-item class="px-6 pb-6">
              <div class="text-subtitle-2 text-mono">{{ product.name }}</div>
              <div class="text-h6 text-mono">{{ product.unit_price }} {{ product.price_unit }}</div>
              <v-btn
                class="mt-2"
                block
                color="primary"
                prepend-icon="mdi-login-variant"
                @click="checkout"
                :loading="submitting"
              >
                Go checkout
              </v-btn>
            </v-card-item>
          </v-card>
        </v-col>
      </v-row>
    </div>

    <v-footer class="border-t" app>
      <v-breadcrumbs density="compact">
        <v-breadcrumbs-item exact :to="{ name: 'explore' }">Explore</v-breadcrumbs-item>
        <v-breadcrumbs-divider />
        <v-breadcrumbs-item disabled>{{ app.name }}</v-breadcrumbs-item>
      </v-breadcrumbs>
    </v-footer>
  </v-container>
</template>

<script lang="ts" setup>
import { useSnackbar } from "@/stores/snackbar"
import { http } from "@/utils/http"
import { onMounted } from "vue"
import { ref } from "vue"
import { useRoute } from "vue-router"

const $route = useRoute()
const $snackbar = useSnackbar()

const submitting = ref(false)

const app = ref<any>({})
const product = ref<any>({})

async function fetch() {
  try {
    app.value = (await http.get(`/api/explore/apps/${$route.params.app}`)).data
    product.value = (await http.get(`/api/explore/apps/${$route.params.app}/product`)).data
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  }
}

async function checkout() {
  try {
    submitting.value = true
    const res = await http.post("/api/library/add", { app: $route.params.app })

    window.location.href = res.data.url
  } catch (e: any) {
    $snackbar.show({ text: `Something wrong... ${e}`, color: "error" })
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetch()
})
</script>
