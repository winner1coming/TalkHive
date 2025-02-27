<!-- 聊天框,上半部分为历史记录，下半部分为输入区-->
<template>
  <div class="chat-box" ref="chatBox">
    <!-- 最上方的聊天头部 -->
    <div class="chat-header">
      <div class="chat-avatar">
        <img :src="selectedChat.avatar" alt="avatar" />
      </div>
      <div class="chat-name">{{ selectedChat.name }}</div>
      <div style="margin-left: auto;">
        <div class="detail-button" @click="clickManagement">···</div>
      </div>
    </div>
    <!-- 上方的消息历史 -->
    <div class="messages" ref="messages" :style="messagesStyle" @scroll="handleScroll">
      <MessageItem 
        v-for="message in messages" 
        :message="message"
        @show-context-menu="showContextMenu"
        @show-profile-card="showProfileCard"
      />
      <!-- 滚动到底部按钮 -->
      <p v-show="showScrollButton" class="scroll-to-bottom" @click="scrollToBottom">
        <img src="@/assets/images/down.png" class="scroll-to-bottom-img"/>
      </p>
    </div>
    
    <!-- 下方的输入框 -->
    <MessageInput 
      @send-message="sendMessage" 
      :isBanned="isBanned"
      :isAllBanned="isAllBanned"
    />
    <!-- 右键菜单 -->
    <ContextMenu ref="contextMenu"  @select-item="handleMenuSelect" />
    <!-- 个人名片 -->
    <PersonProfileCard ref="profileCard"/>
  </div>
</template>
  
<script>
import MessageItem from './MessageItem.vue';
import MessageInput from './MessageInput.vue';
import ContextMenu from '@/components/base/ContextMenu.vue';
import PersonProfileCard from '@/components/base/PersonProfileCard.vue';
import * as chatListAPI from '@/services/chatList';
import { getPersonProfileCard } from '@/services/api';
import { EventBus } from '@/components/base/EventBus';

