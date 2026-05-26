<template>
  <div class="container">
    <!-- 文档信息栏 -->
    <div class="doc-info">
      <div class="left-infos">
        <img src="@/assets/icon/return.png" alt="返回" class="icon" @click="returnToWorkspace"/>
      </div>
      <div class="doc_name">
        {{ currentDoc.doc_name }}
        <span v-if="loadingLargeDoc" class="loading-indicator">加载大型文档中...</span>
      </div>
      <div class="user-info">
        <div class="other-users">正在编辑：
          <div class="other-user"
            v-for="[key, user] in remoteUsers"
            :key="key"
            :style="{ color: user.color }"
          >
            <img
              :src="user.userIcon || user.avatar"
              :alt="user.name"
              style="width: 18px; height: 18px; border-radius: 50%;"
            />
            {{ user.name }}
          </div>
        </div>
        <div class="me-user">
          <img
            :src="currentUser.avatar"
            :alt="currentUser.username"
            style="width: 18px; height: 18px; border-radius: 50%; margin-top:3px; margin-right:5px;"
          />
          <span>{{ currentUser.username }}</span>
        </div>
      </div>
    </div>

    <!-- TipTap 工具栏 -->
    <CollabToolbar :editor="editor" />

    <!-- TipTap 编辑器 -->
    <editor-content :editor="editor" class="tiptap-editor" />
  </div>
</template>

<script>
import * as Y from "yjs"
import { WebsocketProvider } from "y-websocket"
import { Editor, EditorContent } from "@tiptap/vue-3"
import pako from "pako"
import { Buffer } from "buffer"
import * as WorkSpaceAPI from "@/services/workspace_api"
import { getCollabExtensions } from "@/utils/collab-schema"
import { workerEpochCompact, terminateWorker } from "@/utils/collab-worker"
import { getCachedSnapshot, setCachedSnapshot, deleteCachedSnapshot } from "@/utils/collab-cache"
import CollabToolbar from "./CollabToolbar.vue"

