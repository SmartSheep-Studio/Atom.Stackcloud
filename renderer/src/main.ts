import "./assets/style.css"

import { createApp } from "vue"
import { createPinia } from "pinia"

import wrapper from "./app.vue"
import router from "./router"

import "vfonts/Lato.css"
import "vfonts/FiraCode.css"

const app = createApp(wrapper)

app.use(createPinia())
app.use(router)

app.mount("#app")
