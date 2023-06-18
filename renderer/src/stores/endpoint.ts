import { defineStore } from "pinia"
import { http } from "@/utils/http"
import { ref } from "vue"

export const useEndpoint = defineStore("endpoint", () => {
  const isPrepared = ref(true)
  const configuration = ref<any>({})
  const additional = ref<any>({})
  const service = ref<any>({})

  async function fetch() {
    try {
      const res = await http.get("/api/info")
      configuration.value = res.data.configuration
      additional.value = res.data.additional
      service.value = res.data.service
    } catch (e: any) {
      if (e.response.status === 503) {
        isPrepared.value = false
      } else {
        isPrepared.value = true
      }
      throw e
    }

    document.title = `${configuration.value.general.name ?? "Project Atom"} Lineup`
  }

  return { isPrepared, configuration, additional, service, fetch }
})
