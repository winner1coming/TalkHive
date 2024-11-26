<template>
  <div class="accountlogin">
      <h2> Login</h2>
    <div class="input-group">
      <label for="account">账号:</label>
      <input id="account" type="text" v-model="account" placeholder="请输入账号" />
    </div>
    <div class="input-group">
      <label for="password">密码:</label>
      <input id="password" type="password" v-model="password" placeholder="请输入密码" />
    </div>
    <button class="login-button" @click="login">登录</button>
      </div>
</template>

<script>
import { login } from '@/services/api'; // 导入登录 API

export default {
  data() {
    return {
      account: '',
      password: '',
    };
  },
  
  methods: {
    // 登录方法
    async login() {
      try {
        const response = await login(this.account, this.password);
        if (response.success) {
          this.$router.push('/');
        } else {
          alert(response.message || '登录失败');
        }
      } catch (error) {
        alert(error || '登录失败');
      }
    },
  },
};
</script>

<style scoped>
.accountlogin{
    /* 使用 Flexbox 布局，使内容居中 */
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh; /* 设置高度为视口高度 */
    background-color:rgb(236, 245, 245);
    padding: 20px;
    box-sizing: border-box;
}

.accountlogin h2{
    font-size:36px;
    color: black;
    margin-bottom: 20px;
    margin-left: 20px;
}

.input-group {
  display: flex;
  margin-bottom: 20px;
  width: 100%;
  max-width: 300px;
  box-sizing: border-box;
}

.input-group label {
  margin-right: 15px;
  font-size: 14px;
  color: #666;
  white-space: nowrap; /* 防止标签换行 */
}

.input-group input {
  flex: 1; /* 使输入框占据剩余空间 */
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.login-button {
  width: 100%;
  max-width: 100px;
  padding: 10px;
  font-size: 16px;
  color: #fff;
  background-color: #42b983;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  margin-top: 10px;
  margin-left: 20px;
}

.login-button:hover {
  background-color: #369f6e;
}

p {
  margin-top: 10px;
  font-size: 14px;
  color: #666;
}

p a {
  color: #42b983;
  text-decoration: none;
}

p a:hover {
  text-decoration: underline;
}
</style>