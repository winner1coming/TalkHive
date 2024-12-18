<template>
  <div class="login">
    <div class="loginth">
      <!-- 登录选项 -->
        <div class="login-options">
          <button :class="{ active: loginType === 'account' }" @click="loginType = 'account'">账号登录</button>
          <button :class="{ active: loginType === 'sms' }" @click="loginType = 'sms'">短信登录</button>
        </div>
        <div class="select" v-if="loginType === 'account'">
          <AccountLogin />
        </div>
    
        <!-- 短信登录表单 -->
        <div class="select" v-if="loginType === 'sms'">
          <SmsLogin />
        </div>
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

.loginth {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 120vh;
  background-color: #f9f9f9;
  padding: 20px;
  box-sizing: border-box;
}

.login-options {
  display: flex;
  margin-bottom: 40px;
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
  margin-top: -10px;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* 添加阴影效果 */
  border-radius: 8px; /* 添加圆角效果 */
  background-color: #fff; /* 添加背景色 */
  padding: 20px; /* 添加内边距 */
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