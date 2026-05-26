<template>
  <div class="editor-toolbar" v-if="editor">
    <!-- 文本格式化 -->
    <button
      @click="editor.chain().focus().toggleBold().run()"
      :class="{ 'is-active': editor.isActive('bold') }"
      title="粗体"
    ><b>B</b></button>
    <button
      @click="editor.chain().focus().toggleItalic().run()"
      :class="{ 'is-active': editor.isActive('italic') }"
      title="斜体"
    ><i>I</i></button>
    <button
      @click="editor.chain().focus().toggleUnderline().run()"
      :class="{ 'is-active': editor.isActive('underline') }"
      title="下划线"
    ><u>U</u></button>
    <button
      @click="editor.chain().focus().toggleStrike().run()"
      :class="{ 'is-active': editor.isActive('strike') }"
      title="删除线"
    ><s>S</s></button>

    <span class="separator"></span>

    <!-- 标题 -->
    <button
      @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
      :class="{ 'is-active': editor.isActive('heading', { level: 1 }) }"
      title="标题 1"
    >H1</button>
    <button
      @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
      :class="{ 'is-active': editor.isActive('heading', { level: 2 }) }"
      title="标题 2"
    >H2</button>
    <button
      @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
      :class="{ 'is-active': editor.isActive('heading', { level: 3 }) }"
      title="标题 3"
    >H3</button>

    <span class="separator"></span>

    <!-- 列表 -->
    <button
      @click="editor.chain().focus().toggleBulletList().run()"
      :class="{ 'is-active': editor.isActive('bulletList') }"
      title="无序列表"
    >•</button>
    <button
      @click="editor.chain().focus().toggleOrderedList().run()"
      :class="{ 'is-active': editor.isActive('orderedList') }"
      title="有序列表"
    >1.</button>
    <button
      @click="editor.chain().focus().toggleTaskList().run()"
      :class="{ 'is-active': editor.isActive('taskList') }"
      title="任务列表"
    >☑</button>

    <span class="separator"></span>

    <!-- 引用与代码 -->
    <button
      @click="editor.chain().focus().toggleBlockquote().run()"
      :class="{ 'is-active': editor.isActive('blockquote') }"
      title="引用"
    >"</button>
    <button
      @click="editor.chain().focus().toggleCodeBlock().run()"
      :class="{ 'is-active': editor.isActive('codeBlock') }"
      title="代码块"
    >&lt;/&gt;</button>
    <button
      @click="editor.chain().focus().toggleCode().run()"
      :class="{ 'is-active': editor.isActive('code') }"
      title="行内代码"
    >&lt;c&gt;</button>

    <span class="separator"></span>

    <!-- 对齐 -->
    <button
      @click="editor.chain().focus().setTextAlign('left').run()"
      :class="{ 'is-active': editor.isActive({ textAlign: 'left' }) }"
      title="左对齐"
    >≡</button>
    <button
      @click="editor.chain().focus().setTextAlign('center').run()"
      :class="{ 'is-active': editor.isActive({ textAlign: 'center' }) }"
      title="居中对齐"
    >≡</button>
    <button
      @click="editor.chain().focus().setTextAlign('right').run()"
      :class="{ 'is-active': editor.isActive({ textAlign: 'right' }) }"
      title="右对齐"
    >≡</button>

    <span class="separator"></span>

    <!-- 链接与图片 -->
    <button @click="addLink" title="链接">🔗</button>
    <button @click="addImage" title="图片">🖼</button>

    <span class="separator"></span>

    <!-- 清除格式 -->
    <button @click="editor.chain().focus().clearNodes().unsetAllMarks().run()" title="清除格式">✕</button>

    <!-- 撤销/重做 -->
    <span class="separator"></span>
    <button @click="editor.chain().focus().undo().run()" title="撤销" :disabled="!editor.can().undo()">↶</button>
    <button @click="editor.chain().focus().redo().run()" title="重做" :disabled="!editor.can().redo()">↷</button>
  </div>
</template>

<script>
export default {
  name: "CollabToolbar",
  props: {
    editor: { type: Object, default: null },
  },
  methods: {
    addLink() {
      const previousUrl = this.editor.getAttributes('link').href
      const url = window.prompt('输入链接地址', previousUrl)
      if (url === null) return
      if (url === '') {
        this.editor.chain().focus().extendMarkRange('link').unsetLink().run()
        return
      }
      this.editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
    },
    addImage() {
      const url = window.prompt('输入图片地址')
      if (url) {
        this.editor.chain().focus().setImage({ src: url }).run()
      }
    },
  },
}
</script>

<style scoped>
.editor-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
  padding: 6px 8px;
  border-bottom: 1px solid #ddd;
  background: #f8f9fa;
  align-items: center;
}
.editor-toolbar button {
  border: 1px solid transparent;
  background: transparent;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 3px;
  font-size: 14px;
  line-height: 1;
  min-width: 28px;
  text-align: center;
}
.editor-toolbar button:hover {
  background: #e9ecef;
}
.editor-toolbar button.is-active {
  background: #cce5ff;
  border-color: #b8daff;
}
.editor-toolbar button:disabled {
  opacity: 0.4;
  cursor: default;
}
.separator {
  width: 1px;
  height: 20px;
  background: #ddd;
  margin: 0 4px;
  display: inline-block;
}
</style>
