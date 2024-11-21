<template>
  <!-- 消息列表页面容器 -->
  <div class="chat-list">
    <!-- 页面标题 -->
    <h2>消息列表</h2>
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
    <ul>
      <!-- 每个消息项 -->
      <li 
        v-for="chat in chats" 
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
          <span class="chat-time">{{ chat.lastMessageTime }}</span>
          <span v-if="chat.unreadCount" class="unread-count">{{ chat.unreadCount }}</span>   <!--todo-->
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  // 从父组件中接收到消息列表
  props['chats'],
  // 组件的 data 函数，返回一个对象，包含组件的响应式数据
  data() {
    return {
      // 消息标签
      tags: [
        { name: 'all', label: '全部' },
        { name: 'friend', label: '好友' },
        { name: 'group', label: '群聊' },
        { name: 'unread', label: '未读' },
      ],
      // // 消息列表数据，每个消息包含 id 和 title 属性
      // chats: [
      //   { id: 1, title: '消息1' },
      //   { id: 2, title: '消息2' },
      //   // 更多消息
      // ],
      activeTab: 'all',
    };
  },

  computed: {
    // 过滤后的消息列表
    filteredChats() {
      switch (this.activeTag) {
        case 'friend':
          return this.chats.filter(chat => chat.tag === 'friend');
        case 'group':
          return this.chats.filter(chat => chat.tag === 'group');
        case 'unread':
          return this.chats.filter(chat => chat.unreadCount > 0);
        default:
          return this.chats;
      }
    },
  },

  methods: {
    // 选中tag筛选消息
    filterChats(tagName) {
      this.activeTag = tagName;
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
  padding: 20px; /* 设置内边距 */
}
</style>