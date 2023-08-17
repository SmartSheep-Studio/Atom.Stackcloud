import {fileURLToPath, URL} from "node:url"

import {defineConfig} from "vite"
import vue from "@vitejs/plugin-vue"
import AutoImport from "unplugin-auto-import/vite"
import Components from "unplugin-vue-components/vite"
import {NaiveUiResolver} from "unplugin-vue-components/resolvers"
import VueI18nPlugin from "@intlify/unplugin-vue-i18n/vite";

// https://vitejs.dev/config/
export default defineConfig({
    base: "/srv/subapps/stackcloud",
    plugins: [
        vue(),
        AutoImport({
            imports: [
                "vue",
                {
                    "naive-ui": ["useDialog", "useMessage", "useNotification", "useLoadingBar"],
                },
            ],
        }),
        Components({
            resolvers: [NaiveUiResolver()],
        }),
        VueI18nPlugin({runtimeOnly: false}),
    ],
    resolve: {
        alias: {
            "@": fileURLToPath(new URL("./src", import.meta.url)),
        },
    },
    server: {
        port: 5174,
        proxy: {
            "/srv/subapps/stackcloud/api": {
                target: "http://127.0.0.1:9444",
                changeOrigin: true,
                rewrite: (path) => path.replace("/srv/subapps/stackcloud", ""),
            },
        },
    },
})
