<template>
  <div class="note-editor">
    <!-- 工具栏 -->
    <div class="editor-toolbar">
      <div class="toolbar-item">
        <label for="noteFilename">文件名：</label>
        <input type="text" v-model="currentNote.filename" placeholder="输入文件名" />
      </div>

      <div class="toolbar-item">
        <label for="noteCategory">分类：</label>
        <select v-model="currentNote.category">
          <option v-for="category in categories" :key="category" :value="category">{{ category }}</option>
        </select>
      </div>

      <div class="toolbar-item">
        <label for="language">语言：</label>
        <select v-model="selectedLanguage" @change="changeLanguage">
          <option value="markdown">Markdown</option>
        </select>
      </div>

      <div class="toolbar-item">
        <label for="fontSize">字体：</label>
        <select v-model="fontSize" @change="changeFontSize">
          <option value="12">12px</option>
          <option value="14">14px</option>
          <option value="16">16px</option>
          <option value="18">18px</option>
          <option value="20">20px</option>
        </select>
      </div>

      <div class="toolbar-item">
        <label for="bgColor">背景：</label>
        <input type="color" v-model="backgroundColor" @input="changeBackgroundColor">
      </div>

      <div class="actions">
        <button @click="saveNote" style="padding:5px">保存</button>
        <button @click="cancelEdit">取消</button>
      </div>
    </div>

    <!-- Monaco 编辑器 -->
    <div ref="monacoEditor" class="monaco-editor-container"></div>
  </div>
</template>

<script>
import * as monaco from 'monaco-editor';
import { toRaw, nextTick } from 'vue';
import * as WorkSpaceAPI from '@/services/workspace_api';

export default {
  name: "NoteEditor",
  computed: {
    currentNote() {
      return this.$store.getters.getCurrentNote;  // 从 Vuex 获取 currentNote
    },
    categories() {
      this.localCategory = this.$store.getters.getCategories;
      this.localCategory.push('');
      return this.localCategory; 
    }
  },
  data() {
    return {
      noteContent: "",
      selectedLanguage: 'markdown',
      fontSize: '14',
      backgroundColor: '#ffffff',
      editor: null,
      // categories: ['工作', '学习', '个人', '项目'],  // 模拟从后端获取分类
      // 用于处理 v-model 的本地数据
      localFilename: this.note_filename,
      localCategory: this.note_category
    };
  },
  // created() {
  //   // 模拟从后端获取笔记数据
  //   const notes = [
  //     { id: 1, filename: "Vue学习笔记.md", content: "这是 Vue 学习笔记的内容...", category: '学习' },
  //     { id: 2, filename: "项目需求分析.docx", content: "这是项目需求分析的内容...", category: '工作' },
  //     { id: 3, filename: "代码优化方案.txt", content: "这是代码优化方案的内容...", category: '项目' },
  //   ];

  //   // 查找当前笔记，根据传递的 note_id
  //   this.note = notes.find(note => note.id === Number(this.note_id));
  //   if (this.note) {
  //     this.noteContent = this.note.content;
  //   } else {
  //     alert("笔记未找到！");
  //     this.$router.push("/workspace/notes");
  //   }
  // },
  mounted() {
    if (this.currentNote) {
      console.log("currentNote:", this.currentNote);
    } else {
      console.log("currentNote is not loaded yet.");
    }
    nextTick(async() => {
      if (this.$refs.monacoEditor) {
        await this.getNoteContent(this.currentNote.note_id);
        this.editor = monaco.editor.create(this.$refs.monacoEditor, {
          value: this.noteContent,
          language: this.selectedLanguage,
          automaticLayout: true,
          fontSize: parseInt(this.fontSize),
          theme: 'vs',
          backgroundColor: this.backgroundColor,
        });
        // // 在内容变化时更新 noteContent
        // this.editor.onDidChangeModelContent(() => {
        //   this.noteContent = this.editor.getValue(); // 获取 Monaco 编辑器的当前值
        // });
      }
    });
  },
  watch: {
    selectedLanguage(newLang) {
      if (this.editor) {
        this.editor.setModelLanguage(this.editor.getModel(), newLang);
      }
    },
    fontSize(newFontSize) {
      if (this.editor) {
        this.editor.updateOptions({ fontSize: parseInt(newFontSize) });
      }
    },
    backgroundColor(newColor) {
      if (this.editor) {
        const editorContainer = this.$refs.monacoEditor;
        editorContainer.style.backgroundColor = newColor;
      }
    },
  },
  methods: {
    // // 更新笔记文件名
    // updateNoteFilename() {
    //   this.$emit('update:note_filename', this.localFilename); // 通过事件通知父组件更新
    // },

    // updateNoteCategory() {
    //   // 更新分类
    //   this.$store.dispatch('updateCategory', this.currentNote.category);
    // },
    async getNoteContent(noteID) {
      const response = await WorkSpaceAPI.getNoteContent(noteID);
      const data = response.data;
      this.noteContent = data;
    },


    // 语言切换
    changeLanguage() {
      if (this.editor) {
        const model = toRaw(this.editor.getModel());
        if (model) {
          monaco.editor.setModelLanguage(model, this.selectedLanguage);
        }
      }
    },

    // 字体大小调整
    changeFontSize() {
      if (this.editor) {
        this.editor.updateOptions({ fontSize: parseInt(this.fontSize) });
      }
    },

    // 背景颜色调整
    changeBackgroundColor() {
      if (this.editor) {
        const editorContainer = this.$refs.monacoEditor;
        editorContainer.style.backgroundColor = this.backgroundColor;
      }
    },

    // 保存笔记
    async saveNote() {
      try {
        //const response = await WorkSpaceAPI.EditNote(this.newFile.filename + this.newFile.filetype, this.newFile.category);
        // console.log(this.noteContent);
        // this.noteContent = this.editor.getValue();
        // console.log(this.noteContent);
        const model = toRaw(this.editor.getModel());
        const content = model.getValue();
        console.log(content);
        const response = await WorkSpaceAPI.saveEditNote(this.currentNote.note_id, this.currentNote.filename, this.currentNote.category,content);
        if (response.status != 200) {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('无法保存文件:', error);
        alert('保存文件失败！');
      }
      this.$router.push("/workspace/notes");
    },

    // 取消编辑
    cancelEdit() {
      this.$router.push("/workspace/notes");
    },
  },
};
</script>

<style scoped>
.note-editor {
  padding: 10px;
  text-align: left;
  background-color: var(--background-color);
}

.editor-toolbar {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
  align-items: center;
}

.toolbar-item {
  display: flex;
  align-items: center;
}

.toolbar-item label {
  margin-right: 10px;
  font-size: 14px;
  color: var(--text-color);
}

select,
input[type="color"],
input[type="text"] {
  padding: 5px 10px;
  font-size: 14px;
  border-radius: 5px;
  border: 1px solid var(--background-color2);
  font-size: var(--font-size-small);
  color: var(--text-color);
}

.monaco-editor-container {
  height: 700px;
  width: 100%;
  border: 1px solid var(--background-color2);
}

.actions {
  display: flex;
  gap: 10px;
}

button {
  padding: 10px 20px;
  font-size: 16px;
  cursor: pointer;
  color: var(--text-color);
  background-color: var(--button-background-color);
}

button:hover {
  background-color: var(--button-background-color2);
}
</style>
