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
      path: "/apps/:app",
      name: "console.apps",
      component: () => import("@/views/app.vue"),
    },
    {
      path: "/apps/:app/collections/create",
      name: "console.apps.collections.create",
      component: () => import("@/views/actions/create-collection.vue"),
    },
    {
      path: "/apps/:app/collections/:collection",
      name: "console.apps.collections",
      component: () => import("@/views/collection.vue"),
    },
    {
      path: "/apps/:app/collections/:collection/update",
      name: "console.apps.collections.update",
      component: () => import("@/views/actions/update-collection.vue"),
    },
    {
      path: "/apps/create",
      name: "console.apps.create",
      component: () => import("@/views/actions/create-app.vue"),
    },
  ],
})

export default router
