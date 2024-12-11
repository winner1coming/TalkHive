<template>
  <div class="accountlogin">
    <div class = "all">
      <img  class="avatar" src = '@/assets/images/avatar.jpg'/>
      
      <div class="input-group">
        <label for="account">账号:</label>
        <input id="account" type="text" v-model="phone" placeholder="请输入账号" />
      </div>
      
      <div class="input-group">
        <label for="password">密码:</label>
        <input id="password" type="password" v-model="password" placeholder="请输入密码" />
      </div>

      <button class="login-button" @click="testlogin">登录</button>
      
      <div class = "link">
        <!-- 注册链接 -->
        <router-link to = "/register" class = "register">注册</router-link>
        <router-link to = "/forgetpassword" class ="reset">忘记密码</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { login } from '@/services/api'; // 导入登录 API

export default {
  data() {
    return {
      phone: '666666',
      password: '666666',
    };
  },
  
  methods: {
    // 登录方法
    async testlogin(){
      this.$router.push('/home');
    },
    async login() {
      try {
        const response = await login(
          {
            phone: this.phone,
            password: this.password,
          });
        if (response.success) {
          this.$router.push('/home');
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
    height: 60vh; /* 设置高度为视口高度 */
    padding: 10px;
    box-sizing: border-box;
}

.all {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 350px;
  padding: 20px;
  box-sizing: border-box;
}


.avatar{
  width: 100px;
  height: 100px;
  margin-top: 10px;
  margin-bottom: 20px;
  margin-left: 45px;
  border-radius: 100%;
}

.input-group {
  display: flex;
  margin-left: 20px;
  margin-bottom: 20px;
  align-items: center;
  width: 100%;
  max-width: 300px;
  box-sizing: border-box;
}

.input-group label {
  margin-right: 15px;
  font-size: 14px;
  align-items: center;
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

.link {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  width: 100%;
  margin-top: 20px;
}

.register {
  text-decoration: none;
  justify-items: start;
  color: #42b983;
  font-size: 14px;
}

.reset {
  text-decoration: none;
  color: #42b983;
  font-size: 14px;
}

.register:hover, .reset:hover {
  text-decoration: underline;
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