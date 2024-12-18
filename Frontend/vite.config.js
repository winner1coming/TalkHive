import { defineConfig, loadEnv } from 'vite'
import { fileURLToPath, URL } from 'node:url'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import compressPlugin from 'vite-plugin-compression'
import MonacoEditorPlugin from 'vite-plugin-monaco-editor'

export default defineConfig({
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      },
      extensions: ['.js', '.json', 'jsx', '.vue', '.ts'] // 使用路径别名时想要省略的后缀名，可以自己 增减
    },
    assetsInclude: ['./src/assets'],
    plugins: [
      vue(),
      vueJsx({}),
      compressPlugin(),
      MonacoEditorPlugin({
        languages: ['javascript', 'typescript', 'html', 'css', 'json', 'markdown'],
        features: ['!gotoSymbol']  // 可以禁用某些功能来减少包大小
      })
    ],
    // define: {
    //   __APP_ENV__: env.APP_ENV
    // },
})
