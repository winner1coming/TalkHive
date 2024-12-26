<!-- 聊天框,上半部分为历史记录，下半部分为输入区-->
<template><!-- todo List:新消息的处理逻辑 -->
  <div class="chat-box" ref="chatBox">
    <!-- 最上方的聊天头部 -->
    <div class="chat-header">
      <div class="chat-avatar">
        <img :src="selectedChat.avatar" alt="avatar" />
      </div>
      <div class="chat-name">{{ selectedChat.name }}</div>
      <div style="margin-left: auto;" v-if="selectedChat.tags.includes('group')">
        <button class="detail-button" @click="clickGroupManagement">···</button>
      </div>
    </div>
    <!-- 上方的消息历史 -->
    <div class="messages" ref="messages" :style="messagesStyle">
      <MessageItem 
        v-for="message in messages" 
        :message="message"
        :avatar="selectedChat.avatar"
        @show-context-menu="showContextMenu"
        @show-profile-card="showProfileCard"
      />
    </div>
    <!-- 下方的输入框 -->
    <MessageInput @send-message="sendMessage" />
    <!-- 右键菜单 -->
    <ContextMenu ref="contextMenu"  @select-item="handleMenuSelect" />
    <!-- 个人名片 -->
    <ProfileCard ref="profileCard"/>
  </div>
</template>
  
<script>
import MessageItem from './MessageItem.vue';
import MessageInput from './MessageInput.vue';
import ContextMenu from '@/components/base/ContextMenu.vue';
import ProfileCard from '@/components/base/ProfileCard.vue';
import * as chatListAPI from '@/services/chatList';
import { getProfileCard } from '@/services/api';

export default {
  components: {MessageItem, MessageInput, ContextMenu, ProfileCard},
  data() {
    return {
      messages: [],  // 当前聊天的消息历史
      boundD: 0,
      boundR: 0,
      selectedChat: null,  // 当前选中的聊天
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
    async selectNewChat(account_id) {
      // 加载消息历史   
      const response = await chatListAPI.getMessages(account_id);
      // 若被禁言  
      //todo
      this.messages = response.data.messages;
      this.$nextTick(() => {
        this.scrollToBottom();
      });
      
    },
    sendMessage(content) {
      // todo api
      this.messages.push({
        message_id: '0',  // 消息编号
        send_account_id: this.$store.state.user.id,  // 发送者的id
        content: content,
        sender: this.$store.state.user.username,   // 发送者的备注
        create_time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),  // 发送时间
        type: 'text',   // 消息类型
      });
      this.scrollToBottom();
    },
    clickGroupManagement() {
      // 打开群聊管理弹窗
      this.$emit('clickGroupManagement');
    },
    showContextMenu(event, message) {
      const items = ['引用', '转发', '删除', '撤回', '复制', '多选', '收藏', '置顶'];
      this.$refs.contextMenu.show(event, items, message, this.boundD, this.boundR);
    },
    async handleMenuSelect(option, message){   // todo api没搞完
      if(option === '引用'){
        // todo
        this.$emit('reply', message);
      }else if(option === '转发'){
        this.$emit('forward', message);
      }else if(option === '删除'){
        chatListAPI.deleteMessage(message.message_id);
      }else if(option === '撤回'){
        chatListAPI.recallMessage(message.message_id);
      }else if(option === '复制'){
        // todo
      }else if(option === '多选'){
        // todo
      }else if(option === '收藏'){
        chatListAPI.collectMessage(message.message_id);
      }else if(option === '置顶'){
        chatListAPI.topMessage(message.message_id);  // todo
      }
    
    },
    async showProfileCard(event, send_account_id){
      let group_id = null;
      if(this.selectedChat.tags.includes('group')){
        group_id = this.selectedChat.id;
      }
      const response = await getProfileCard(send_account_id, group_id); 
      const profile = response.data;
      // const profile = {
      //   tid: '0',  // tid
      //   avatar: new URL('@/assets/images/avatar.jpg', import.meta.url).href,  // 头像地址
      //   remark: '',  // 备注
      //   nickname: '', // 对方设置的昵称
      //   groupNickname: '',  // 对方的群昵称
      //   tag: '',  // 分组
      //   signature: '',  // 个性签名
      // };
      
      this.$refs.profileCard.show(event, profile, this.boundD, this.boundR);
    },

    scrollToBottom(){
      const messages = this.$refs.messages;
      if (messages) {
        messages.scrollTop = messages.scrollHeight;
      }
    }
  },
  mounted() {
    this.boundD = this.$refs.chatBox.getBoundingClientRect().bottom;
    this.boundR = this.$refs.chatBox.getBoundingClientRect().right;
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
    background-color: #687aec91;
  }
  .chat-name {
    margin-left: 10px;
    font-weight: bold;
  }
  .detail-button {
    background-color:transparent;
    border: none;
    cursor: pointer;
  }
  .messages {
    flex: 1;
    padding: 10px;
    background-color: #f0f0f0;
    display: block;
    overflow-y: auto; /* 允许垂直滚动 */
    overflow-x: hidden; /* 隐藏水平滚动条 */
  }
  </style>