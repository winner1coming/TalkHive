 <!-- 聊天框,上半部分为历史记录，下半部分为输入区-->
<template>
    <div class="chat-box">
      <ChatHeader :chat="selectedChat" />
      <div class="messages" ref="messages">
        <MessageItem 
          v-for="message in messages" 
          :key="message.id" 
          :message="message"
          @message-action="handleMessageAction"
        />
      </div>
      <MessageInput @send-message="sendMessage" />
    </div>
  </template>
  
  <script>
  import ChatHeader from './ChatHeader.vue';
  import MessageItem from './MessageItem.vue';
  import MessageInput from './MessageInput.vue';
  
  export default {
    components: { ChatHeader, MessageItem, MessageInput },
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
    },
  };
  </script>
  
  <style scoped>
  .chat-box {
    display: flex;
    flex-direction: column;
    height: 100%;
  }
  .messages {
    flex: 1;
    overflow-y: auto;
    padding: 10px;
    background-color: #f0f0f0;
  }
  </style>
  