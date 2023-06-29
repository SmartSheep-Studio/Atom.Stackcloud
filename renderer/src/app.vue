<template>
  <n-dialog-provider>
    <n-message-provider>
      <div class="w-full h-screen relative">
        <n-layout has-sider position="absolute">
          <n-layout-sider bordered collapse-mode="width" collapsed :collapsed-width="64" class="pt-2">
            <n-menu
              collapsed
              v-model:value="menuKey"
              :collapsed-width="64"
              :collapsed-icon-size="22"
              :options="menuOptions"
            />
          </n-layout-sider>

          <n-layout class="w-full h-full">
            <data-provider>
              <gatekeeper>
                <router-view />
              </gatekeeper>
            </data-provider>
          </n-layout>
        </n-layout>
      </div>
    </n-message-provider>
  </n-dialog-provider>
</template>

<script lang="ts" setup>
import DataProvider from "@/data-provider.vue"
import Gatekeeper from "@/components/global/gatekeeper.vue"
import { useEndpoint } from "@/stores/endpoint"
import { usePrincipal } from "@/stores/principal"
import { h, type Component, computed, type Ref, ref } from "vue"
import { RouterLink, useRoute, useRouter } from "vue-router"
import { NIcon, type MenuOption } from "naive-ui"
import { TerminalRound, ExploreRound } from "@vicons/material"

const $route = useRoute()
const $router = useRouter()
const $endpoint = useEndpoint()
const $principal = usePrincipal()

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuKey = ref($route.name)
const menuOptions: Ref<MenuOption[]> = computed(() =>
  $principal.isLoggedIn
    ? [
        {
          label: () => h(RouterLink, { to: { name: "landing" } }, { default: () => "Explore" }),
          icon: renderIcon(ExploreRound),
          key: "landing",
        },
        {
          label: () => h(RouterLink, { to: { name: "console" } }, { default: () => "Console" }),
          icon: renderIcon(TerminalRound),
          key: "console",
        },
      ]
    : []
)
</script>

<style>
.n-layout-header {
  height: 72px;
}

.n-layout-footer {
  padding: 24px;
}

.h-max {
  height: 100vh;
}

.w-dialog {
  width: 520px;
}
</style>
