<template>
  <!-- 消息列表页面容器 -->
  <div class="chat-list">
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
    <!-- 消息列表，使用 v-for 指令循环渲染 chats 数组中的每个消息 -->
    <ul class="chat-items">
      <!-- 每个消息项 -->
      <li 
        v-for="chat in filteredChats" 
        :key="chat.id"
        @contextmenu.prevent="showChatMenu($event, chat)"
        @click = selectChat(chat)
      >
        <div class="chat-avatar">   <!-- 头像-->
          <img :src="chat.avatar" alt="avatar" />
        </div>
        <div class="chat-info">   <!-- 信息-->
          <div class="chat-name">{{ chat.name }}</div>
          <div class="chat-last-chat">{{chat.lastMessage.length > this.maxChars ? chat.lastMessage.slice(0, this.maxChars) + '...' : chat.lastMessage}}</div>
        </div>
        <div class="chat-meta">   <!-- 时间和未读-->
          <div class="chat-time">{{ chat.lastMessageTime }}</div>
          <div v-if="chat.unreadCount" class="unread-count">{{ chat.unreadCount }}</div>   <!--todo-->
        </div>
      </li>
    </ul>
    <!-- 添加好友弹窗 -->
    <AddFriendGroup
      v-if="isAddModalVisible"
      @close="isAddModalVisible = false"
      @add-friend="handleAddFriend"
    />
    <!-- 新建群聊弹窗 -->
    <BuildGroup
      v-if="isBuildModalVisible"
      @close="isBuildModalVisible = false"
      @build-group="handleBuildGroup"
    />
    <ContextMenu ref="contextMenu"  @select-item="handleMenuSelect" />
  </div>
</template>

