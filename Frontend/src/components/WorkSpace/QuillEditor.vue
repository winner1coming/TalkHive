<template>
  <div class="container">
    <!-- 顶部工具栏 -->
    <div class="editor-toolbar">
      <div class="toolbar-item">
        <label for="noteFilename">文件名：</label>
        <input type="text" v-model="currentNote.filename" placeholder="输入文件名" />
        <label for="noteCategory">分类：</label>
        <select v-model="currentNote.category">
          <option v-for="category in categories" :key="category" :value="category">{{ category }}</option>
        </select>
      </div>
      <div class="actions">
        <button class="btn" @click="saveContent">保存</button>
        <button class="btn" @click="cancelEdit">返回</button>
      </div>
    </div>

    <!-- TipTap 工具栏 -->
    <CollabToolbar :editor="editor" />

    <!-- TipTap 编辑器 -->
    <editor-content :editor="editor" class="tiptap-editor" />
  </div>
</template>

<script>
import { Editor, EditorContent } from "@tiptap/vue-3"
import * as WorkSpaceAPI from "@/services/workspace_api"
import { getLocalExtensions } from "@/utils/collab-schema"
import CollabToolbar from "./CollabToolbar.vue"

// Quill Delta → HTML 简易转换（兼容旧笔记格式）
function deltaToHtml(delta) {
  if (!delta || !delta.ops || !Array.isArray(delta.ops)) {
    return '<p></p>'
  }
  let html = ''
  for (const op of delta.ops) {
    const text = String(op.insert || '')
    const attrs = op.attributes || {}

    // 图片/嵌入内容
    if (op.insert && typeof op.insert === 'object') {
      if (op.insert.image) {
        html += `<img src="${op.insert.image}" />`
      }
      continue
    }

    // 换行 = 段落结束
    if (text === '\n') {
      html += '</p><p>'
      continue
    }

    let wrapped = text
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')

    if (attrs.bold) wrapped = `<strong>${wrapped}</strong>`
    if (attrs.italic) wrapped = `<em>${wrapped}</em>`
    if (attrs.underline) wrapped = `<u>${wrapped}</u>`
    if (attrs.strike) wrapped = `<s>${wrapped}</s>`
    if (attrs.code) wrapped = `<code>${wrapped}</code>`
    if (attrs.link) wrapped = `<a href="${attrs.link}">${wrapped}</a>`

    if (attrs.background) wrapped = `<span style="background:${attrs.background}">${wrapped}</span>`
    if (attrs.color) wrapped = `<span style="color:${attrs.color}">${wrapped}</span>`

    html += wrapped
  }
  return `<p>${html}</p>`
}

export default {
  name: "NoteEditor",
  components: { EditorContent, CollabToolbar },

  computed: {
    currentNote() {
      return this.$store.getters.getCurrentNote
    },
    categories() {
      const cats = this.$store.getters.getCategories
      cats.push('')
      return cats
    },
  },

  data() {
    return {
      editor: null,
    }
  },

  mounted() {
    // 创建 TipTap 编辑器（无 Yjs，本地模式）
    this.editor = new Editor({
      extensions: getLocalExtensions(),
      content: '',
    })

    // 加载笔记内容
    this.loadContent()
  },

  beforeUnmount() {
    if (this.editor) {
      this.editor.destroy()
      this.editor = null
    }
  },

  methods: {
    async loadContent() {
      try {
        const res = await WorkSpaceAPI.getNoteContent(this.currentNote.note_id)
        const content = res.data.content
        if (!content) return

        // 尝试解析：Quill Delta / TipTap JSON / HTML
        let html = ''
        try {
          const parsed = JSON.parse(content)
          if (parsed && parsed.ops) {
            // 旧格式：Quill Delta → HTML
            html = deltaToHtml(parsed)
          } else if (parsed && parsed.type) {
            // TipTap JSON 格式
            this.editor.commands.setContent(parsed)
            return
          } else {
            html = content
          }
        } catch {
          // 纯文本或 HTML
          html = content
        }

        if (html) {
          this.editor.commands.setContent(html)
        }
      } catch (e) {
        console.error("加载文档内容失败：", e)
      }
    },

    async saveContent() {
      // 以 HTML 格式保存（通用格式，兼容性强）
      const html = this.editor.getHTML()
      try {
        await WorkSpaceAPI.saveEditNote(
          this.currentNote.note_id,
          this.currentNote.filename,
          this.currentNote.category,
          html,
        )
        console.log("保存成功")
      } catch (e) {
        console.error("保存失败：", e)
      }
    },

    cancelEdit() {
      this.$router.push("/workspace/notes")
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
.editor-toolbar {
  display: flex;
  gap: 10px;
  justify-content: space-between;
  margin-bottom: 6px;
  align-items: center;
}
.toolbar-item {
  display: flex;
  align-items: center;
}
.toolbar-item label,
.toolbar-item input,
.toolbar-item select {
  margin-right: 5px;
  margin-left: 10px;
  font-size: 14px;
  color: var(--text-color);
}
.actions {
  display: flex;
  gap: 15px;
  margin-right: 10px;
}
.btn {
  cursor: pointer;
  padding: 2px 5px;
}
.tiptap-editor {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 20px;
  border: 1px solid #ddd;
  border-top: none;
}
.tiptap-editor :deep(.ProseMirror) {
  outline: none;
  min-height: 100%;
}
.tiptap-editor :deep(.ProseMirror p) {
  margin: 0.5em 0;
}
.tiptap-editor :deep(.ProseMirror h1),
.tiptap-editor :deep(.ProseMirror h2),
.tiptap-editor :deep(.ProseMirror h3) {
  margin: 0.8em 0 0.4em;
}
</style>
