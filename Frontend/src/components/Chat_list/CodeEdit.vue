<template>
  <div v-show="visible" class="editor-container" :style="{ top: `${y}px`, left: `${x}px` }">
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
    <div ref="editor" class="editor"></div>
    <button @click="sendCode">发送</button>
  </div>
</template>

<script>
import * as monaco from 'monaco-editor';
import { EventBus } from '@/components/base/EventBus';
import { toRaw } from 'vue';
export default {
  data() {
    return {
      visible: false,
      x: 0,
      y: 0,
      selectedLanguage: 'javascript',
      value:'',
    };
  },
  methods:{
    changeLanguage(){
      if (this.editor) {
        const model = toRaw(this.editor.getModel()); // 获取当前模型
        if (model) {
          // 使用 monaco.editor.setModelLanguage 来更改模型的语言
          monaco.editor.setModelLanguage(model, this.selectedLanguage);
        }
      }
    },
    sendCode(){
      this.$emit('send-code', this.editor.getValue(), this.selectedLanguage);
      this.hide();
    },
    show(event, boundR) {  // boundD, boundR 为边界的坐标
      EventBus.emit('float-component-opened', this); // 通知其他组件
      const cardWidth = 343;
      const cardHeight = 400;
      const x = event.clientX + cardWidth > boundR ? event.clientX - cardWidth : event.clientX;
      const y = event.clientY - cardHeight -20 < 0 ? 0 : event.clientY - cardHeight -20;
      this.x = x;
      this.y = y;
      this.visible = true;
      console.log('show');
      EventBus.emit('float-component-open', this); // 通知其他组件
    },
    hide() {
      this.visible = false;
      EventBus.emit('hide-float-component'); // 通知其他组件
    },
  },
  mounted() {
    this.editor = monaco.editor.create(this.$refs.editor, {
      value: this.value,
      language: this.selectedLanguage,
      automaticLayout: true,
      lineNumbersMinChars: 2, // 设置行号的最小字符数
      tabSize: 2, // 设置制表符宽度
      minimap: {
        enabled: false, // 禁用右侧的迷你地图
      },
      fontSize: 14, // 设置字体大小
      lineHeight: 20, // 设置行高
      padding: {
        top: 10,
        bottom: 10,
      },
    });

    this.editor.onDidChangeModelContent(() => {
      this.$emit('input', this.editor.getValue());
    });

    EventBus.on('other-float-component', (component) => {
      if (this.visible && component !== this) {
        this.hide();
      }
    });
    EventBus.on('close-float-component', (clickedElement) => {
      if (this.visible && !this.$el.contains(clickedElement)) {
        this.hide();
      }
    });
  },
  watch: {
    value(newValue) {
      if (newValue !== this.editor.getValue()) {
        this.editor.setValue(newValue);
      }
    },
  },
  beforeDestroy() {
    if (this.editor) {
      this.editor.dispose();
    }
    EventBus.off('other-float-component');
    EventBus.off('close-float-component');
  },
};
</script>

<style scoped>
.editor-container {
  position: fixed;
  width: 343px;
  height: 400px;
  z-index: 1000;
  background-color: white;
}

.editor {
  width: 100%;
  height: 85%;
  font-family: 'Fira Code', monospace;
  box-sizing: border-box;
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
</style>