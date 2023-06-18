import { createApp } from "vue"
import { createPinia } from "pinia"

import wrapper from "./app.vue"
import router from "./router"

import "vuetify/styles"
import { createVuetify } from "vuetify"
import * as components from "vuetify/components"
import * as directives from "vuetify/directives"
import { aliases, mdi } from "vuetify/iconsets/mdi"

import "@mdi/font/css/materialdesignicons.css"

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: "light",
    themes: {
      light: {
        colors: {
          background: "#ffffff",
          surface: "#ffffff",
          primary: "#3f51b5",
          secondary: "#2196f3",
          accent: "#673ab7",
          error: "#f44336",
          info: "#03a9f4",
          success: "#4caf50",
          warning: "#ff9800"
        }
      }
    }
  },
  icons: {
    defaultSet: "mdi",
    aliases,
    sets: {
      mdi
    }
  }
})

const app = createApp(wrapper)

app.use(createPinia())
app.use(vuetify)
app.use(router)

app.mount("#app")
