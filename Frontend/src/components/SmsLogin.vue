<template>
    <div class="smslogin">
      <div class="input-group">
        <label for="phoneNumber">手机号</label>
        <input id="phoneNumber" type="text" v-model="phoneNumber" placeholder="请输入手机号" />
      </div>
      <div class="input-group">
        <label for="smsCode">验证码</label>
        <input id="smsCode" type="text" v-model="smsCode" placeholder="请输入验证码" />
        <button class="send-sms-button" @click="sendSmsCode">发送验证码</button>
      </div>
      <button class="login-button" @click="login">登录</button>
    </div>
  </template>
  
  <script>
  import { login, sendSmsCode } from '@/services/api'; // 导入登录和发送验证码 API
  
  export default {
    data() {
      return {
        phoneNumber: '',
        smsCode: '',
      };
    },
    
    methods: {
      // 登录方法
      async login() {
        try {
          const response = await login(this.phoneNumber, this.smsCode);
          if (response.success) {
            this.$router.push('/');
          } else {
            alert(response.message || '登录失败');
          }
        } catch (error) {
          alert(error || '登录失败');
        }
      },
      
      // 发送验证码方法
      async sendSmsCode() {
        if (!this.validatePhoneNumber(this.phoneNumber)) {
          alert('请输入有效的手机号码');
          return;
        }
        
        try {
          const response = await sendSmsCode(this.phoneNumber);
          if (response.success) {
            alert('验证码已发送');
          } else {
            alert(response.message || '发送验证码失败');
          }
        } catch (error) {
          alert(error || '发送验证码失败');
        }
      },
      
      // 验证手机号
      validatePhoneNumber(phoneNumber) {
        return /^1[3-9]\d{9}$/.test(phoneNumber);
      },
    },
  };
  </script>
  
  <style scoped>
  .input-group {
    display: flex;
    align-items: center;
    margin-bottom: 15px;
    width: 100%;
    max-width: 300px;
  }
  
  .input-group label {
    margin-right: 10px;
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
  
  .send-sms-button {
    padding: 10px;
    font-size: 14px;
    color: #fff;
    background-color: #42b983;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    margin-left: 10px;
  }
  
  .login-button {
    width: 100%;
    max-width: 300px;
    padding: 10px;
    font-size: 16px;
    color: #fff;
    background-color: #42b983;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    margin-top: 10px;
  }
  
  .login-button:hover {
    background-color: #369f6e;
  }
  </style>