export default {
  name: "CollabEditor",
  components: { EditorContent, CollabToolbar },

  computed: {
    currentDoc() {
      return this.$store.getters.getCurrentDoc
    },
    currentUser() {
      return this.$store.state.user
    },
  },

  data() {
    return {
      ydoc: null,
      provider: null,
      editor: null,
      remoteUsers: new Map(),
      snapshotTimer: null,
      loadingLargeDoc: false,
      // 版本计数器变更检测：每次编辑递增，保存时对比，跳过未变更的保存
      changeVersion: 0,
      lastSavedVersion: 0,
    }
  },

  async created() {
    // 1. 创建 Yjs 文档
    this.ydoc = new Y.Doc()

    // 2. 连接 WebSocket 协作服务器
    this.provider = new WebsocketProvider(
      `ws://localhost:1234?room=${this.currentDoc.doc_id}`,
      this.currentDoc.doc_id,
      this.ydoc
    )

    // 3. 设置 awareness（用户在线状态）
    this.setupAwareness()

    // 4. 先加载快照再创建编辑器（避免 ySyncPlugin observer 在 Y.XmlFragment
    //    尚未完整时被远程更新触发导致 crash）
    await this.loadSnapshot()
    this.initEditor()
  },

  mounted() {
    // 计算视口高度用于容器尺寸
    const vh = window.innerHeight * 0.01
    document.documentElement.style.setProperty("--vh", `${vh}px`)
  },

  beforeUnmount() {
    window.removeEventListener("beforeunload", this.destroyHandler)
    this.destroyHandler()
  },

  methods: {
    initEditor() {
      // 创建 TipTap 编辑器（此时 Y.XmlFragment 已被快照填充）
      this.editor = new Editor({
        extensions: getCollabExtensions(this.ydoc, this.provider, {
          id: this.currentUser.id,
          username: this.currentUser.username,
          avatar: this.currentUser.avatar,
        }),
      })

      // 跟踪本地变更（用于跳过未变更的保存）
      this.editor.on('update', () => { this.changeVersion++ })

      // 监听远程用户状态变化
      if (this.provider?.awareness) {
        this.provider.awareness.on("change", () => {
          const states = this.provider.awareness.getStates()
          const remoteUsers = new Map()
          states.forEach((state, clientID) => {
            const user = state.user
            if (user && user.id !== this.currentUser.id && !remoteUsers.has(user.id)) {
              remoteUsers.set(user.id, user)
            }
          })
          this.remoteUsers = remoteUsers
        })
      }

      // 每 30 秒保存快照
      this.snapshotTimer = setInterval(() => {
        this.saveSnapshot()
      }, 30000)

      window.addEventListener("beforeunload", this.destroyHandler)
    },

    setupAwareness() {
      const randomColor = () => `#${Math.floor(Math.random()*0xffffff).toString(16).padStart(6, '0')}`
      this.provider.awareness.setLocalStateField("user", {
        id: this.currentUser.id,
        name: this.currentUser.username,
        color: randomColor(),
        userIcon: this.currentUser.avatar,
        avatar: this.currentUser.avatar,
      })
    },

    // 保存快照：变更检测 + Web Worker Epoch 压缩（回退主线程）+ gzip 上传
    async saveSnapshot() {
      if (this.changeVersion === this.lastSavedVersion) return
      this.lastSavedVersion = this.changeVersion

      const fullState = Y.encodeStateAsUpdate(this.ydoc)

      // 优先在 Web Worker 中做 epoch 压缩（不阻塞主线程），失败则回退到主线程
      let compactState
      try {
        const buf = await workerEpochCompact(fullState)
        compactState = new Uint8Array(buf)
      } catch (err) {
        // Worker 不可用（不支持、加载失败等），在主线程执行
        const freshDoc = new Y.Doc()
        Y.applyUpdate(freshDoc, fullState)
        compactState = Y.encodeStateAsUpdate(freshDoc)
        freshDoc.destroy()
      }

      const compressed = pako.gzip(compactState)
      const base64 = Buffer.from(compressed).toString("base64")
      try {
        const response = await WorkSpaceAPI.saveSnapshot(this.currentDoc.doc_id, base64)
        if (response.status !== 200) {
          console.log("保存协作文档快照失败")
        } else {
          // 更新本地缓存，下次打开无需重新下载
          await setCachedSnapshot(this.currentDoc.doc_id, base64).catch(() => {})
        }
      } catch (err) {
        console.log("保存快照出错：", err)
      }
    },

    // 从后端加载快照，使用 IndexedDB 缓存避免重复下载和解压
    // 策略：stale-while-revalidate — 先应用缓存（即时），再异步刷新服务器版本
    async loadSnapshot() {
      // Phase 1: 从 IndexedDB 缓存加载（即时，无网络开销）
      const cached = await getCachedSnapshot(this.currentDoc.doc_id).catch(() => null)
      if (cached?.snapshot) {
        try {
          const compressed = Uint8Array.from(atob(cached.snapshot), c => c.charCodeAt(0))
          const update = pako.ungzip(compressed)
          // 大文档缓存也用 rAF 推迟，避免首次渲染卡顿
          const LARGE_DOC_THRESHOLD = 500000
          if (update.length > LARGE_DOC_THRESHOLD) {
            await new Promise(resolve => requestAnimationFrame(resolve))
          }
          Y.applyUpdate(this.ydoc, update)
          console.log("Snapshot restored from local cache")
        } catch (err) {
          console.log("Cache load failed, clearing corrupted cache:", err)
          // 清除损坏的缓存，避免下次再加载同一份坏数据
          deleteCachedSnapshot(this.currentDoc.doc_id).catch(() => {})
        }
      }

      // Phase 2: 总是从服务器获取最新版本（stale-while-revalidate）
      // 如果 Phase 1 已应用缓存，这次刷新在后台完成，用户无感知
      try {
        const response = await WorkSpaceAPI.getSnapshot(this.currentDoc.doc_id)
        const data = response.data
        if (!data.snapshot) {
          if (!cached) console.log("No snapshot found, starting fresh")
          return
        }

        // 将服务器响应写入缓存（下次打开直接走 Phase 1）
        await setCachedSnapshot(this.currentDoc.doc_id, data.snapshot).catch(() => {})

        let update
        try {
          const compressed = Uint8Array.from(atob(data.snapshot), c => c.charCodeAt(0))
          update = pako.ungzip(compressed)
        } catch {
          // 兼容旧格式（未压缩）
          update = Uint8Array.from(atob(data.snapshot), c => c.charCodeAt(0))
        }

        // 对大文档使用 rAF 推迟 Y.applyUpdate，让浏览器先渲染 UI 再处理重计算
        // 注意：Yjs 更新是自包含二进制格式，不可在任意字节边界切片
        const LARGE_DOC_THRESHOLD = 500000
        if (update.length > LARGE_DOC_THRESHOLD) {
          // 如果缓存已应用，用户已在编辑，不显示 loading 遮罩
          if (!cached) {
            this.loadingLargeDoc = true
          }
          await new Promise(resolve => requestAnimationFrame(resolve))
          Y.applyUpdate(this.ydoc, update)
          this.loadingLargeDoc = false
        } else {
          Y.applyUpdate(this.ydoc, update)
        }
        console.log("Snapshot restored from server")
      } catch (err) {
        console.log("Server snapshot decode failed, starting from current state:", err)
        this.loadingLargeDoc = false
      }
    },

    returnToWorkspace() {
      this.$router.push("/workspace/collabdocs")
    },

    destroyHandler() {
      if (this.provider?.awareness) {
        this.provider.awareness.setLocalState(null)
      }
      this.provider?.disconnect()

      if (this.editor) {
        this.editor.destroy()
        this.editor = null
      }
      if (this.snapshotTimer) {
        clearInterval(this.snapshotTimer)
        this.snapshotTimer = null
      }
      terminateWorker()
    },
  },
}
</script>

