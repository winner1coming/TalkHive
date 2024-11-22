<template>
    <div class="chat-view">
      <!-- 左侧聊天列表 -->
      <ChatList 
        :chats="chatsList" 
        @chat-selected="selectChat"
      />
  
      <!-- 右侧聊天详情 -->
      <div v-if="selectedChat" class="chat-details">
        <!-- 顶部：显示聊天信息
        <ChatHeader :chat="selectedChat" @manage-group="openGroupManagement" /> -->
        <!-- 消息历史 -->
        <ChatBox 
          :selectedChat="selectedChat" 
          :messages="messages" 
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
  // import GroupManagement from './GroupManagement.vue';
  
  export default {
    components: { ChatList, ChatBox },
    data() {
      return {
        chatsList: [{
          id: 0,
          avatar: '',
          name: 'Alice',
          lastMessage: 'hi',
          lastMessageTime: '10:00',
          unreadCount: 1,
          tags: ['unread','pinned'],
        },
        {
          id: 1,
          avatar: '',
          name: 'Bob',
          lastMessage: 'hello',
          lastMessageTime: '11:00',
          unreadCount: 0,
          tags: ['unread'],
        }], // 消息列表（从后端获取）
        selectedChat: null, // 当前选中的聊天
        messages: [{
          id: 0,
          content: 'Hello',
          sender: 'Alice',
          timestamp: new Date(),
          
        },
        {
          id: 1,
          content: 'Hi',
          sender: 'Bob',
          timestamp: new Date(),
          
        }], // 消息列表，格式：{ chatId: [{...}, {...}] }
        showGroupManagement: false, // 是否显示群聊管理弹窗
      };
    },
    methods: {
      async fetchChatList() {
        // 从后端获取聊天列表
        this.chatList = await this.apiGet('/chats');
      },
      async selectChat(chat) {
        this.selectedChat = chat;
  
        // 如果消息为空，加载消息历史   todo debug（应该不管是否为空都要加载）
        if (!this.messages) {
          this.messages = await this.apiGet(`/messages/${chat.id}`);
        }else{
          this.messages.forEach(message => {
            // message.read = true;
          });
        }
      },
      async sendMessage(content) {
        if (!this.selectedChat) return;
  
        // 发送消息到后端
        const newMessage = await this.apiPost(`/messages/${this.selectedChat.id}`, { content });
        this.messages[this.selectedChat.id].push(newMessage);
      },
      handleMessageAction(action, message) {
        // 处理消息的各种操作（复制、删除、多选等）
        console.log(`Action: ${action}`, message);
      },
      openGroupManagement() {
        if (this.selectedChat && this.selectedChat.type === 'group') {
          this.showGroupManagement = true;
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
      async apiGet(url) {
        // 模拟后端请求
        console.log(`GET ${url}`);
        return [];
      },
      async apiPost(url, data) {
        // 模拟后端请求
        console.log(`POST ${url}`, data);
        return { id: Date.now(), content: data.content, sender: 'You', timestamp: new Date() };
      },
    },
    mounted() {
      this.fetchChatList(); // 初始化加载聊天列表
    },
  };
  </script>
  
  <style scoped>
  .chat-view {
    display: flex;
    height: 100%;
  }
  .chat-details {
    flex: 3;
    display: flex;
    flex-direction: column;
  }
  .welcome-message {
    flex: 3;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #888;
    font-size: 1.5rem;
  }
  </style>
  