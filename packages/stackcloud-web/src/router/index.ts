import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "console",
      component: () => import("@/views/landing.vue"),
    },
    {
      path: "/apps/-/:app",
      name: "console.apps",
      component: () => import("@/views/app.vue"),
      children: [
        {
          path: "/apps/-/:app",
          name: "console.apps.intro",
          component: () => import("@/views/parts/intro.vue"),
        },
        {
          path: "/apps/-/:app/settings",
          name: "console.apps.settings",
          component: () => import("@/views/parts/settings.vue"),
        },
        {
          path: "/apps/-/:app/functions/create",
          name: "console.apps.functions.create",
          component: () => import("@/views/actions/create-function.vue"),
        },
        {
          path: "/apps/-/:app/functions/:function",
          name: "console.apps.functions",
          component: () => import("@/views/function.vue"),
        },
        {
          path: "/apps/-/:app/functions/:function/update",
          name: "console.apps.functions.update",
          component: () => import("@/views/actions/update-function.vue"),
        },
        {
          path: "/apps/-/:app/collections/create",
          name: "console.apps.collections.create",
          component: () => import("@/views/actions/create-collection.vue"),
        },
        {
          path: "/apps/-/:app/collections/:collection",
          name: "console.apps.collections",
          component: () => import("@/views/collection.vue"),
        },
        {
          path: "/apps/-/:app/collections/:collection/update",
          name: "console.apps.collections.update",
          component: () => import("@/views/actions/update-collection.vue"),
        },
      ],
    },
    {
      path: "/apps/create",
      name: "console.apps.create",
      component: () => import("@/views/actions/create-app.vue"),
    },
  ],
})

export default router
