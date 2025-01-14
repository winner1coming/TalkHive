<template>
  <div class="note-editor">
    <!-- 工具栏 -->
    <div class="editor-toolbar">
      <div class="toolbar-item">
        <label for="noteFilename">文件名：</label>
        <input type="text" v-model="currentCode.filename" placeholder="输入文件名" />
      </div>

      <div class="toolbar-item">
        <label for="language">语言：</label>
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
        <button @click="saveCode">保存</button>
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
    currentCode() {
      return this.$store.getters.getCurrentCode;  // 从 Vuex 获取 currentCode
    },
  },
  data() {
    return {
      noteContent: "",
      selectedLanguage: '',
      fontSize: '14',
      backgroundColor: '#ffffff',
      editor: null,
      // categories: ['工作', '学习', '个人', '项目'],  // 模拟从后端获取分类
      // 用于处理 v-model 的本地数据
      localFilename: this.note_filename,
    };
  },

  mounted() {
    if (this.currentCode) {
      console.log("currentCode:", this.currentCode);
    } else {
      console.log("currentCode is not loaded yet.");
    }
    // 获取文件后缀并计算语言
    const fileExtension = this.getFileExtension(this.currentCode.filename);
    const language = this.getLanguageByExtension(fileExtension);
    console.log(fileExtension,language);
    nextTick(async() => {
      if (this.$refs.monacoEditor) {
        console.log("currentCode:",this.currentCode.code_id);
        await this.getCodeContent(this.currentCode.code_id);
        this.editor = monaco.editor.create(this.$refs.monacoEditor, {
          value: this.noteContent,
          language: language,
          automaticLayout: true,
          fontSize: parseInt(this.fontSize),
          theme: 'vs',
          backgroundColor: this.backgroundColor,
        });
      }
    });
    this.selectedLanguage=language;
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
    // 定义一个函数来映射后缀名到语言标识符
    getLanguageByExtension(extension) {
        switch (extension) {
            case ".js":
            return "javascript";
            case ".java":
            return "java";
            case ".cpp":
            return "cpp";
            case ".java":
            return "java";
            case ".py":
            return "python";
            case ".html":
            return "html";
            case ".css":
            return "css";
            case ".md":
            return "markdown";
            // 添加更多的语言和后缀名映射
            default:
            return "plaintext"; // 默认语言
        }
    },
    getFileExtension(fname) {
        let dotIndex = fname.lastIndexOf("."); // 找到最后一个点的位置
        let extension = fname.slice(dotIndex);  // 提取点号后的部分
        return extension;
    },
    async getCodeContent(codeID) {
      console.log(codeID);
      const response = await WorkSpaceAPI.getCodeContent(codeID);
      const data = response.data;
      console.log(codeID);
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

    async saveCode() {
      try {
        console.log("进入saveCode");
        const model = toRaw(this.editor.getModel());
        const content = model.getValue();
        console.log(content);
        let dotIndex = this.currentCode.filename.lastIndexOf("."); // 找到最后一个点的位置
        let name = this.currentCode.filename.slice(0, dotIndex);  // 提取点号前的部分
        let extension = this.currentCode.filename.slice(dotIndex);  // 提取点号后的部分
        //let parts = this.currentCode.filename.split(/(?=\.)/); 
        console.log(name,extension);
        const response = await WorkSpaceAPI.saveEditCode(this.currentCode.code_id, name, extension, content);
        if (response.status != 200) {
          alert(response.data.message);
        }
      } catch (error) {
        console.error('无法保存文件:', error);
        alert('保存文件失败！');
      }
      this.$router.push("/workspace/code");
    },


    // 取消编辑
    cancelEdit() {
      this.$router.push("/workspace/code");
    },
  },
};
</script>

<style scoped>
.note-editor {
  padding: 10px;
  text-align: left;
  background-color: var(--background-color);
  color: var(--text-color);
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
  font-size: var(--font-size-small);
  color: var(--text-color);
}

select,
input[type="color"],
input[type="text"] {
  padding: 5px 10px;
  font-size: var(--font-size-small);
  border-radius: 5px;
  border: 1px solid var(--background-color2);
  color: var(--text-color);
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
  font-size: 16px;
  cursor: pointer;
  color: var(--button-text-color);
  background-color: var(--button-background-color);
}

button:hover {
  background-color: var(--button-background-color2);
}
</style>