<style scoped>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  overflow: hidden;
  box-sizing: border-box;
}
* {
  box-sizing: inherit;
}
.container {
  display: flex;
  flex-direction: column;
  height: calc(var(--vh, 1vh) * 100 - 40px);
  width: 100%;
  padding: 20px;
}
.doc-info {
  display: flex;
  flex-direction: row;
  flex-shrink: 0;
  align-items: center;
  justify-content: space-between;
  background-color: lightgrey;
  padding: 4px 8px;
}
.left-infos {
  display: flex;
  flex-direction: row;
  align-items: center;
}
.doc_name {
  font-size: 20px;
  margin-left: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
}
.user-info {
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  align-items: center;
}
.other-users {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 6px;
  font-size: 14px;
}
.other-user, .me-user {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 4px 6px;
  font-size: 14px;
}
.me-user {
  margin-left: 10px;
}
.icon {
  width: 25px;
  height: 25px;
  margin-left: 5px;
  cursor: pointer;
}
.loading-indicator {
  font-size: 13px;
  color: #666;
  animation: pulse 1.5s ease-in-out infinite;
}
@keyframes pulse {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 1; }
}
.tiptap-editor {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 20px;
  border: 1px solid #ddd;
  border-top: none;
}
/* 让浏览器跳过离屏 DOM 的布局/绘制 */
.tiptap-editor :deep(.ProseMirror) {
  outline: none;
  min-height: 100%;
}
/* 为每个块级节点启用 content-visibility，浏览器可跳过离屏节点的布局/绘制 */
.tiptap-editor :deep(.ProseMirror > *) {
  content-visibility: auto;
  contain-intrinsic-size: auto 1.5rem;
}
.tiptap-editor :deep(.ProseMirror p) {
  margin: 0.5em 0;
}
.tiptap-editor :deep(.ProseMirror h1),
.tiptap-editor :deep(.ProseMirror h2),
.tiptap-editor :deep(.ProseMirror h3) {
  margin: 0.8em 0 0.4em;
}
/* 远程用户光标：默认 cursorBuilder 只设了 border-color，需补全 border 样式 */
.tiptap-editor :deep(.ProseMirror-yjs-cursor) {
  border-left: 2px solid;
  height: 1em;
  position: relative;
  content-visibility: visible;
}
/* 光标上的用户名标签 */
.tiptap-editor :deep(.ProseMirror-yjs-cursor > div) {
  position: absolute;
  top: -1.4em;
  left: -2px;
  font-size: 11px;
  color: #fff;
  padding: 1px 4px;
  white-space: nowrap;
  border-radius: 2px 2px 2px 0;
  pointer-events: none;
}
/* 远程用户选中区域半透明高亮 */
.tiptap-editor :deep(.ProseMirror-yjs-selection) {
  opacity: 0.3;
  content-visibility: visible;
}
</style>
