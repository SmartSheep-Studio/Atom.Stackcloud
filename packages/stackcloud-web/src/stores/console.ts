import { http } from "@/utils/http"
import { useMessage } from "naive-ui"
import { defineStore } from "pinia"
import { useRoute } from "vue-router"
import { ref } from "vue"

export const useConsole = defineStore("console", () => {
  const isPrepared = ref(false)
  const focusApp = ref<any>({})

  const $route = useRoute()
  const $message = useMessage()

  async function fetch() {
    try {
      isPrepared.value = false
      focusApp.value = (await http.get(`/api/apps/${$route.params.app}`)).data
      focusApp.value.cloud_functions = (await http.get(`/api/apps/${$route.params.app}/functions`)).data
      focusApp.value.cloud_collections = (await http.get(`/api/apps/${$route.params.app}/records`)).data
      isPrepared.value = true
    } catch (e: any) {
      $message.error(`Something went wrong... ${e}. Retry in 3 seconds...`)
      setTimeout(() => {
        fetch()
      }, 3000)
    }
  }  

  return { isPrepared, focusApp, fetch }
})
