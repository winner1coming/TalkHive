<template>
  <!-- 消息列表页面容器 -->
  <div class="message-list">
    <!-- 页面标题 -->
    <h2>消息列表</h2>
    
    <!-- 消息列表，使用 v-for 指令循环渲染 messages 数组中的每个消息 -->
    <ul>
      <!-- 每个消息项 -->
      <li v-for="message in messages" :key="message.id">
        <!-- 使用 router-link 组件生成跳转到聊天页面的链接，链接路径为 /chat/{消息ID} -->
        <router-link :to="`/chat/${message.id}`">{{ message.title }}</router-link>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  // 组件的 data 函数，返回一个对象，包含组件的响应式数据
  data() {
    return {
      // 消息列表数据，每个消息包含 id 和 title 属性
      messages: [
        { id: 1, title: '消息1' },
        { id: 2, title: '消息2' },
        // 更多消息
      ],
    };
  },
};
</script>

<style scoped>
/* 消息列表页面的样式 */
.message-list {
  padding: 20px; /* 设置内边距 */
}
</style>