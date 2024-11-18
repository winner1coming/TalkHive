<template>
  <!-- 注册页面容器 -->
  <div class="register">
    <!-- 页面标题 -->
    <h2>注册</h2>
    
    <!-- 用户名输入框，使用 v-model 双向绑定到 data 中的 username 属性 -->
    <input type="text" v-model="username" placeholder="用户名" />
    
    <!-- 密码输入框，使用 v-model 双向绑定到 data 中的 password 属性 -->
    <input type="password" v-model="password" placeholder="密码" />
    
    <!-- 确认密码输入框，使用 v-model 双向绑定到 data 中的 confirmPassword 属性 -->
    <input type="password" v-model="confirmPassword" placeholder="确认密码" />
    
    <!-- 注册按钮，点击时触发 register 方法 -->
    <button @click="register">注册</button>
    
    <!-- 提示已有账号的用户可以点击链接跳转到登录页面 -->
    <p>已有账号？<router-link to="/login">登录</router-link></p>
  </div>
</template>

<script>
export default {
  // 组件的 data 函数，返回一个对象，包含组件的响应式数据
  data() {
    return {
      // 用户名输入框的值
      username: '',
      // 密码输入框的值
      password: '',
      // 确认密码输入框的值
      confirmPassword: '',
    };
  },
  
  // 组件的方法定义
  methods: {
    // 注册方法，处理用户点击注册按钮时的逻辑
    register() {
      // 检查密码和确认密码是否一致
      if (this.password !== this.confirmPassword) {
        // 如果不一致，弹出提示框并返回
        alert('密码不一致');
        return;
      }
      
      // 如果密码一致，调用 Vuex store 中的 register 动作，传递用户名和密码
      this.$store.dispatch('register', { username: this.username, password: this.password });
    },
  },
};
</script>

<style scoped>
/* 注册页面的样式 */
.register {
  /* 使用 Flexbox 布局，使内容居中 */
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh; /* 设置高度为视口高度 */
}
</style>