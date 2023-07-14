import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/console",
      name: "console",
      component: () => import("@/views/landing.vue"),
    },
    {
      path: "/console/:app",
      name: "console.apps",
      component: () => import("@/views/app.vue"),
    },
    {
      path: "/console/apps/create",
      name: "console.apps.create",
      component: () => import("@/views/actions/create-app.vue"),
    },
  ],
})

export default router
