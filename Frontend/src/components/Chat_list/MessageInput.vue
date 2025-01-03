<template>
  <div class="message-input">
    <div class="input-toolbar">
      <button @click="toggleEmojiPicker">表情</button>
      <button @click="sendFile">文件</button>
      <input type="file" ref="fileInput" style="display: none;" @change="handleFile" />
      <button @click="sendScreenshot">截图</button>
      <button @click="sendCodeBlock">代码块</button>
    </div>
    <textarea 
      v-model="content" 
      placeholder="输入消息..." 
      @keydown.enter.prevent="sendMessage"
    ></textarea>
    <button class="send-button" @click="sendMessage">发送</button>
    <Emoji 
      ref="emojiPicker"
      @emoji-click="addEmoji"
    />
    <CodeEdit 
      ref="codeEditor"
      @send-code="sendCode"
    />
        <!-- 截图区域选择工具 -->
    <div v-if="isSelecting" class="screenshot-overlay" @mousedown="startSelection" @mousemove="resizeSelection" @mouseup="endSelection">
      <div class="selection-box" :style="selectionBoxStyle"></div>
    </div>
      <!-- 确认和取消按钮 -->
    <div v-if="showConfirmButtons" class="confirm-buttons" :style="confirmButtonsStyle">
      <button @click.stop="confirmScreenshot">确认</button>
      <button @click.stop="cancelScreenshot">取消</button>
    </div>
  </div>

</template>

<script>
import html2canvas from 'html2canvas';
import 'emoji-picker-element';
import Emoji from './Emoji.vue';
import CodeEdit from './CodeEdit.vue';
export default {
  components: { Emoji, CodeEdit },
  data() {
    return {
      content: '',
      isSelecting: false, // 是否正在选择截图区域
      startX: 0, // 选择区域的起始 X 坐标
      startY: 0, // 选择区域的起始 Y 坐标
      endX: 0, // 选择区域的结束 X 坐标
      endY: 0, // 选择区域的结束 Y 坐标
      isDragging: false, // 是否正在拖动选择区域
      showConfirmButtons: false, // 是否显示确认和取消按钮
    };
  },
  computed: {
    // 计算选择区域的样式
    selectionBoxStyle() {
      const width = Math.abs(this.endX - this.startX);
      const height = Math.abs(this.endY - this.startY);
      const left = Math.min(this.startX, this.endX);
      const top = Math.min(this.startY, this.endY);
      return {
        width: `${width}px`,
        height: `${height}px`,
        left: `${left}px`,
        top: `${top}px`,
      };
    },
    // 计算确认和取消按钮的位置
    confirmButtonsStyle() {
      const left = Math.min(this.startX, this.endX) + Math.abs(this.endX - this.startX) + 10;
      const top = Math.min(this.startY, this.endY);
      return {
        left: `${left}px`,
        top: `${top}px`,
      };
    },
  },
  methods: {
    sendMessage() {
      if (this.content.trim()) {
        this.$emit('send-message', this.content, 'text');
        this.content = '';
      }
    },
    sendCode(code, language) {
      this.$emit('send-message', code, language);
    },
    toggleEmojiPicker(event) {
      this.$refs.emojiPicker.show(event, window.innerHeight, window.innerWidth);
    },
    addEmoji(event) {
      this.content += event.detail.unicode;
    },
    sendFile() {
      // 发送文件逻辑
      this.$refs.fileInput.click();
      
    },
    handleFile(event) {
      // 处理文件逻辑
      const file = event.target.files[0];
      if (file) {
        this.$emit('send-message', { type: 'file', content: file });
      }
    },
    startScreenshot() {
      this.isSelecting = true; // 进入截图模式
      this.showConfirmButtons = false; // 隐藏确认和取消按钮
    },
    startSelection(event) {
      this.isDragging = true;
      this.startX = event.clientX;
      this.startY = event.clientY;
      this.endX = event.clientX;
      this.endY = event.clientY;
    },
    resizeSelection(event) {
      if (this.isDragging) {
        this.endX = event.clientX;
        this.endY = event.clientY;
      }
    },
    endSelection() {
      this.isDragging = false;
      this.showConfirmButtons = true; // 显示确认和取消按钮
    },
    async confirmScreenshot() {
      const canvas = await html2canvas(document.body, {
        x: Math.min(this.startX, this.endX),
        y: Math.min(this.startY, this.endY),
        width: Math.abs(this.endX - this.startX),
        height: Math.abs(this.endY - this.startY),
      });

      // 将 Canvas 转换为图片文件
      const imageData = canvas.toDataURL('image/png');
      this.$emit('send-message',imageData,'image');

      this.cancelScreenshot(); // 退出截图模式
    },
    cancelScreenshot() {
      console.log("我退出啦！");
      this.isSelecting = false; // 退出截图模式
      this.showConfirmButtons = false; // 隐藏确认和取消按钮
    },
    sendScreenshot() {
      // 发送截图逻辑
      console.log("在截图，勿扰");
      this.startScreenshot();
    },
    sendCodeBlock(event) {
      // 发送代码块逻辑
      this.$refs.codeEditor.show(event, window.innerHeight, window.innerWidth);
    },
  }
};
</script>

<style scoped src="@/assets/css/chatList.css"></style>
<style scoped>
.message-input {
  display: flex;
  flex-direction: column;
  padding: 10px;
  background-color: #f5f5f5;
}
textarea {
  flex: 1;
  resize: none;
  padding-right: 50px;
  border:none;
  background-color: #f5f5f5;
  outline: none;
}
.input-toolbar {
  display: flex;
  align-items: center;
  position: relative;
}
.send-button{
  margin-left: auto;
  width: auto;
}

/* 截图区域选择工具 */
.screenshot-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000; /* 确保 z-index 足够高 */
  cursor: crosshair;
}

.selection-box {
  position: absolute;
  border: 2px dashed #fff;
  background: rgba(255, 255, 255, 0.2);
  z-index: 1001;
}

/* 确认和取消按钮 */
.confirm-buttons {
  position: absolute;
  display: flex;
  gap: 10px;
  z-index: 1002;
}

.confirm-buttons button {
  padding: 5px 10px;
  background-color: #4CAF50;
  color: white;
  border: none;
  cursor: pointer;
  pointer-events: auto;
}

.confirm-buttons button:last-child {
  background-color: #f44336;
}

</style>