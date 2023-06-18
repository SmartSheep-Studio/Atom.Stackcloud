<template>
  <v-app id="atom">
    <initializing :show="connecting" />

    <router-view v-if="!connecting" />

    <snackbar />
  </v-app>
</template>

<script lang="ts" setup>
import { onMounted, ref, watch } from "vue"
import { useAccount } from "@/stores/account"
import { keepGate } from "@/utils/gatekeeper"
import { useRoute } from "vue-router"
import { useSnackbar } from "@/stores/snackbar"
import { useEndpoint } from "@/stores/endpoint"
import { getEndpointPath } from "@/utils/endpoint"
import Snackbar from "@/components/global/snackbar.vue"
import Initializing from "@/components/global/initializing.vue"

const connecting = ref(true)

const $route = useRoute()

const $account = useAccount()
const $endpoint = useEndpoint()

function fetch() {
  Promise.all([$endpoint.fetch(), $account.fetch()])
    .then(() => {
      connecting.value = false
    })
    .catch((e) => {
      const $snackbar = useSnackbar()
      $snackbar.show({
        text: "Initializing failed, cannot establish connection with server. Retry in 5s...",
        color: "error",
        timeout: 4850
      })
      setTimeout(() => fetch(), 5000)
    })
}

onMounted(() => {
  fetch()

  watch($account, () => {
    if (!keepGate($route)) {
          window.location.href = getEndpointPath('/auth/sign-in', `redirect_uri=${window.location.href}`)
    }
  })

  watch(
    $route,
    (v) => {
      if (!connecting.value) {
        if (!keepGate(v)) {
          window.location.href = getEndpointPath('/auth/sign-in', `redirect_uri=${window.location.href}`)
        }
      }
    },
    { deep: true }
  )
})
</script>

<style scoped>
.loading-placeholder {
  display: flex;
  justify-content: center;
  align-items: center;
}

.loading-inner-placeholder {
  width: 100%;
  height: fit-content;
  display: flex;
  flex-direction: row;
}
</style>