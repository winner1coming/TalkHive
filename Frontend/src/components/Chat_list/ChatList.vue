<template>
  <!-- 消息列表页面容器 -->
  <div class="chat-list">

    <!-- 消息列表的头部 -->
    <div class="chat-header">
      <!-- 搜索框-->
      <SearchBar 
        @search="handleSearch" 
        @button-click="showNewContextMenu($event)"
        :isImmidiate="true"
      />
      <!-- 筛选标签-->
      <div class="chat-list-header">
        <div class="chat-tag">
          <button 
            v-for="tag in tags" 
            :key="tag.name" 
            :class = "{ active: activeTag === tag.name }"
            @click="filterChats(tag.name)"
          >
            {{ tag.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- 消息列表，使用 v-for 指令循环渲染 chats 数组中的每个消息 -->
    <ul class="chat-items">
      <!-- 每个消息项 -->
      <li 
        v-for="chat in filteredChats" 
        :key="chat.id"
        @contextmenu.prevent="showChatMenu($event, chat)"
        @click = selectChat(chat)
        :class="{pinned: chat.tags.includes('pinned'), selected: selectedChat && chat.id === selectedChat.id}"
      >
        <div class="left-part">
          <!-- 头像-->
          <div class="chat-avatar">  
            <img :src="chat.avatar" alt="avatar" />
          </div>
          <!-- 信息-->
          <div class="chat-info">   
            <div class="chat-name" :style="{width: `${chatListWidth-155}px`}">{{ chat.name }}</div>
            <div class="chat-last-chat" :style="{width: `${chatListWidth-155}px`}">{{chat.lastMessage}}</div>
          </div>
        </div>
        <div class="chat-meta">   <!-- 时间和未读-->
          <div class="chat-time">{{ formatTime(chat.lastMessageTime) }}</div>
          <div v-if="chat.tags.includes('mute')" class="mute">🔇</div>
          <div v-else-if="chat.unreadCount" class="unread-count">{{ chat.unreadCount }}</div>
        </div>
      </li>
    </ul>

    <!-- 添加好友弹窗 -->
    <AddFriendGroup
      v-if="isAddModalVisible"
      @close="isAddModalVisible = false"
    />
    <!-- 新建群聊弹窗 -->
    <BuildGroup
      v-if="isBuildModalVisible"
      @close="isBuildModalVisible = false"
    />
    <ContextMenu ref="contextMenu"  @select-item="handleMenuSelect" />
  
  </div>
</template>

<script>
import SearchBar from '@/components/base/SearchBar.vue';
import ContextMenu from '@/components/base/ContextMenu.vue';
import { EventBus } from '@/components/base/EventBus';  
import * as chatListAPI from '@/services/chatList';
import AddFriendGroup from '@/components/Chat_list/AddFriendGroup.vue';
import BuildGroup from '@/components/Chat_list/BuildGroup.vue';
export default {
  components: {
    SearchBar,
    ContextMenu,
    AddFriendGroup,
    BuildGroup,
  },
  props:['chatListWidth'],
  // 组件的 data 函数，返回一个对象，包含组件的响应式数据
  data() {
    return {
      // 消息列表（从后端获取）
      chats: [],
      
      // 选中的聊天
      selectedChat: null,
      // 消息标签
      tags: [
        { name: 'all', label: '全部' },
        { name: 'friend', label: '好友' },
        { name: 'group', label: '群聊' },
        { name: 'unread', label: '未读' },
        { name: 'pinned', label: '置顶' },
        { name: 'blocked', label: '屏蔽' },
      ],
      activeTag: 'all',
      isAddModalVisible: false,
      isBuildModalVisible: false,
      menuType: '',
    };
  },

  computed: {
    // 过滤后的消息列表
    filteredChats() {
      if(!this.chats) {
        return [];
      }
      let chats = this.chats;
      if(this.activeTag === 'blocked') {
        chats = chats.filter(chat => chat.tags.includes(this.activeTag));
      }else if (this.activeTag !== 'all') {
        chats = chats.filter(chat => chat.tags.includes(this.activeTag) && !chat.tags.includes('blocked'));
      }else{   // all不显示被屏蔽的消息
        chats = chats.filter(chat => !chat.tags.includes('blocked'));
      }
      if(!chats) {
        return chats;
      }
      // 将置顶的消息排在前面
      return chats.sort((a, b) => {
        const aPinned = a.tags.includes('pinned') ? 1 : 0;
        const bPinned = b.tags.includes('pinned') ? 1 : 0;
        return bPinned - aPinned;
      });
    },
    // maxChars(){  // 可以显示的字体个数
    //   return Math.floor((this.chatListWidth - 120) / parseInt(this.$store.state.settings.fontSize,10));
    // },
  },
  watch:{
    '$store.state.currentChat': {
      handler: function(val) {
        if(val){
          if(this.selectedChat && val.id!==this.selectedChat.id) this.selectChat(val);
          this.chats = this.chats.map(chat => chat.id === val.id? val : chat);
          if(val.unreadCount > 0){
            this.readMessages(val);
          }
        }
      },
      immediate: true,
      deep: true,
    },
    '$store.state.creatingChat': {
      handler: function(val) {
        if(val){
          this.$store.dispatch('setCreatingChat', false);
          this.$nextTick(() => {
            let chat=null;
            if(this.chats) chat = this.chats.find(chat => chat.id === data.id);
            if(chat){
              this.selectChat(chat);
            }else{
              //console.log(this.$store.state.newChat);
              this.selectChat(null, this.$store.state.newChat.id, this.$store.state.newChat.is_group);
            }
          });
          
        }
      },
      immediate: true,
    },
  },
  methods: {
    async fetchChatList() {
      try{
        // 从后端获取聊天列表
        const response = await chatListAPI.getChatList();
        if(response.status === 200) {
          this.chats = response.data.data;
          if(!this.chats){
            this.chats = [];
          }
        }
        else{
          this.$root.notify(response.data.message, 'error');
        }
      }catch(e){
        console.log(e);
      }
    },
    // 选中tag筛选消息
    filterChats(tagName) {
      this.activeTag = tagName;
    },
    // 选中消息，切换到对应的聊天
    async selectChat(chat, tid=null, is_group=false) {
      console.log(chat, tid, is_group);
      if (!chat) {
        if(!tid) return;
        try{
          console.log('CreateChat');
          const response = await chatListAPI.getChat(tid, is_group);
          console.log(response);
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
            return;
          }
          chat = response.data.data[0];
          
          this.fetchChatList();
        }catch(e){
          console.log(e);
        }
      }
      this.selectedChat = chat;   
      this.$store.dispatch('setChat', chat);
      // 已读消息
      if(chat.tags.includes('unread')) {
        chat.tags = chat.tags.filter(tag => tag !== 'unread');
        chat.unreadCount = 0;
        try{
          let response;
          if(chat.tags.includes('friend')) {
            response = await chatListAPI.readMessages(chat.id, true, false);
          }else{
            response = await chatListAPI.readMessages(chat.id, true, true);
          }
          if(response.status !== 200){
            this.$root.notify(response.data.message, 'error');
          }
        }catch(e){
          console.log(e);
        }
      }
    },
    // 格式化时间
    formatTime(time) {
      if(!time) return '';
      const now = new Date();
      const messageTime = new Date(time);
      const isToday = now.toDateString() === messageTime.toDateString();
      const isYesterday = new Date(now.setDate(now.getDate() - 1)).toDateString() === messageTime.toDateString();

      if (isToday) {
        return messageTime.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      } else if (isYesterday) {
        return '昨天';
      } else {
        return messageTime.toLocaleDateString();
      }
    },
    // 搜索消息
    async handleSearch(keyword) {
      try{
        // 搜索聊天列表
        const response = await chatListAPI.searchChats(keyword);
        if(response.status === 200) {
          this.chatList = response.data.data;
        }else{
          this.$root.notify(response.data.message, 'error');
        }
      }catch(e){
        console.log(e);
      }
    },
    // 已读消息
    async readMessages(chat) {
      // 标记为已读
      try{
        let response;
        if(chat.tags.includes('friend')) {
          response = await chatListAPI.readMessages(chat.id, true, false);
        }else{
          response = await chatListAPI.readMessages(chat.id, true, true);
        }
        chat.tags = chat.tags.filter(tag => tag !== 'unread');
        chat.unreadCount = 0;
      }catch(e){
        console.log(e);
      }
    },
    // 显示新建消息的菜单
    showNewContextMenu(event) {
      this.menuType = 'new';
      const items = [
        '添加好友/群聊',
        '新建群聊',
      ];
      this.$refs.contextMenu.show(event, items, null, null, null);
    },
    // 右键聊天列表后的菜单
    showChatMenu(event, obj) {
      this.menuType = 'chat';
      let items = [];
      if(obj.id !== this.$store.state.currentChat.id) {
        if(obj.tags.includes('unread')) {
          items.push('标记为已读');
        } else {
          items.push('标记为未读');
        }
      }
      if(obj.tags.includes('pinned')) {
        items.push('取消置顶');
      } else {
        items.push('置顶');
      }
      items.push('删除');
      if(obj.tags.includes('mute')) {
        items.push('取消消息免打扰');
      } else {
        items.push('消息免打扰');
      }
      if(obj.tags.includes('blocked')) {
        items.push('取消屏蔽');
      } else {
        items.push('屏蔽');
      }
      this.$refs.contextMenu.show(event, items, obj, this.boundD, this.boundR);
    },
    // 处理新建消息的菜单点击事件
    async handleNewMenu(option) {
      if(option === '添加好友/群聊') {
        this.isAddModalVisible = true;
      }else if(option === '新建群聊') {
        
        this.isBuildModalVisible = true;
      }
    },
    // 处理聊天列表的菜单点击事件
    async handleChatMenu(option, chat){
      if(option === '置顶') {
        // 告知服务器修改消息的置顶状态（并且本地更新）
        try{
          let response;
          if(chat.tags.includes('friend')){
            response = await chatListAPI.pinChat(chat.id, true, false);
          }else{
            response = await chatListAPI.pinChat(chat.id, true, true);
          }
          if(response.status === 200) {
            chat.tags.push('pinned');
          }else{
            this.$root.notify(response.data.message, 'error');
          }
        }catch(e){
          console.log(e);
        }
      }else if(option === '取消置顶') {
        try{
          let response;
          if(chat.tags.includes('friend')){
            response = await chatListAPI.pinChat(chat.id, false, false);
          }else{
            response = await chatListAPI.pinChat(chat.id, false, true);
          }
          if(response.status === 200) {
            chat.tags = chat.tags.filter(tag => tag !== 'pinned');
          }else{
            this.$root.notify(response.data.message, 'error');
          }
        }catch(e){
          console.log(e);
        }
      }else if(option === '删除') {
        // 删除聊天
        try{
          let response;
          if(chat.tags.includes('friend')){
            response = await chatListAPI.deleteChat(chat.id, false);
          }else{
            response = await chatListAPI.deleteChat(chat.id, true);
          }
          if(response.status === 200) {
            this.chats = this.chats.filter(onechat => onechat.id !== chat.id);
          }else{
            this.$root.notify(response.data.message, 'error');
          }
        }catch(e){
          console.log(e);
        }
      }else if(option === '标记为已读') {
        // 标记为已读
        try{
          let response;
          if(chat.tags.includes('friend')) {
            response = await chatListAPI.readMessages(chat.id, true, false);
          }else{
            response = await chatListAPI.readMessages(chat.id, true, true);
          }
          if(response.status !== 200) {
            this.$root.notify(response.data.message, 'error');
          }else{
            chat.tags = chat.tags.filter(tag => tag !== 'unread');
            chat.unreadCount = 0;
          }
        }catch(e){
          console.log(e);
        }
      }else if(option === '标记为未读') {
        // 标记为未读
        try{
          let response;
          if(chat.tags.includes('friend')) {
            response = await chatListAPI.readMessages(chat.id, false, false);
          }else{
            response = await chatListAPI.readMessages(chat.id, false, true);
          }
          if(response.status !== 200) {
            this.$root.notify(response.data.message, 'error');
          }else{
            chat.tags.push('unread');
            chat.unreadCount = 1;
          }
        }catch(e){
          console.log(e);
        }
      }else if(option === '消息免打扰') {
        try{
          let response;
          if(chat.tags.includes('friend')) {
            response = await chatListAPI.setMute(chat.id, true, false);
          }else{
            response = await chatListAPI.setMute(chat.id, true, true);
          }
          if(response.status === 200) {
            chat.tags.push('mute');
          }else{
            this.$root.notify(response.data.message, 'error');
          }
        }catch(e){
          console.log(e);
        }
      }else if(option === '取消消息免打扰') {
        try{
          let response;
          if(chat.tags.includes('friend')) {
            response = await chatListAPI.setMute(chat.id, false, false);
          }else{
            response = await chatListAPI.setMute(chat.id, false, true);
          }
          if(response.status === 200) {
            chat.tags = chat.tags.filter(tag => tag !== 'mute');
          }else{
            this.$root.notify(response.data.message, 'error');
          }
        }catch(e){
          console.log(e);
        }
      }else if(option === '屏蔽') {
        try{
          let response;
          if(chat.tags.includes('friend')) {
            response = await chatListAPI.blockChat(chat.id, true, false);
          }else{
            response = await chatListAPI.blockChat(chat.id, true, true);
          }
          if(response.status === 200) {
            chat.tags.push('blocked');
          }else{
            this.$root.notify(response.data.message, 'error');
          }
        }catch(e){
          console.log(e);
        }
      }else if(option === '取消屏蔽') {
        try{
          let response;
          if(chat.tags.includes('friend')) {
            response = await chatListAPI.blockChat(chat.id, false, false);
          }else{
            response = await chatListAPI.blockChat(chat.id, false, true);
          }
          if(response.status === 200) {
            chat.tags = chat.tags.filter(tag => tag !== 'blocked');
          }else{
            this.$root.notify(response.data.message, 'error');
          }
        }catch(e){
          console.log(e);
        }
      }
    },
    // 处理菜单的点击事件
    handleMenuSelect(item, obj) {
      if(this.menuType === 'new') this.handleNewMenu(item);
      if(this.menuType === 'chat') this.handleChatMenu(item, obj);
    },
  },
  created () {
    this.fetchChatList();
    if(this.$store.state.currentChat){
      this.selectChat(this.$store.state.currentChat);
      if(this.$store.state.currentChat.unreadCount > 0){
        this.readMessages(this.$store.state.currentChat);
      }
    }
  },
  mounted() {
    EventBus.on('set-mute', (tid, is_mute) => {
      for (let i = 0; i < this.chats.length; i++) {
        if (this.chats[i].id === tid) {
          if (is_mute) {
            this.chats[i].tags.push('mute');
          } else {
            this.chats[i].tags = this.chats[i].tags.filter(tag => tag !== 'mute');
          }
          break; 
        }
      } 
    });
    EventBus.on('set-pinned', (tid, is_pinned) => {
      for (let i = 0; i < this.chats.length; i++) {
        if (this.chats[i].id === tid) {
          if (is_pinned) {
            this.chats[i].tags.push('pinned');
          } else {
            this.chats[i].tags = this.chats[i].tags.filter(tag => tag !== 'pinned');
          }
          break; 
        }
      } 
    });
    EventBus.on('set-blocked', (tid, is_blocked) => {
      for (let i = 0; i < this.chats.length; i++) {
        if (this.chats[i].id === tid) {
          if (is_blocked) {
            this.chats[i].tags.push('blocked');
          } else {
            this.chats[i].tags = this.chats[i].tags.filter(tag => tag !== 'blocked');
          }
          break; 
        }
      } 
    });
    EventBus.on('set-blacklist', (tid, is_blacklist) => {
      for (let i = 0; i < this.chats.length; i++) {
        if (this.chats[i].id === tid) {
          if (is_blacklist) {
            this.chats = this.chats.filter(chat => chat.id !== tid);
          } else {
            
          }
          break; 
        }
      } 
    });
    EventBus.on('update-chat', (newChat) => {
      if(!this.chats){
        this.fetchChatList();
        return;
      }
      const chat = this.chats.find(chat => chat.id === newChat.id);
      if(!chat){
        this.fetchChatList();
        return;
      }
      this.chats = this.chats.filter(chat => chat.id !== newChat.id);
      if(this.selectedChat&&newChat.id===this.selectedChat.id){
        // chat.tags = newChat.tags.filter(tag => tag !== 'unread');
        // newChat.unreadCount = 0;
        // if(newChat.tags.includes('friend')){
        //   this.chatListAPI.readMessages(newChat.id, true, false);
        // }
        // else{
        //   this.chatListAPI.readMessages(newChat.id, true, true);
        // }
        
      }else{
        chat.unreadCount = chat.unreadCount+1;
        chat.tags = chat.tags.filter(tag => tag !== 'unread');
        chat.tags.push('unread');
      }
      chat.lastMessage = newChat.lastMessage;
      chat.lastMessageTime = newChat.lastMessageTime;
      this.chats.unshift(chat); 
      // this.fetchChatList();
    });
  },
  beforeUnmount() {
    console.log('destroy');
    EventBus.off('set-mute');
    EventBus.off('set-pinned');
    EventBus.off('set-blocked');
    EventBus.off('set-blacklist');
    EventBus.off('update-chat');
  },
};
</script>

<style scoped src="@/assets/css/chatList.css"></style>
<style scoped>
.chat-list {
  width: 30%;
  height: 100%;
  background-color: var(--background-color);
  color: var(--text-color);
  display: flex;
  flex-direction: column;
  min-width: 170px;
}
.chat-header{
  flex: 1;
}
.chat-tag{
  font-size: var(--font-size-small);
}
.chat-items {
  flex: 9;
  list-style: none;
  padding: 0;
  overflow-x: hidden;
  overflow-y: auto;
}

.chat-items li {
  display: flex;
  align-items: center;
  padding: 10px;
  padding-bottom: 0px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
  justify-content: space-between;
}
.chat-items li.pinned {
  background-color: var(--background-color1);
  color: var(--text-color);
}
.chat-items li.selected {
  background-color: var(--background-color2);
  color: var(--text-color);
}
.left-part {
  display: flex;
  align-items: center;
  flex-direction: row;
}
.chat-avatar img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}
.chat-info {
  margin-left: 10px;
  text-align: left;
}
.chat-name{
  font-weight: bold;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.chat-last-chat {
  color: #888;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: var(--font-size-small-small);

}
.chat-meta {
  text-align: right;
  width: 71px;
}
.chat-time {
  color: #888;
  font-size: var(--font-size-small-small);
}
.unread-count {
  background-color: #d63131df;
  color: white;
  width: 15px;
  height: 20px;
  display: inline-block;
  padding: 2px 5px;
  border-radius: 50%;
  text-align: center;
}
.mute{
  color: #888;
}
</style>