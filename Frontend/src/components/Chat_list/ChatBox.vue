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
          @show-context-menu="showContextMenu"
          @show-profile-card="showProfileCard"
        />
      </div>
      <MessageInput @send-message="sendMessage" />
      <ContextMenu ref="contextMenu"  @select-item="handleMenuSelect" />
      <ProfileCard ref="profileCard" @go-to-chat="goToChat" />
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
    props: ['selectedChat', 'messages'], // 当前选中的聊天信息和聊天记录
    data() {
    },
    methods: {
      sendMessage(content) {
        // 通知父组件发送消息
        this.$emit('send-message', content);
      },
      clickGroupManagement() {
        // 打开群聊管理弹窗
        this.$emit('clickGroupManagement');
      },
      showContextMenu(event, message) {
        const items = ['引用', '转发', '删除', '撤回', '复制', '多选', '收藏', '置顶'];
        this.$refs.contextMenu.show(event, event.clientX, event.clientY, items, message);
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
        // const profile = await this.getProfileCard(send_account_id); todo
        const profile = {
          tid: '0',  // tid
          avatar: new URL('@/assets/images/avatar.jpg', import.meta.url).href,  // 头像地址
          remark: '',  // 备注
          nickname: '', // 对方设置的昵称
          groupNickname: '',  // 对方的群昵称
          tag: '',  // 分组
          signature: '',  // 个性签名
        };
        this.$refs.profileCard.show(event, profile);
      },
      async goToChat(tid){
        this.$emit('go-to-chat', tid);
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