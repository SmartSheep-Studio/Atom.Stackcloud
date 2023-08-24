<template>
  <div>
    <n-card size="small" class="px-2 h-[48px] flex content-center rounded-none">
      <div class="flex items-center justify-between gap-2">
        <n-button text @click="$router.push({ name: 'console' })">
          <template #icon>
            <n-icon :component="ArrowBackRound" />
          </template>
        </n-button>
        <div class="text-gray-400 hover:text-gray-900 transition-colors cursor-default">
          {{ $console.focusApp.name }} - Stackcloud Console
        </div>
        <n-button text>
          <template #icon>
            <n-icon :component="SettingsRound" />
          </template>
        </n-button>
      </div>
    </n-card>
    <n-spin :show="!$console.isPrepared">
      <splitpanes class="h-screen-inner">
        <pane :min-size="15" :max-size="30">
          <div class="h-full p-4">
            <n-tree block-line expand-on-click :data="navNodes" :node-props="navProps" />
          </div>
        </pane>
        <pane>
          <div class="w-full h-full">
            <router-view :key="$route.fullPath" />
          </div>
        </pane>
      </splitpanes>
    </n-spin>
  </div>
</template>

<script lang="ts" setup>
import "splitpanes/dist/splitpanes.css"
import { Splitpanes, Pane } from "splitpanes"
import { type TreeOption, NIcon, NSpace, NTag } from "naive-ui"
import { onMounted, computed, h } from "vue"
import { AppsFilled, FunctionsRound, LayersRound, ArrowBackRound, SettingsRound, AddBoxRound } from "@vicons/material"
import { useRoute, useRouter } from "vue-router"
import { useConsole } from "@/stores/console"

const $route = useRoute()
const $router = useRouter()
const $console = useConsole()

const navNodes = computed(() => {
  return [
    {
      label: "Application",
      key: "application",
      component: { name: "console.apps.settings", params: { app: $console.focusApp.slug } },
      prefix: () => h(NIcon, null, { default: () => h(AppsFilled) }),
    },
    {
      label: "Cloud Functions",
      key: "cloud-functions",
      prefix: () => h(NIcon, null, { default: () => h(FunctionsRound) }),
      suffix: () =>
        h(
          NIcon,
          {
            onClick: (e: Event) => {
              e.preventDefault()
              e.stopPropagation()
              $router.push({ name: "console.apps.functions.create", params: { app: $console.focusApp.slug } })
            },
          },
          {
            default: () => h(AddBoxRound),
          }
        ),
      children: $console.focusApp.cloud_functions?.map((v: any) => ({
        label: v.name,
        key: `cloud-functions-${v.slug}`,
        component: { name: "console.apps.functions", params: { app: $console.focusApp.slug, function: v.slug } },
        suffix: () =>
          h(
            NSpace,
            { size: "small" },
            v.tags.map((tag: string) => h(NTag, { size: "small", bordered: false, type: "primary" }, { default: tag }))
          ),
      })),
    },
    {
      label: "Cloud Collections",
      key: "cloud-collections",
      prefix: () => h(NIcon, null, { default: () => h(LayersRound) }),
      suffix: () =>
        h(
          NIcon,
          {
            onClick: (e: Event) => {
              e.preventDefault()
              e.stopPropagation()
              $router.push({ name: "console.apps.collections.create", params: { app: $console.focusApp.slug } })
            },
          },
          {
            default: () => h(AddBoxRound),
          }
        ),
      children: $console.focusApp.cloud_collections?.map((v: any) => ({
        label: v.name,
        key: `cloud-collections-${v.slug}`,
        component: { name: "console.apps.collections", params: { app: $console.focusApp.slug, collection: v.slug } },
        suffix: () =>
          h(
            NSpace,
            { size: "small" },
            v.tags.map((tag: string) => h(NTag, { size: "small", bordered: false, type: "primary" }, { default: tag }))
          ),
      })),
    },
  ]
})
const navProps = ({ option }: { option: TreeOption }) => {
  return {
    onClick() {
      if (option.component) {
        $router.push(option.component)
      }
    },
  }
}

onMounted(() => {
  $console.fetch()
})

// Use for dynamic calculate height
const isUnderShadow = computed(() => {
  return (window as any).__POWERED_BY_WUJIE__ != null;
});
</script>

<style>
.h-screen-inner {
  height: calc(100vh - 48px);
}

.splitpanes__pane {
  background-color: #ffffff;
}

.splitpanes__splitter {
  min-width: unset !important;
  width: 0.667px;
  background-color: #dedede;
  position: relative;
}

.splitpanes__splitter:before {
  content: "";
  position: absolute;
  left: 0;
  top: 0;
  transition: opacity 0.4s;
  background-color: rgba(255, 0, 0, 0.3);
  opacity: 0;
  z-index: 1;
}

.splitpanes__splitter:hover:before {
  opacity: 1;
}

.splitpanes--vertical > .splitpanes__splitter:before {
  left: -15px;
  right: -15px;
  height: 100%;
}

.splitpanes--horizontal > .splitpanes__splitter:before {
  top: -15px;
  bottom: -15px;
  width: 100%;
}
</style>
