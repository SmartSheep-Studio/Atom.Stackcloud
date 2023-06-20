<template>
  <v-app-bar v-if="isBuiltin">
    <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>

    <v-toolbar-title>
      <span @click="$router.push({ name: 'explore' })" style="cursor: pointer">Matrix Marketplace</span>
    </v-toolbar-title>

    <v-spacer></v-spacer>
  </v-app-bar>

  <v-navigation-drawer v-if="isBuiltin" v-model="drawer" width="300">
    <template #prepend>
      <v-list-item
        exact
        v-if="$account.isLoggedIn"
        :title="$account.profile?.nickname"
        :subtitle="usePlaceholder('description', $account.profile?.description)"
        class="py-3"
      >
        <template #prepend>
          <v-avatar icon="mdi-account-circle" :image="usePlaceholder('avatar', $account.profile?.avatar_url)" />
        </template>
        <template #append>
          <v-btn icon="mdi-logout-variant" color="error" size="small" variant="text" @click="$account.logout()" flat />
          <v-btn
            icon="mdi-account-circle"
            color="primary"
            size="small"
            variant="text"
            :href="getEndpointPath('/users/profile')"
            flat
          />
        </template>
      </v-list-item>

      <v-list-item
        v-else
        class="py-3"
        title="Guest"
        subtitle="Forbidden"
        exact
        :href="getEndpointPath('/auth/sign-in', `redirect_uri=${fallback}`)"
      />
    </template>

    <v-divider></v-divider>

    <v-list density="compact" color="primary" nav>
      <v-list-group
        prepend-icon="mdi-console"
        title="Console"
        value="console"
        v-if="hasUserPermissions('matrix.console.view')"
      >
        <template v-slot:activator="{ props }">
          <v-list-item v-bind="props" />
        </template>

        <v-list-item exact :to="{ name: 'console.dashboard' }" title="Dashboard" append-icon="mdi-chart-arc" />
      </v-list-group>
    </v-list>
  </v-navigation-drawer>

  <v-main class="bg-grey-lighten-2">
    <router-view />
  </v-main>
</template>

<script lang="ts" setup>
import { useAccount } from "@/stores/account"
import { usePlaceholder } from "@/utils/placeholders"
import { getEndpointPath } from "@/utils/endpoint"
import { hasUserPermissions } from "@/utils/gatekeeper"
import { computed } from "vue"
import { ref } from "vue"

const drawer = ref(true)
const fallback = `${window.location.origin}/api/auth/request`

const $account = useAccount()

const isBuiltin = computed(() => {
  try {
    return window.self === window.top
  } catch (e) {
    return false
  }
})
</script>
