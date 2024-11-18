<template>
  <!-- 编辑资料页面容器 -->
  <div class="edit-profile">
    <!-- 页面标题 -->
    <h2>编辑资料</h2>
    
    <!-- 用户名输入框，使用 v-model 双向绑定到 data 中的 username 属性 -->
    <input type="text" v-model="username" placeholder="用户名" />
    
    <!-- ID 输入框，使用 v-model 双向绑定到 data 中的 id 属性 -->
    <input type="text" v-model="id" placeholder="ID" />
    
    <!-- 保存按钮，点击时触发 saveProfile 方法 -->
    <button @click="saveProfile">保存</button>
  </div>
</template>

<script>
export default {
  // 组件的 data 函数，返回一个对象，包含组件的响应式数据
  data() {
    return {
      // 用户名输入框的值
      username: '',
      // ID 输入框的值
      id: '',
    };
  },
  
  // 组件的方法定义
  methods: {
    // 保存资料方法，处理用户点击保存按钮时的逻辑
    saveProfile() {
      // 调用 Vuex store 中的 updateProfile 动作，传递用户名和 ID
      this.$store.dispatch('updateProfile', { username: this.username, id: this.id });
    },
  },
};
</script>

<style scoped>
/* 编辑资料页面的样式 */
.edit-profile {
  padding: 20px; /* 设置内边距 */
}
</style>