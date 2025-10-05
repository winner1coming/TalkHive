<template>
  <div class="container">
    <!-- 工具栏 -->
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
    <!--  Quill 编辑器 -->
    <div ref="quillEditor" class="quill-editor"></div>

  </div>
</template>
<script>
import Quill from "quill";
import 'quill/dist/quill.snow.css';
import * as WorkSpaceAPI from '@/services/workspace_api';

let quill = null;
export default {
  name: "QuillEditor",

  computed: {
    currentNote() {
      return this.$store.getters.getCurrentNote;
    },
    categories() {
      this.localCategory = this.$store.getters.getCategories;
      this.localCategory.push('');
      return this.localCategory; 
    },
  },

  // data() {
  //   return {
  //     quill: null
  //   }
  // },

  mounted() {
    const toolbarOptions = [
      ['bold', 'italic', 'underline', 'strike'],
      ['blockquote', 'code-block'],
      ['link', 'image', 'video', 'formula'],
      [{ 'header': 1 }, { 'header': 2 }],
      [{ 'list': 'ordered'}, { 'list': 'bullet' }, { 'list': 'check' }],
      [{ 'script': 'sub'}, { 'script': 'super' }],
      [{ 'indent': '-1'}, { 'indent': '+1' }],
      [{ 'direction': 'rtl' }],
      [{ 'size': ['small', false, 'large', 'huge'] }],
      [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
      [{ 'color': [] }, { 'background': [] }],
      [{ 'font': [] }],
      [{ 'align': [] }],
      ['clean']
    ];

    quill = new Quill(this.$refs.quillEditor, {
      theme: "snow",
      modules: { toolbar: toolbarOptions }
    });

    // 初始化时加载内容
    this.loadContent();
  },
  beforeRouteLeave(){
    console.log("Quill Editor：触发 beforeRouteLeave");
  },
  beforeUnmount(){
    //this.destroyQuill();
    console.log("QuillEditor 触发 beforeUnmount");
  },

  methods: {
    async loadContent() {
      // 从后端取内容
      try {
        const res = await WorkSpaceAPI.getNoteContent(this.currentNote.note_id);
        const content = res.data.content;   // 假设返回的是JSON字符串
        if (content) {
          quill.setContents(JSON.parse(content));
        }
      } catch (e) {
        console.error("加载文档内容失败：", e);
      }
    },

    async saveContent() {
      // 获取 Delta 格式
      const delta = quill.getContents();
      try {
        await WorkSpaceAPI.saveEditNote(
          this.currentNote.note_id,
          this.currentNote.filename,
          this.currentNote.category,
          JSON.stringify(delta)  // 存字符串
        );
        console.log("保存成功");
      } catch (e) {
        console.error("保存失败：", e);
      }
    },

    // 清理quill实例
    destroyQuill() {
        if (quill) {
        // 移除编辑器 DOM
          quill.off && quill.off(); // 某些 Quill 版本支持
          quill = null;
          const editor = this.$refs.quillEditor;
          if (editor) {
            editor.innerHTML = '';  // 清空 DOM
          }
        }
    },

    // 取消编辑
    cancelEdit() {
      //this.destroyQuill();
      this.$router.push("/workspace/notes");
    },
  }
}
</script>

<style scoped>
  html, body{
    margin: 0;
    padding: 0;
    height: 100%;
    width: 100%;
    overflow: hidden;   /* 防止全局滚动 */
    box-sizing: border-box;
  }
  *{
    box-sizing: inherit;
  }
  .container{
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
  .toolbar-item select{
    margin-right: 5px;
    margin-left: 10px;
    font-size: 14px;
    color: var(--text-color);
  }

  .actions {
    display: flex;
    gap: 15px;
    margin-right:10px;
  }
  .btn {
    cursor: pointer;
    padding-top: 2px;
    padding-bottom: 2px;
    padding-left: 5px;
    padding-right:5px;
  }
  .ql-undo, .ql-redo {
    font-size: 18px;
    padding: 4px 8px;
    border: none;
    background: none;
    cursor: pointer;
  }
  .quill-editor{
    flex: 1;
    min-height: 0;
    padding: 20px;
  }

</style>
