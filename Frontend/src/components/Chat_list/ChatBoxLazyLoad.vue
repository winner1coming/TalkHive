<template>
    <div class="chat-box" ref="chatBox">
      <div class="chat-header">
        <div class="chat-avatar">
          <img :src="selectedChat.avatar" alt="avatar" />
        </div>
        <div class="chat-name">{{ selectedChat.name }}</div>
        <div style="margin-left: auto;">
          <div class="detail-button" @click="clickManagement">···</div>
        </div>
      </div>
  
      <!-- 虚拟列表消息区域 -->
      <RecycleScroller
        class="messages"
        :items="messages"
        :item-size="80"
        key-field="message_id"
        @scroll="handleVirtualScroll"
        ref="scroller"
        :min-item-size="50"
      >
        <template #default="{ item }">
          <MessageItem 
            :message="item"
            @show-context-menu="showContextMenu"
            @show-profile-card="showProfileCard"
          />
        </template>
      </RecycleScroller>
  
      <!-- 滚动到底部按钮 -->
      <p v-show="showScrollButton" class="scroll-to-bottom" @click="scrollToBottom">
        <img src="@/assets/images/down.png" class="scroll-to-bottom-img"/>
      </p>
  
      <MessageInput 
        @send-message="sendMessage" 
        :isBanned="isBanned"
        :isAllBanned="isAllBanned"
      />
      <ContextMenu ref="contextMenu" @select-item="handleMenuSelect" />
      <PersonProfileCard ref="profileCard"/>
    </div>
  </template>
  
  <script>
  import { RecycleScroller } from 'vue-virtual-scroller';
  import MessageItem from './MessageItem.vue';
  import MessageInput from './MessageInput.vue';
  import ContextMenu from '@/components/base/ContextMenu.vue';
  import PersonProfileCard from '@/components/base/PersonProfileCard.vue';
  import * as chatListAPI from '@/services/chatList';
  import { getPersonProfileCard } from '@/services/api';
  import { EventBus } from '@/components/base/EventBus';
  
  export default {
    components: {
      MessageItem,
      MessageInput,
      ContextMenu,
      PersonProfileCard,
      RecycleScroller
    },
    data() {
      return {
        messages: [],
        isBanned: false,
        isAllBanned: false,
        selectedChat: null,
        showScrollButton: false,
        loading: false,
        noMoreHistory: false,
        backgroundImage: this.$store.state.settings.background,
      };
    },
    computed: {
      messagesStyle() {
        return {
          backgroundImage: `url(${this.backgroundImage})`,
          backgroundPosition: 'center',
          backgroundRepeat: 'no-repeat',
          backgroundSize: 'cover',
        };
      },
    },
    watch: {
      '$store.state.currentChat': {
        immediate: true,
        handler(newVal) {
          if (newVal && (!this.selectedChat || this.selectedChat.id !== newVal.id)) {
            this.selectedChat = newVal;
            this.loadInitialMessages();
          }
        }
      }
    },
    methods: {
      async loadInitialMessages() {
        this.messages = [];
        this.noMoreHistory = false;
        const response = await chatListAPI.getMessages(this.selectedChat.id);
        if (response.status === 200 && response.data.data) {
          this.messages = response.data.data.messages;
          this.isBanned = response.data.data.is_banned;
          this.isAllBanned = response.data.data.is_all_banned;
          this.$nextTick(() => this.scrollToBottom());
        }
      },
      async loadMoreMessages() {
        if (this.loading || this.noMoreHistory) return;
        this.loading = true;
        const firstMsgId = this.messages[0]?.message_id;
        const response = await chatListAPI.getMoreMessages(this.selectedChat.id, firstMsgId);
        if (response.status === 200 && response.data.data) {
          const newMsgs = response.data.data.messages;
          if (newMsgs.length === 0) {
            this.noMoreHistory = true;
          } else {
            this.messages.unshift(...newMsgs);
          }
        }
        this.loading = false;
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
      handleVirtualScroll({ scrollTop }) {
        this.showScrollButton = scrollTop + 50 < this.$refs.scroller.$el.scrollHeight - this.$refs.scroller.$el.clientHeight;
        if (scrollTop === 0) {
          this.loadMoreMessages();
        }
      },
      scrollToBottom() {
        this.$nextTick(() => {
          const scroller = this.$refs.scroller;
          if (scroller) scroller.scrollToItem(this.messages.length - 1);
        });
      },
      clickManagement() {
        this.$emit('click-management');
      },
      showContextMenu(event, message) {
        const items = ['删除', '撤回', '复制', '收藏'];
        this.$refs.contextMenu.show(event, items, message);
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
      async showProfileCard(event, send_account_id) {
        let group_id = this.selectedChat.tags.includes('group') ? this.selectedChat.id : null;
        const response = await getPersonProfileCard(send_account_id, group_id);
        if (response.status === 200) {
          const profile = response.data.data;
          this.$refs.profileCard.show(event, profile);
        }
      },
    },
    mounted() {
      EventBus.on('new-message', (newMessage) => {
        this.messages.push(newMessage);
        const scroller = this.$refs.scroller;
        if (scroller && scroller.$el.scrollTop + 900 > scroller.$el.scrollHeight - scroller.$el.clientHeight) {
          this.scrollToBottom();
        }
      });
    },
    beforeDestroy() {
      EventBus.off('new-message');
    },
  };
  </script>
  
  <style scoped>
  @import "vue-virtual-scroller/dist/vue-virtual-scroller.css";
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
    background-color: transparent;
    border: none;
    cursor: pointer;
  }
  .messages {
    flex: 1;
    overflow-y: auto;
    background-color: var(--background-color);
    color: var(--text-color);
  }
  .scroll-to-bottom {
    position: sticky;
    bottom: 5px;
    right: 5px;
    background-color: #f0f4f9cf;
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
  .scroll-to-bottom-img {
    width: 20px;
    height: 20px;
  }
  </style>
  