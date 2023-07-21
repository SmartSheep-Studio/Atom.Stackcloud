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
import { computed } from "vue";

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

.h-max-inner {
  height: calc(100vh - 72px - 48px);
}

.h-screen-inner {
  height: calc(100vh - 48px);
}

.w-dialog {
  width: 520px;
}
</style>
