<template>
  <div class="chat-view">
    <!-- 左侧聊天列表 -->
    <ChatList 
      ref="chatList"
      :chatListWidth="chatListWidth"
      :style="{ width: chatListWidth + 'px' }"
    />
    <!-- 拖动条 -->
    <div class="resizer" @mousedown="startResize"></div>
    <!-- 右侧聊天详情 -->
    <div v-if="this.$store.state.currentChat" class="chat-details">
      <!-- 消息历史 -->
      <ChatBox 
        @clickGroupManagement="clickGroupManagement"
        @send-message="sendNewMessage"
      />
    </div>

    <!-- 如果没有选择聊天 -->
    <div v-else class="welcome-message">请选择一个聊天开始！</div>

    <!-- 群聊管理弹窗 -->
    <GroupManagement 
      ref="groupManagement"
      :group="this.$store.state.currentChat" 
      @close="closeGroupManagement"
      @update-group="updateGroupDetails"
    />
  </div>
</template>
  
<script>
import ChatList from '@/components/Chat_list/ChatList.vue';
import ChatBox from '@/components/Chat_list/ChatBox.vue';
import GroupManagement from '@/components/Chat_list/GroupManagement.vue';
import { EventBus } from '@/components/base/EventBus';

export default {
  components: { ChatList, ChatBox, GroupManagement },
  data() {
    return {
      
      // messages: [{
      //   message_id: '0',  // 消息编号
      //   send_account_id: '0',  // 发送者的id
      //   content: 'Hello',
      //   sender: 'Alice',   // 发送者的备注
      //   create_time: '11:00',   // 发送时间
      //   type: 'text',   // 消息类型
      // },
      // {
      //   send_account_id: '1',
      //   content: 'Hi',
      //   sender: 'Bob',
      //   timestamp: '12:00',
      // }], 
      showGroupManagement: false, // 是否显示群聊管理弹窗
      showFriendManagement: false, // 是否显示好友管理弹窗
      chatListWidth: 300,  // 聊天列表的宽度
    };
  },
  methods: {
    async sendNewMessage(content) {   // todo 目前只有发送文字的功能
      if (!this.$store.state.currentChat) return;
      // 发送消息到后端
      await sendMessage(this.$store.state.currentChat.id, content);
      // this.messages[this.selectedChat.id].push(newMessage);  todo 消息发送后，是否需要接收自己发送的消息
    },
    // async goToChat(tid) {
    //   // 跳转到指定聊天
    //   const chat = this.$ref.chatList.chats.find(chat => chat.id === tid);
    //   this.selectChat(chat);
    // },
    clickGroupManagement() {
      if (this.$store.state.currentChat.tags.includes('group')) {
        this.$refs.groupManagement.show();
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
      // if (this.selectedChat.id === updatedGroup.id) {
      //   this.selectedChat = updatedGroup;
      // }
    },
    handleNewMessage(message) {  // todo！
      // 处理新消息
      if (this.selectedChat && this.selectedChat.id === message.chatId) {
        this.messages.push(message);
      }
    },

    // 拖动条的逻辑
    startResize(event) {
      this.isResizing = true;
      document.addEventListener('mousemove', this.resize);
      document.addEventListener('mouseup', this.stopResize);
    },
    resize(event) {
      if (this.isResizing) {
        this.chatListWidth = event.clientX;
      }
    },
    stopResize() {
      this.isResizing = false;
      document.removeEventListener('mousemove', this.resize);
      document.removeEventListener('mouseup', this.stopResize);
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
  overflow: hidden;
}
.resizer {
  width: 3px;
  height: 100%;
  cursor: ew-resize;
  background-color: #ccc;
}
.chat-details {
  height: 100%;
  width: 100%;
  flex: 3;
  display: flex;
  flex-direction: column;
}
.welcome-message {
  height: 100%;
  width: 100%;
  flex: 3;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #888;
  font-size: 1.5rem;
}
</style>
  