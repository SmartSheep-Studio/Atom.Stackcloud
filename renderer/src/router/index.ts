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
    },
    {
      path: "/console/:app",
      name: "console.app",
      component: () => import("@/views/console/app-details.vue"),
    },
    {
      path: "/console/app/create",
      name: "console.app.create",
      component: () => import("@/views/console/actions/create-app.vue"),
    },
  ],
})

export default router
