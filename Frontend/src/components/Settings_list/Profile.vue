<template>
  <!-- 个人信息页面容器 -->
  <div class="profile">
    <!-- 页面标题 -->
    <h2>个人信息</h2>
    
    <!-- 显示用户名 -->
    <p>用户名: {{ user.username }}</p>
    
    <!-- 显示用户ID -->
    <p>ID: {{ user.id }}</p>
    
    <!-- 编辑资料链接，点击后跳转到编辑资料页面 -->
    <router-link to="/edit-profile">编辑资料</router-link>
  </div>
</template>

<script>
export default {
  // 计算属性，用于从 Vuex store 中获取当前用户信息
  computed: {
    user() {
      // 返回 Vuex store 中的 user 状态
      return this.$store.state.user;
    },
  },
};
</script>

<style scoped>
/* 个人信息页面的样式 */
.profile {
  padding: 20px; /* 设置内边距 */
}
</style>