import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "layouts.main",
      component: () => import("@/layouts/main-layout.vue"),
      children: [
        {
          path: "/",
          name: "explore",
          component: () => import("@/views/explore.vue")
        },

        {
          path: "/console",
          name: "console.dashboard",
          component: () => import("@/views/console/dashboard.vue")
        },
        {
          path: "/console/:app",
          name: "console.apps",
          component: () => import("@/views/console/apps/details.vue")
        },
      ]
    }
  ]
})

export default router
