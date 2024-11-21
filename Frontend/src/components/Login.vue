<template>
  <!-- 登录组件的根元素 -->
  <div class="login">
    <!-- 标题 -->
    <h2>登录</h2>

    <!-- 用户名输入框 -->
    <input type="text" v-model="username" placeholder="用户名" />

    <!-- 密码输入框 -->
    <input type="password" v-model="password" placeholder="密码" />

    <!-- 登录按钮 -->
    <button @click="login">登录</button>

    <!-- 注册链接 -->
    <p>没有账号？<router-link to="/register">注册</router-link></p>
  </div>
</template>

<script>
export default {
  // 组件名称
  name: 'Login',

  // 组件的数据
  data() {
    return {
      // 用户名
      username: '',

      // 密码
      password: '',
    };
  },

  // 组件的方法
  methods: {
    // 登录方法
    login() {
      // 分发登录动作到 Vuex store
      this.$store.dispatch('login', { username: this.username, password: this.password });
    },
  },
};
</script>

<style scoped>
/* 登录组件的样式 */
.login {
  /* 使用 Flexbox 布局 */
  display: flex;

  /* 设置主轴方向为垂直 */
  flex-direction: column;

  /* 水平居中对齐 */
  align-items: center;

  /* 垂直居中对齐 */
  justify-content: center;

  /* 设置高度为视口高度 */
  height: 100vh;
}
</style>