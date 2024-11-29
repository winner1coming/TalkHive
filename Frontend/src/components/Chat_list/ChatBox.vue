<!-- 聊天框,上半部分为历史记录，下半部分为输入区-->
<template>
    <div class="chat-box">
      <div class="chat-header">
        <div class="chat-avatar">
          <img :src="selectedChat.avatar" alt="avatar" />
        </div>
        <div class="chat-name">{{ selectedChat.name }}</div>
        <div style="margin-left: auto;" v-if="selectedChat.tags.includes('group')">
          <button class="group-button" @click="clickGroupManagement">···</button>
        </div>
      </div>
      <div class="messages" ref="messages">
        <MessageItem 
          v-for="message in messages" 
          :message="message"
          :avatar="selectedChat.avatar"
          @message-action="handleMessageAction"
        />
      </div>
      <MessageInput @send-message="sendMessage" />
    </div>

  </template>
  
  <script>
  import MessageItem from './MessageItem.vue';
  import MessageInput from './MessageInput.vue';
  
  export default {
    components: {MessageItem, MessageInput },
    props: ['selectedChat', 'messages'], // 当前选中的聊天信息和聊天记录
    methods: {
      sendMessage(content) {
        // 通知父组件发送消息
        this.$emit('send-message', content);
      },
      handleMessageAction(action, message) {
        // 处理消息的各种操作
        this.$emit('message-action', action, message);
      },
      clickGroupManagement() {
        // 打开群聊管理弹窗
        this.$emit('clickGroupManagement');
      }
    },
  };
  </script>
  
  <style scoped>
  .chat-box {
    display: flex;
    flex-direction: column;
    height: 100%;
  }
  .chat-header {
    display: flex;
    align-items: center;
    padding: 10px;
    background-color: #687aec91;
  }
  .chat-avatar img {
    width: 40px;
    height: 40px;
    border-radius: 50%;
  }
  .chat-name {
    margin-left: 10px;
    font-weight: bold;
  }
  .group-button {
    background-color:transparent;
    border: none;
    cursor: pointer;
  }
  .messages {
    flex: 1;
    padding: 10px;
    background-color: #f0f0f0;
    display: block;
  }
  </style>