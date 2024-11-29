<template>
    <div class="chat-view">
      <!-- 左侧聊天列表 -->
      <ChatList 
        @chat-selected="selectChat"
      />
  
      <!-- 右侧聊天详情 -->
      <div v-if="selectedChat" class="chat-details">
        <!-- 消息历史 -->
        <ChatBox 
          :selectedChat="selectedChat" 
          :messages="messages" 
          @clickGroupManagement="clickGroupManagement"
          @send-message="sendMessage"
          @message-action="handleMessageAction"
        />
      </div>
  
      <!-- 如果没有选择聊天 -->
      <div v-else class="welcome-message">请选择一个聊天开始！</div>
  
      <!-- 群聊管理弹窗 -->
      <GroupManagement 
        v-if="showGroupManagement" 
        :group="selectedChat" 
        @close="closeGroupManagement"
        @update-group="updateGroupDetails"
      />
    </div>
  </template>
  
  <script>
  import ChatList from '@/components/Chat_list/ChatList.vue';
  import ChatBox from '@/components/Chat_list/ChatBox.vue';
  import GroupManagement from '@/components/Chat_list/GroupManagement.vue';
  import {getMessages} from '@/services/chatList';
  import { EventBus } from '@/components/base/EventBus';
  
  export default {
    components: { ChatList, ChatBox, GroupManagement },
    data() {
      return {
        selectedChat: null, // 当前选中的聊天
        messages: [{
          send_account_id: '0',  // 发送者的id
          content: 'Hello',
          sender: 'Alice',   // 发送者的备注
          timestamp: '11:00',   // 发送时间
        },
        {
          send_account_id: '1',
          content: 'Hi',
          sender: 'Bob',
          timestamp: '12:00',
        }], 
        showGroupManagement: false, // 是否显示群聊管理弹窗
      };
    },
    methods: {
      async selectChat(chat) {
        this.selectedChat = chat;
        // 加载消息历史   
        this.messages = await getMessages(chat.id);
      },
      async sendMessage(content) {   // todo 目前只有发送文字的功能
        if (!this.selectedChat) return;
        // 发送消息到后端
        const newMessage = await this.apiPost(`/messages/${this.selectedChat.id}`, { content });
        // this.messages[this.selectedChat.id].push(newMessage);  todo 消息发送后，是否需要接收自己发送的消息
      },
      handleMessageAction(action, message) {
        // 处理消息的各种操作（复制、删除、多选等）
        console.log(`Action: ${action}`, message);
      },
      clickGroupManagement() {
        if (this.selectedChat && this.selectedChat.is_groupchat) {
          this.showGroupManagement = !this.showGroupManagement;
        }
      },
      closeGroupManagement() {
        this.showGroupManagement = false;
      },
      updateGroupDetails(updatedGroup) {
        // 更新群聊信息
        this.chatList = this.chatList.map(chat =>
          chat.id === updatedGroup.id ? updatedGroup : chat
        );
        if (this.selectedChat.id === updatedGroup.id) {
          this.selectedChat = updatedGroup;
        }
      },
      handleNewMessage(message) {  // todo
        // 处理新消息
        if (this.selectedChat && this.selectedChat.id === message.chatId) {
          this.messages.push(message);
        }
      },
    },
    created() {
      EventBus.on('new-message', (message)=>{this.handleNewMessage(message)});
    },
    beforeDestroy() {
      EventBus.off('new-message');
    }
  };
  </script>
  
  <style scoped>
  .chat-view {
    display: flex;
    align-items: flex-start;
    height: 100%;
    width: 100%;
  }
  .chat-details {
    height: 100%;
    flex: 3;
    display: flex;
    flex-direction: column;
  }
  .welcome-message {
    height: 100%;
    flex: 3;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #888;
    font-size: 1.5rem;
  }
  </style>
  