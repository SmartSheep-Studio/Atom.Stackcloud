import { defineStore } from "pinia"
import { useCookies } from "@vueuse/integrations/useCookies"
import { http } from "@/utils/http"
import { computed, ref } from "vue"
import { useLocalStorage } from "@vueuse/core"

export interface User {
  avatar_url?: string
  banner_url?: string
  contacts: Contact[]
  created_at?: Date
  deleted_at?: Date
  description?: string
  friends: User[]
  group_id?: number
  id?: number
  locked_at?: Date
  nickname?: string
  passcodes?: any[]
  permissions?: string[]
  sessions: Session[]
  updated_at?: Date
  user_assets?: UserAsset[]
  username?: string
}

export interface Contact {
  user_id?: number
  content?: string
  created_at?: Date
  deleted_at?: null
  description?: string
  id?: number
  name?: string
  type?: string
  updated_at?: Date
  verified_at?: null
}

export interface Session {
  access?: string
  available?: boolean
  client_id?: number
  code?: string
  created_at?: Date
  deleted_at?: Date
  expired_at?: Date
  id?: number
  ip?: string
  location?: string
  refresh?: string
  scope?: string[]
  type?: number
  updated_at?: Date
  user_id?: number
}

export interface UserAsset {
  id: number
  created_at: Date
  updated_at: Date
  deleted_at: null
  name: string
  size: number
  type: number
  storage_id: string
  storage_policy: null
  user_id: number
}

export const useAccount = defineStore("account", () => {
  const isLoggedIn = ref(false)
  const cookies = useCookies(["lineup_authorization"])
  const token = computed(() => cookies.get("lineup_authorization"))

  const profile = useLocalStorage<User | null>("atom-profile", null, {
    deep: true,
    listenToStorageChanges: true,
    serializer: {
      read(v) {
        try {
          return JSON.parse(v)
        } catch {
          return null
        }
      },
      write(v) {
        if (v != null) {
          return JSON.stringify(v)
        } else {
          return "null"
        }
      }
    }
  })

  async function fetch() {
    if (cookies.get("lineup_authorization") != null) {
      try {
        profile.value = (await http.get("/api/auth")).data
        isLoggedIn.value = true
      } catch {
        profile.value = null
        isLoggedIn.value = false
      }
    }
  }

  function logout() {
    cookies.remove("lineup_authorization")
    profile.value = null
    window.location.reload()
  }

  return { profile, isLoggedIn, token, fetch, logout }
})