export default {
  components: {MessageItem, MessageInput, ContextMenu, PersonProfileCard},
  data() {
    return {
      messages: [],  // 当前聊天的消息历史
      isBanned:false,  // 是否被禁言
      isAllBanned:false,  // 是否全员禁言
      group_role: null,  // 群聊角色
      boundD: 0,
      boundR: 0,
      selectedChat: null,  // 当前选中的聊天
      showScrollButton: false, // 控制滚动按钮的显示
      backgroundImage:this.$store.state.settings.background,
    }
  },
  computed:{
    messagesStyle(){
      return{
        backgroundImage: `url(${this.backgroundImage})`, // 动态绑定背景图片
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat',
        backgroundSize: 'cover',
      };
    },
  },
  watch: {
    '$store.state.currentChat': {
      immediate: true,
      handler: function(newVar) {
        if(newVar) {
          if(this.selectedChat && this.selectedChat.id === newVar.id) return;
          this.selectedChat = newVar;
          this.selectNewChat(newVar.id);
        }
      },
    },
  },
  methods: {
    async fetchMessages(account_id){
      try{
        const response = await chatListAPI.getMessages(account_id, this.selectedChat.tags.includes('group') ? true : false);
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        if(!response.data.data){
          this.messages = [];
        }
        else{
          this.messages = response.data.data.messages;
          this.isBanned = response.data.data.is_banned;
          this.isAllBanned = response.data.data.is_all_banned;
          this.group_role = response.data.data.group_role;
        }
        this.$nextTick(() => {
          this.scrollToBottom();
        });
      }catch(e){
        console.log(e);
      }
    },
    selectNewChat(account_id) {
      this.fetchMessages(account_id);
      
    },
    async sendMessage(content, type) {
      try{
        let response;
        if(type === 'file'){
            const formData = new FormData();
            formData.append('tid', this.selectedChat.id);
            formData.append('content', content);
            formData.append('is_group', this.selectedChat.tags.includes('group') ? true : false);
            response = await chatListAPI.sendFile(formData);
        }else{
          response = await chatListAPI.sendMessage(this.selectedChat.id, content, type, this.selectedChat.tags.includes('group') ? true : false);
        }
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }else{
          this.messages.push({
            message_id: response.data.data.message_id,  // 消息编号
            send_account_id: this.$store.state.user.id,  // 发送者的id
            avatar:this.$store.state.user.avatar,
            content: content,
            sender: this.$store.state.user.username,   // 发送者的备注
            create_time: response.data.data.create_time,  // 发送时间  todo 改为前端创建
            type: type,   // 消息类型
          });
          let newChat = this.$store.state.currentChat;
          if(type==='text') newChat.lastMessage = content;
          else if(type==='image')newChat.lastMessage = '[图片]';
          else if(type==='file')newChat.lastMessage = '[文件]';
          else newChat.lastMessage = '[代码块]';
          this.$store.dispatch('setChat', newChat);
          this.scrollToBottom();
        }
      }catch(e){
        console.log(e);
      }
      
    },
    clickManagement() {
      // 打开群聊管理弹窗
      this.$emit('click-management');
    },
    showContextMenu(event, message) {
      const items = ['删除', '撤回', '复制','收藏'];
      this.$refs.contextMenu.show(event, items, message, this.boundD, this.boundR);
    },
    async handleMenuSelect(option, message){   
      if(option === '删除'){
        try{
          const response = await chatListAPI.deleteMessage(message.message_id);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
          }else{
            this.messages = this.messages.filter(item => item.message_id !== message.message_id);
          }
        }catch(e){
          console.log(e);
        }
      }
      else if(option === '撤回'){
        try{
          const response = await chatListAPI.recallMessage(message.message_id);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
          }else{
            this.messages = this.messages.filter(item => item.message_id !== message.message_id);
            //逻辑有待完善
          }
        }catch(e){
          console.log(e);
        }  
      }
      else if(option === '复制'){
        try{
          await navigator.clipboard.writeText(message.content);
          this.$root.notify('复制成功','success');
        }catch (err){
          console.error('复制失败:',err);
          this.$root.notify('复制失败','error');
        }
      }
      else if(option === '收藏'){
        try{
          const response = await chatListAPI.collectMessage({table_name:"message",message_id: message.message_id});
          if(response.status != 200){
            this.$root.notify(response.data.message, 'error');
          }else{
            this.$root.notify("收藏成功", 'success');
          }
        }catch(e){
          console.log(e);
        }
      }
      // else if(option === '置顶'){
      //   chatListAPI.topMessage(message.message_id);  // todo
      // }
    
    },
    async showProfileCard(event, send_account_id){
      let group_id = null;
      if(this.selectedChat.tags.includes('group')){
        group_id = this.selectedChat.id;
      }
      try{
        const response = await getPersonProfileCard(send_account_id, group_id); 
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        const profile = response.data.data;
        this.$refs.profileCard.show(event, profile, this.boundD, this.boundR);
      }catch(e){
        console.log(e);
      }
    },

    handleScroll() {
      const messages = this.$refs.messages;
      this.showScrollButton = messages.scrollTop + 50 < messages.scrollHeight - messages.clientHeight;
    },
    scrollToBottom(){
      this.$nextTick(() => {
        const messages = this.$refs.messages;
        if (messages) {
          messages.scrollTop = messages.scrollHeight;
          console.log('scroll to bottom');
        }
      });
    }
  },
  mounted() {
    this.boundD = this.$refs.chatBox.getBoundingClientRect().bottom;
    this.boundR = this.$refs.chatBox.getBoundingClientRect().right;
    EventBus.on('new-message', (newMessage) => {
      this.messages.push(newMessage);
      const messagesContainer = this.$refs.messages;
      if(!messagesContainer) return;
      //console.log(messagesContainer.scrollTop + 900, messagesContainer.scrollHeight - messagesContainer.clientHeight);
      if(messagesContainer.scrollTop + 900 > messagesContainer.scrollHeight - messagesContainer.clientHeight){
        this.scrollToBottom();
      }
      //this.fetchMessages(this.selectedChat.id);
    });
  },
  beforeDestroy() {
    EventBus.off('new-message');
  },
};
</script>
  
<style scoped src="@/assets/css/chatList.css"></style>
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
  background-color: var(--sidebar-background-color);
  color: var(--sidebar-text-color);
}
.chat-name {
  margin-left: 10px;
  font-weight: bold;
}
.detail-button {
  padding: 15px;
  background-color:transparent;
  border: none;
  cursor: pointer;
}
.messages {
  position: relative;
  flex: 1;
  padding: 10px;
  background-color: var(--background-color);
  color: var(--text-color);
  display: block;
  overflow-y: auto; /* 允许垂直滚动 */
  overflow-x: hidden; /* 隐藏水平滚动条 */
}

.scroll-to-bottom {
  position: sticky;
  bottom: 5px;
  right: 5px;
  background-color: #f0f4f9cf;
  color: white;
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  margin-left: auto;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}
.scroll-to-bottom:hover {
  background-color: #c3c7cbe7;
}
.scroll-to-bottom-img{
  width: 20px;
  height: 20px;
}
  </style>