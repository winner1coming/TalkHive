<template>
  <div class="note-editor">
    <!-- 工具栏 -->
    <div class="editor-toolbar">
      <div class="toolbar-item">
        <label for="language">语言选择：</label>
        <select v-model="selectedLanguage" @change="changeLanguage">
          <option value="markdown">Markdown</option>
          <option value="javascript">JavaScript</option>
          <option value="python">Python</option>
          <option value="cpp">C++</option>
          <option value="java">Java</option>
          <option value="html">HTML</option>
        </select>
      </div>

      <div class="toolbar-item">
        <label for="fontSize">字体大小：</label>
        <select v-model="fontSize" @change="changeFontSize">
          <option value="12">12px</option>
          <option value="14">14px</option>
          <option value="16">16px</option>
          <option value="18">18px</option>
          <option value="20">20px</option>
        </select>
      </div>

      <div class="toolbar-item">
        <label for="bgColor">背景颜色：</label>
        <input type="color" v-model="backgroundColor" @input="changeBackgroundColor">
      </div>

      <div class="actions">
        <button @click="saveNote">保存</button>
        <button @click="cancelEdit">取消</button>
      </div>
    </div>

    <!-- Monaco 编辑器 -->
    <div ref="monacoEditor" class="monaco-editor-container"></div>

  </div>
</template>

<script>
import * as monaco from 'monaco-editor';
import { ref, onMounted,toRaw, nextTick } from 'vue';

export default {
  name: "NoteEditor",
  props: ["id"],
  data() {
    return {
      note: null,
      noteContent: "",
      selectedLanguage: 'markdown',
      fontSize: '14',
      backgroundColor: '#ffffff',
      editor: null,
    };
  },
  created() {
    const notes = [
      { id: 1, filename: "Vue学习笔记.md", content: "这是 Vue 学习笔记的内容..." },
      { id: 2, filename: "项目需求分析.docx", content: "这是项目需求分析的内容..." },
      { id: 3, filename: "代码优化方案.txt", content: "这是代码优化方案的内容..." },
    ];
    this.note = notes.find(note => note.id === Number(this.id));
    if (this.note) {
      this.noteContent = this.note.content;
    } else {
      alert("笔记未找到！");
      this.$router.push("/workspace/notes");
    }
  },
  mounted() {
    nextTick(() => {
      if (this.$refs.monacoEditor) {
        this.editor = monaco.editor.create(this.$refs.monacoEditor, {
          value: this.noteContent,
          language: this.selectedLanguage,
          automaticLayout: true,
          fontSize: parseInt(this.fontSize),  // 初始化字体大小
          theme: 'vs',  // 默认主题
          backgroundColor: this.backgroundColor,
        });
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
        // 更新背景颜色
        const editorContainer = this.$refs.monacoEditor;
        editorContainer.style.backgroundColor = newColor;
      }
    },
  },
  methods: {
    changeLanguage() {
      if (this.editor) {
        const model = toRaw(this.editor.getModel()); // 获取当前模型
        if (model) {
          // 使用 monaco.editor.setModelLanguage 来更改模型的语言
          monaco.editor.setModelLanguage(model, this.selectedLanguage);
        }
      }
    },
    changeFontSize() {
      if (this.editor) {
        this.editor.updateOptions({ fontSize: parseInt(this.fontSize) });
      }
    },
    changeBackgroundColor() {
      if (this.editor) {
        const editorContainer = this.$refs.monacoEditor;
        console.log(editorContainer)
        console.log(this.$refs.monacoEditor)
        editorContainer.style.backgroundColor = this.backgroundColor;
        console.log(this.backgroundColor)
        console.log(editorContainer.style.backgroundColor)
      }
    },
    saveNote() {
      alert(`笔记已保存:\n${this.noteContent}`);
      this.$router.push("/workspace/notes");
    },
    
    // async saveFile() {
    //   try {
    //     // 发送请求到后端，保存新建的文件
    //     const response = await axios.post('/workspace/create-file', {
    //       filename: this.newFile.filename + this.newFile.filetype,  // 文件名和文件格式拼接
    //     });

    //     if (response.data.status === 200) {
    //       alert('文件创建成功');
    //       this.showCreateFile = false;  // 关闭编辑框
    //     } else {
    //       alert(response.data.message);
    //     }
    //   } catch (error) {
    //     console.error('无法创建文件:', error);
    //     alert('创建文件失败！');
    //   }
    // },
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
}

.editor-toolbar {
  display: flex;
  gap: 20px;
  margin-bottom: 10px;
  align-items: center;
}

.toolbar-item {
  display: flex;
  align-items: center;
}

.toolbar-item label {
  margin-right: 10px;
  font-size: var(--font-size-small);
}

select,
input[type="color"] {
  padding: 5px 10px;
  font-size: var(--font-size-small);
  border-radius: 5px;
  border: 1px solid #ccc;
}

.monaco-editor-container {
  height: 700px;
  width: 100%;
  border: 1px solid #ccc;
}

.actions {
  display: flex;
  gap: 10px;
}

button {
  padding: 10px 20px;
  font-size: var(--font-size);
  cursor: pointer;
}

button:hover {
  background-color: #eaeaea;
}
</style>
