import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "landing",
      component: () => import("@/views/landing.vue"),
    },

    {
      path: "/console",
      name: "console",
      component: () => import("@/views/console/landing.vue"),
    }
  ],
})

export default router
