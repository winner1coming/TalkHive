<template>
    <div class="loginth">
      <!-- 页面标题 -->
      <h2>登录界面</h2>
      <!-- 登录选项 -->
        <div class="login-options">
          <button :class="{ active: loginType === 'account' }" @click="loginType = 'account'">账号登录</button>
          <button :class="{ active: loginType === 'sms' }" @click="loginType = 'sms'">短信登录</button>
        </div>
        <div class = "select">
        <!-- 账号登录表单 -->
        <AccountLogin v-if="loginType === 'account'" />
        
        <!-- 短信登录表单 -->
        <SmsLogin v-if="loginType === 'sms'" />
        </div>
    </div>
  </template>
  
  <script>
  import AccountLogin from '@/components/AccountLogin.vue';
  import SmsLogin from '@/components/SmsLogin.vue';
  
  export default {
    components: { AccountLogin, SmsLogin},
    data() {
      return {
        loginType: 'account', // 登录类型：account 或 sms
      };
    },
  };
  </script>
  
  <style scoped>
  .login-th {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;
    background-color: #f9f9f9;
    padding: 20px;
    box-sizing: border-box;
  }
  
  h2 {
    margin-bottom: 20px;
    font-size: 24px;
    color: #333;
  }

  .avatar img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}
  
  .login-options {
    display: flex;
    margin-bottom: 20px;
  }
  
  .login-options button {
    padding: 10px 20px;
    font-size: 16px;
    color: #666;
    background-color: #fff;
    border: 1px solid #ccc;
    border-radius: 4px;
    cursor: pointer;
    margin-right: 10px;
  }
  
  .login-options button.active {
    background-color: #42b983;
    color: #fff;
    border-color: #42b983;
  }

  .select {
  width: 100%;
  max-width: 400px;
  margin-bottom: 20px;
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