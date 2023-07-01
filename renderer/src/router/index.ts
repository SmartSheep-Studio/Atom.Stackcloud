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
      name: "console.apps",
      component: () => import("@/views/console/app-details.vue"),
    },
    {
      path: "/console/:app/posts/create",
      name: "console.apps.posts.create",
      component: () => import("@/views/console/actions/create-post.vue"),
    },
    {
      path: "/console/:app/posts/:post/update",
      name: "console.apps.posts.update",
      component: () => import("@/views/console/actions/update-post.vue"),
    },
    {
      path: "/console/apps/create",
      name: "console.apps.create",
      component: () => import("@/views/console/actions/create-app.vue"),
    },
  ],
})

export default router
