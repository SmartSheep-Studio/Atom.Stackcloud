<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-dialog-provider>
      <n-message-provider>
        <div class="w-full relative" :class="isUnderShadow ? 'h-max' : 'h-screen'">
          <n-layout position="absolute">
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
  </n-config-provider>
</template>

<script lang="ts" setup>
import DataProvider from "@/data-provider.vue";
import Gatekeeper from "@/components/global/gatekeeper.vue";
import { useEndpoint } from "@/stores/endpoint";
import { usePrincipal } from "@/stores/principal";
import { type Component, computed, h } from "vue";
import { useRoute, useRouter } from "vue-router";
import { NIcon } from "naive-ui";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const $route = useRoute();
const $router = useRouter();
const $endpoint = useEndpoint();
const $principal = usePrincipal();

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

const themeOverrides = {
  common: {
    primaryColor: "#ca4d4dFF",
    primaryColorHover: "#DF5656FF",
    primaryColorPressed: "#C04747FF",
    primaryColorSuppl: "#A84141FF"
  }
};

// Use for dynamic calculate height
const isUnderShadow = computed(() => {
  return (window as any).__POWERED_BY_WUJIE__ != null;
});
</script>

<style>
.n-layout-header {
  height: 72px;
}

.n-layout-footer {
  padding: 24px;
}

.h-max {
  height: calc(100vh - 72px);
}

.w-dialog {
  width: 520px;
}
</style>
