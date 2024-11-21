<template>
  <!-- 单聊页面容器 -->
  <div class="single-chat">
    <!-- 页面标题 -->
    <h2>单聊</h2>
    
    <!-- 聊天窗口，显示消息列表 -->
    <div class="chat-window">
      <!-- 使用 v-for 指令循环渲染 messages 数组中的每个消息 -->
      <div v-for="message in messages" :key="message.id" class="message">
        <!-- 显示消息内容 -->
        <p>{{ message.content }}</p>
      </div>
    </div>
    
    <!-- 输入框，使用 v-model 双向绑定到 data 中的 newMessage 属性 -->
    <input type="text" v-model="newMessage" placeholder="输入消息" />
    
    <!-- 发送按钮，点击时触发 sendMessage 方法 -->
    <button @click="sendMessage">发送</button>
  </div>
</template>

<script>
export default {
  // 组件的 data 函数，返回一个对象，包含组件的响应式数据
  data() {
    return {
      // 新消息输入框的值
      newMessage: '',
      // 消息列表数据，每个消息包含 id 和 content 属性
      messages: [
        { id: 1, content: '你好' },
        { id: 2, content: '你好，有什么事吗？' },
        // 更多消息
      ],
    };
  },
  
  // 组件的方法定义
  methods: {
    // 发送消息方法，处理用户点击发送按钮时的逻辑
    sendMessage() {
      // 调用 Vuex store 中的 sendMessage 动作，传递新消息内容
      this.$store.dispatch('sendMessage', { content: this.newMessage });
      
      // 清空输入框
      this.newMessage = '';
    },
  },
};
</script>

<style scoped>
/* 单聊页面的样式 */
.single-chat {
  padding: 20px; /* 设置内边距 */
}

/* 聊天窗口的样式 */
.chat-window {
  height: 300px; /* 设置高度 */
  overflow-y: scroll; /* 设置垂直滚动条 */
  border: 1px solid #ccc; /* 设置边框 */
  padding: 10px; /* 设置内边距 */
}

/* 消息项的样式 */
.message {
  margin-bottom: 10px; /* 设置底部外边距 */
}
</style>