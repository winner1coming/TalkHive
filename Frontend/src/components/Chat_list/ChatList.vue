<template>
  <!-- 消息列表页面容器 -->
  <div class="chat-list">
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
        @click = selectChat(chat)
      >
        <div class="chat-avatar">   <!-- 头像-->
          <img :src="chat.avatar" alt="avatar" />
        </div>
        <div class="chat-info">
          <div class="chat-name">{{ chat.name }}</div>
          <div class="chat-last-chat">{{ chat.lastMessage }}</div>
        </div>
        <div class="chat-meta">
          <div class="chat-time">{{ chat.lastMessageTime }}</div>
          <div v-if="chat.unreadCount" class="unread-count">{{ chat.unreadCount }}</div>   <!--todo-->
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  // 从父组件中接收到消息列表
  props:['chats'],
  // 组件的 data 函数，返回一个对象，包含组件的响应式数据
  data() {
    return {
      // 消息标签
      tags: [
        { name: 'all', label: '全部' },
        { name: 'friend', label: '好友' },
        { name: 'group', label: '群聊' },
        { name: 'unread', label: '未读' },
        { name: 'pinned', label: '置顶' },
        { name: 'blocked', label: '屏蔽' },
      ],
      // // 消息列表数据，每个消息包含 id 和 title 属性
      // chats: [
      //   { id: 1, title: '消息1' },
      //   { id: 2, title: '消息2' },
      //   // 更多消息
      // ],
      activeTag: 'all',
    };
  },

  computed: {
    // 过滤后的消息列表
    filteredChats() {
      let chats = this.chats;
      if (this.activeTag !== 'all') {
        chats = chats.filter(chat => chat.tags.includes(this.activeTag));
      }
      // 将置顶的消息排在前面
      return chats.sort((a, b) => b.pinned - a.pinned);
    },
  },

  methods: {
    // 选中tag筛选消息
    filterChats(tagName) {
      this.activeTag = tagName;
      console.log('tag selected:', tagName);  // debug
    },
    // 选中消息，切换到对应的聊天
    selectChat(chat) {
      this.$emit('chat-selected', chat);
    },
  },
};
</script>

<style scoped>
/* 消息列表页面的样式 */
.chat-list {
  width: 30%;
  background-color: #f5f5f5;
  padding: 10px;
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