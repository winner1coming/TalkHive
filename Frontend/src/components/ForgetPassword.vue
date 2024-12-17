<template>
    <!-- 找回密码页面容器 -->
    <div class="forgotpassword">
      <div class="container">
      <!-- 页面标题 -->
      <h2>找回密码</h2>
      
      <!-- 手机号输入框 -->
      <div class="input-group">
        <label for="email">邮箱:</label>
        <input id="email" type="text" v-model="email" placeholder="邮箱" @blur="validateEmail" />
      </div>
      <p v-if="errors.email" class="error">{{ errors.email }}</p>
      
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
        <button class="send-verification-code" @click="sendSmsCode" :disabled="isCountingDown" :class="{ 'counting-down': isCountingDown }">
          {{ isCountingDown ? `${countdown}秒后重试` : '获取' }}</button>
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
    </div>
  </template>
  
  <script>
  import { sendSmsCode, resetPassword } from '@/services/api'; // 导入发送验证码和重置密码 API
  
  export default {
    data() {
      return {
        email: '',
        newPassword: '',
        confirmPassword: '',
        verificationCode: '',
        errors: {
          email: '',
          newPassword: '',
          confirmPassword: '',
          verificationCode: '',
        },
        successMessage: '',
        Code:'',
        isCountingDown:false,
        countdown:60,
      };
    },
    
    methods: {
      // 验证邮箱
      validateEmail() {
        if (!this.email) {
          this.errors.email = '邮箱不能为空';
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(this.email)) {
          this.errors.email = '邮箱格式不正确';
        } else {
          this.errors.email = '';
        }
      },
      
      // 验证新密码
      validateNewPassword() {
          if (!this.newPassword) {
            this.errors.newPassword = '新密码不能为空';
          } else if (this.newPassword.length < 6) {
            this.errors.newPassword = '密码长度不能少于6位';
          } else if (/\s/.test(this.newPassword)) {
          this.errors.newPassword = '密码不能包含空格';
          } else if (!/^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{6,}$/.test(this.newPassword)) {
            this.errors.newPassword = '密码只能是包含数字和字母的组合';
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

      // 启动倒计时
      startCountdown() {
        this.isCountingDown = true;
        this.countdown = 60;

        const timer = setInterval(() => {
          this.countdown--;
          if (this.countdown <= 0) {
            clearInterval(timer);
            this.isCountingDown = false;
          }
        }, 1000);
      },

      async validateCode(){
        if(Code){
          if(Code !== this.verificationCode){
              alert('验证码错误');
              return;
          }
        }
        else{
          alter('请先获取验证码！');
          return;
        }
      },
      
      // 发送验证码
      async sendSmsCode() {
        if (!this.validateEmail()) {
          return;
        }
        
        try {
          const response = await sendSmsCode(
            {
              command:'resetPassword',
              email:this.email,
            }
          );
          if (response.success) {
            alert('验证码已发送');
            this.Code = response.code;
            this.startCountdown();
          } else {
            alert(response.message || '发送验证码失败');
          }
        } catch (error) {
          alert(error || '发送验证码失败');
        }
      },
      
      // 提交找回密码请求
      async submit() {
        this.validateEmail();
        this.validateNewPassword();
        this.validateConfirmPassword();
        this.validateVerificationCode();
        
        if (Object.values(this.errors).some(error => error)) {
          return;
        }
        this.validateCode();
        
        try {
          const response = await resetPassword({
            email: this.email,
            password: this.newPassword,
          });
          
          if (response.success) {
            this.successMessage = '找回密码成功，请返回重新登录';
            this.goToLogin();
          } else {
            alert(response.message || '找回密码失败');
          }
        } catch (error) {
          alert(error || '找回密码失败');
        }
      },
      
      // 跳转到登录页面
      goToLogin() {
        this.$router.push('/loginth');
      },
    },
  };
  </script>
  
  <style scoped>
  .forgotpassword {
    display: grid;
    grid-template-columns: 1fr;
    gap: 20px;
    align-items: center;
    justify-items: center;
    height: 80vh;
    padding: 10px;
    box-sizing: border-box;
  }

  .container {
    display:flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* 添加阴影效果 */
    border-radius: 8px; /* 添加圆角效果 */
    background-color: #fff; /* 添加背景色 */
    padding: 20px; /* 添加内边距 */
    margin-top: 120px;
  }

  h2 {
    margin-bottom: 30px;
    align-items: center;
    font-size: 24px;
    color: #333;
  }

  .input-group {
    display: flex;
    grid-template-columns: 100px 1fr;
    gap: 10px;
    align-items: center;
    margin-bottom: 20px;
    width: 100%;
    max-width: 400px;
  }

  .input-group label {
    font-size: 14px;
    color: #666;
    text-align: left;
    width: 89px;
  }

  .input-group input {
    padding: 10px;
    font-size: 16px;
    border: 1px solid #ccc;
    border-radius: 4px;
    width: 100%;
    flex: 1;
  }

  .send-verification-code {
    padding: 10px;
    width: 20%;
    font-size: 14px;
    color: #fff;
    background-color: #42b983;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    margin-left: 0px;
  }

  .error {
    color: red;
    font-size: 12px;
    margin-top: 5px;
    text-align: left;
    width: 100%;
    max-width: 400px;
  }

  .submit-button {
    width: 100%;
    max-width: 100px;
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