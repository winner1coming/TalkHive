<template>
    <!-- 找回密码页面容器 -->
    <div class="forgotpassword">
      <!-- 页面标题 -->
      <h2>找回密码</h2>
      
      <!-- 手机号输入框 -->
      <div class="input-group">
        <label for="phoneNumber">手机号:</label>
        <input id="phoneNumber" type="text" v-model="phoneNumber" placeholder="手机号" @blur="validatePhoneNumber" />
      </div>
      <p v-if="errors.phoneNumber" class="error">{{ errors.phoneNumber }}</p>
      
      <!-- 新密码输入框 -->
      <div class="input-group">
        <label for="newPassword">新密码:</label>
        <input id="newPassword" type="password" v-model="newPassword" placeholder="新密码" @blur="validateNewPassword" />
      </div>
      <p v-if="errors.newPassword" class="error">{{ errors.newPassword }}</p>
      
      <!-- 确认密码输入框 -->
      <div class="input-group">
        <label for="confirmPassword">确认密码:</label>
        <input id="confirmPassword" type="password" v-model="confirmPassword" placeholder="确认密码" @blur="validateConfirmPassword" />
      </div>
      <p v-if="errors.confirmPassword" class="error">{{ errors.confirmPassword }}</p>
      
      <!-- 验证码输入框 -->
      <div class="input-group">
        <label for="verificationCode">验证码:</label>
        <input id="verificationCode" type="text" v-model="verificationCode" placeholder="验证码" @blur="validateVerificationCode" />
        <button class="send-verification-code" @click="sendVerificationCode">获取验证码</button>
      </div>
      <p v-if="errors.verificationCode" class="error">{{ errors.verificationCode }}</p>
      
      <!-- 提交按钮 -->
      <button class="submit-button" @click="submit">提交</button>
      
      <!-- 找回密码成功提示 -->
      <div v-if="successMessage" class="success-message">
        <p>{{ successMessage }}</p>
        <button class="confirm-button" @click="goToLogin">确定</button>
      </div>
    </div>
  </template>
  
  <script>
  import { sendVerificationCode, resetPassword } from '@/services/api'; // 导入发送验证码和重置密码 API
  
  export default {
    data() {
      return {
        phoneNumber: '',
        newPassword: '',
        confirmPassword: '',
        verificationCode: '',
        errors: {
          phoneNumber: '',
          newPassword: '',
          confirmPassword: '',
          verificationCode: '',
        },
        successMessage: '',
      };
    },
    
    methods: {
      // 验证手机号
      validatePhoneNumber() {
        if (!this.phoneNumber) {
          this.errors.phoneNumber = '手机号不能为空';
        } else if (!/^1[3-9]\d{9}$/.test(this.phoneNumber)) {
          this.errors.phoneNumber = '手机号格式不正确';
        } else {
          this.errors.phoneNumber = '';
        }
      },
      
      // 验证新密码
      validateNewPassword() {
        if (!this.newPassword) {
          this.errors.newPassword = '新密码不能为空';
        } else if (this.newPassword.length < 6) {
          this.errors.newPassword = '密码长度不能少于6位';
        } else {
          this.errors.newPassword = '';
        }
      },
      
      // 验证确认密码
      validateConfirmPassword() {
        if (!this.confirmPassword) {
          this.errors.confirmPassword = '确认密码不能为空';
        } else if (this.confirmPassword !== this.newPassword) {
          this.errors.confirmPassword = '两次输入的密码不一致';
        } else {
          this.errors.confirmPassword = '';
        }
      },
      
      // 验证验证码
      validateVerificationCode() {
        if (!this.verificationCode) {
          this.errors.verificationCode = '验证码不能为空';
        } else {
          this.errors.verificationCode = '';
        }
      },
      
      // 发送验证码
      async sendVerificationCode() {
        if (!this.validatePhoneNumber()) {
          return;
        }
        
        try {
          const response = await sendVerificationCode(this.phoneNumber);
          if (response.success) {
            alert('验证码已发送');
          } else {
            alert(response.message || '发送验证码失败');
          }
        } catch (error) {
          alert(error || '发送验证码失败');
        }
      },
      
      // 提交找回密码请求
      async submit() {
        this.validatePhoneNumber();
        this.validateNewPassword();
        this.validateConfirmPassword();
        this.validateVerificationCode();
        
        if (Object.values(this.errors).some(error => error)) {
          return;
        }
        
        try {
          const response = await resetPassword({
            phoneNumber: this.phoneNumber,
            newPassword: this.newPassword,
            verificationCode: this.verificationCode,
          });
          
          if (response.success) {
            this.successMessage = '找回密码成功，请返回重新登录';
          } else {
            alert(response.message || '找回密码失败');
          }
        } catch (error) {
          alert(error || '找回密码失败');
        }
      },
      
      // 跳转到登录页面
      goToLogin() {
        this.$router.push('/login');
      },
    },
  };
  </script>
  
  <style scoped>
  .forgotpassword {
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
  
  .input-group {
    display: flex;
    align-items: center;
    margin-bottom: 15px;
    width: 100%;
    max-width: 300px;
    box-sizing: border-box;
  }
  
  .input-group label {
    width: 100px;
    margin-right: 10px;
    font-size: 14px;
    color: #666;
    text-align: right;
  }
  
  .input-group input {
    flex: 1;
    padding: 10px;
    font-size: 16px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  
  .send-verification-code {
    padding: 10px;
    font-size: 14px;
    color: #fff;
    background-color: #42b983;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    margin-left: 10px;
  }
  
  .error {
    color: red;
    font-size: 12px;
    margin-top: 5px;
    margin-left: 110px;
  }
  
  .submit-button {
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
  
  .submit-button:hover {
    background-color: #369f6e;
  }
  
  .success-message {
    margin-top: 20px;
    text-align: center;
  }
  
  .success-message p {
    font-size: 16px;
    color: #333;
    margin-bottom: 10px;
  }
  
  .confirm-button {
    padding: 10px;
    font-size: 16px;
    color: #fff;
    background-color: #42b983;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .confirm-button:hover {
    background-color: #369f6e;
  }
  </style>