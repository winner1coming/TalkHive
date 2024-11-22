<template>
  <div class="message-input">
    <div class="input-toolbar">
      <button @click="sendEmoji">表情</button>
      <button @click="sendFile">文件</button>
      <input type="file" ref="fileInput" style="display: none;" @change="handleFile" />
      <button @click="sendScreenshot">截图</button>
      <button @click="sendCodeBlock">代码块</button>
      <button @click="sendPoll">群投票</button>
    </div>
    <textarea 
      v-model="inputText" 
      placeholder="输入消息..." 
      @keydown.enter.prevent="sendMessage"
    ></textarea>
    <div class="actions">
      <button @click="send">发送</button>
      
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      content: ''
    };
  },
  methods: {
    sendMessage() {
      if (this.content.trim()) {
        this.$emit('send-message', this.content);
        this.content = '';
      }
    },
    sendEmoji() {
      // 发送表情逻辑
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
    sendCodeBlock() {
      // 发送代码块逻辑
    },
    sendPoll() {
      // 发送群投票逻辑
    }
  }
};
</script>

<style scoped>
.message-input {
  display: flex;
  flex-direction: column;
}
.input-toolbar {
  display: flex;
  justify-content: space-around;
}
textarea {
  flex: 1;
  resize: none;
}
</style>