<script>
import SearchBar from '@/components/base/SearchBar.vue';
import ContextMenu from '@/components/base/ContextMenu.vue';
import * as chatListAPI from '@/services/chatList';
import { addFriendGroup, createGroup } from '@/services/api';
import AddFriendGroup from '@/components/base/AddFriendGroup.vue';
import BuildGroup from '@/components/base/BuildGroup.vue';
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
      // chats: [{
      //     id: '0',   // 好友的tid
      //     avatar: new URL('@/assets/images/avatar.jpg', import.meta.url).href,
      //     name: 'Alice',  // 好友的备注 remark
      //     lastMessage: 'hi',
      //     lastMessageTime: '10:00',
      //     unreadCount: 1,
      //     tags: ['unread','pinned'],   // friend, group, unread, pinned, blocked
      //   },
      //   {
      //     id: '1',
      //     avatar: new URL('@/assets/images/avatar.jpg', import.meta.url).href,
      //     name: 'Bob',
      //     lastMessage: 'hello',
      //     lastMessageTime: '11:00',
      //     unreadCount: 0,
      //     tags: ['unread', 'group'],
      //   }], 
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
      let chats = this.chats;
      if (this.activeTag !== 'all') {
        chats = chats.filter(chat => chat.tags.includes(this.activeTag));
      }
      if(!chats) {
        return chats;
      }
      // 将置顶的消息排在前面
      return chats.sort((a, b) => b.pinned - a.pinned);
    },
    maxChars(){  // 可以显示的字体个数
      return Math.floor((this.chatListWidth - 30) / 12);
    },
  },
  watch:{
    '$store.state.currentChat': {
      handler: function(val) {
        if(val){
          this.selectChat(val);
        }
      },
      immediate: true,
    }
  },
  methods: {
    async fetchChatList() {
      // 从后端获取聊天列表
      let response = await chatListAPI.getChatList();
      if(response.status === 200) {
        this.chats = response.data;
      }
      else{
        console.error('获取聊天列表失败:', response.data);
      }
    },
    // 选中tag筛选消息
    filterChats(tagName) {
      this.activeTag = tagName;
    },
    // 选中消息，切换到对应的聊天
    async selectChat(chat, tid=null) {
      if (!chat) {
        const response = await chatListAPI.getChat(tid);
        chat = response.data;
        this.chats.unshift(chat);
      }
      this.selectedChat = chat;   // todo 滚动到chat
      this.$store.dispatch('setChat', chat);
      
    },
    // 搜索消息
    async handleSearch(keyword) {
      // 搜索聊天列表
      this.chatList = await chatListAPI.searchChats(keyword);
    },
    // 显示新建消息的菜单
    showNewContextMenu(event) {
      this.menuType = 'new';
      const items = [
        '添加好友',
        '新建群聊',
      ];
      this.$refs.contextMenu.show(event, event.clientX, event.clientY, items, null);
    },
    // 右键聊天列表后的菜单
    showChatMenu(event, obj) {
      this.menuType = 'chat';
      let items = [];
      if(obj.tags.includes('unread')) {
        items.push('标记为已读');
      } else {
        items.push('标记为未读');
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
      this.$refs.contextMenu.show(event, event.clientX, event.clientY, items, obj);
    },
    // 处理新建消息的菜单点击事件
    async handleNewMenu(option) {
      if(option === '添加好友') {
        this.isAddModalVisible = true;
      }else if(option === '新建群聊') {
        
        this.isBuildModalVisible = true;
      }
    },
    // 处理聊天列表的菜单点击事件
    async handleChatMenu(option, chat){
      if(option === '置顶') {
        // 告知服务器修改消息的置顶状态（并且本地更新）
        // 置顶聊天
        chat.tags.push('pinned');
        // 告知服务器
        await chatListAPI.pinChat(chat.id, true);
      }else if(option === '取消置顶') {
        // 取消置顶聊天
        chat.tags = chat.tags.filter(tag => tag !== 'pinned');
        // 告知服务器
        await chatListAPI.pinChat(chat.id, false);
      }else if(option === '删除') {
        // 删除聊天
        // 告知服务器
        await chatListAPI.deleteChat(chat.id);
        // 本地删除
        this.chats = this.chats.filter(chat => chat.id !== chat.id);
      }else if(option === '标记为已读') {
        // 标记为已读
        chat.tags = chat.tags.filter(tag => tag !== 'unread');
        // 告知服务器
        await chatListAPI.readMessages(chat.id, true);
      }else if(option === '标记为未读') {
        // 标记为未读
        chat.tags.push('unread');
        // 告知服务器
        await chatListAPI.readMessages(chat.id, false);
      }else if(option === '消息免打扰') {
        // 消息免打扰
        chat.tags.push('mute');
        // 告知服务器
        await chatListAPI.setMute(chat.id, true);
      }else if(option === '取消消息免打扰') {
        // 取消消息免打扰
        chat.tags = chat.tags.filter(tag => tag !== 'mute');
        // 告知服务器
        await chatListAPI.setMute(chat.id, false);
      }else if(option === '屏蔽') {
        // 屏蔽
        chat.tags.push('blocked');
        // 告知服务器
        await chatListAPI.blockChat(chat.id, true);
      }else if(option === '取消屏蔽') {
        // 取消屏蔽
        chat.tags = chat.tags.filter(tag => tag !== 'blocked');
        // 告知服务器
        await chatListAPI.blockChat(chat.id, false);
      }
    },
    // 处理菜单的点击事件
    handleMenuSelect(item, obj) {
      if(this.menuType === 'new') this.handleNewMenu(item);
      if(this.menuType === 'chat') this.handleChatMenu(item, obj);
    },
    // 处理添加好友/群聊的逻辑
    async handleAddFriendGroup(key) {
      try {
        await addFriendGroup(key);
        // 添加成功后的逻辑，如提示用户
        alert(`添加成功`);
      } catch (error) {
        console.error('添加失败:', error);
        alert('添加失败，请重试。');
      }
    },
    // 处理新建群聊的逻辑
    async handleBuildGroup(tids) {
      await createGroup(tids);
    },
  },
  created () {
    this.fetchChatList();
  },
};
</script>

<style scoped>
/* 消息列表页面的样式 */
.chat-list {
  width: 30%;
  height: 100%;
  background-color: #f5f5f5;
}
.chat-header button {
  margin-right: 10px;
  padding: 5px 10px;
  cursor: pointer;
}
.chat-header button.active {
  background-color: #007bff;
  color: white;
}
.chat-items {
  list-style: none;
  padding: 0;
}
.chat-items li {
  display: flex;
  align-items: center;
  padding: 10px;
  padding-bottom: 0px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
}
.chat-items li.unread {
  font-weight: bold;
}
.chat-items li.pinned {
  font-weight: bold;
}
.chat-avatar img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}
.chat-info {
  flex: 5;
  margin-left: 10px;
  text-align: left;
}
.chat-name{
  font-weight: bold;
}
.chat-last-chat {
  color: #888;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.chat-meta {
  text-align: right;
  flex: 1;
}
.chat-time {
  color: #888;
  font-size: 0.8rem;
}
.unread-count {
  background-color: #ff0000;
  color: white;
  width: 15px;
  height: 20px;
  display: inline-block;
  padding: 2px 5px;
  border-radius: 50%;
  text-align: center;
}
</style>