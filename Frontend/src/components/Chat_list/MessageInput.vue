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
    />
  </div>
</template>

<script>
import 'emoji-picker-element';
import Emoji from './Emoji.vue';
import CodeEdit from './CodeEdit.vue';
export default {
  components: { Emoji, CodeEdit },
  data() {
    return {
      content: '',
    };
  },
  methods: {
    sendMessage() {
      if (this.content.trim()) {
        this.$emit('send-message', this.content, 'text');
        this.content = '';
      }
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
    sendScreenshot() {
      // 发送截图逻辑
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



</style>