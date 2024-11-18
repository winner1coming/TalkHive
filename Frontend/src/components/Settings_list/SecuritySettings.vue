<template>
  <!-- 安全设置页面容器 -->
  <div class="security-settings">
    <!-- 页面标题 -->
    <h2>安全设置</h2>
    
    <!-- 旧密码输入框，使用 v-model 双向绑定到 data 中的 oldPassword 属性 -->
    <input type="password" v-model="oldPassword" placeholder="旧密码" />
    
    <!-- 新密码输入框，使用 v-model 双向绑定到 data 中的 newPassword 属性 -->
    <input type="password" v-model="newPassword" placeholder="新密码" />
    
    <!-- 确认新密码输入框，使用 v-model 双向绑定到 data 中的 confirmPassword 属性 -->
    <input type="password" v-model="confirmPassword" placeholder="确认新密码" />
    
    <!-- 修改密码按钮，点击时触发 changePassword 方法 -->
    <button @click="changePassword">修改密码</button>
  </div>
</template>

<script>
export default {
  // 组件的 data 函数，返回一个对象，包含组件的响应式数据
  data() {
    return {
      // 旧密码输入框的值
      oldPassword: '',
      // 新密码输入框的值
      newPassword: '',
      // 确认新密码输入框的值
      confirmPassword: '',
    };
  },
  
  // 组件的方法定义
  methods: {
    // 修改密码方法，处理用户点击修改密码按钮时的逻辑
    changePassword() {
      // 检查新密码和确认新密码是否一致
      if (this.newPassword !== this.confirmPassword) {
        // 如果不一致，弹出提示框并返回
        alert('新密码不一致');
        return;
      }
      
      // 如果新密码一致，调用 Vuex store 中的 changePassword 动作，传递旧密码和新密码
      this.$store.dispatch('changePassword', { oldPassword: this.oldPassword, newPassword: this.newPassword });
    },
  },
};
</script>

<style scoped>
/* 安全设置页面的样式 */
.security-settings {
  padding: 20px; /* 设置内边距 */
}
</style>