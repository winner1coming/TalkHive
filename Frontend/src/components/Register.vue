<template>
  <!-- 注册页面容器 -->
  <div class="register">
    <!-- 页面标题 -->
    <h2>用户注册</h2>
    
    <!-- ID 输入框 -->
    <div class="input-group">
      <label for="id">ID:</label>
      <input id="id" type="text" v-model="id" placeholder="ID" @blur="validateId" />
    </div>
    <p v-if="errors.id" class="error">{{ errors.id }}</p>
    
    <!-- 昵称输入框 -->
    <div class="input-group">
      <label for="nickname">昵称:</label>
      <input id="nickname" type="text" v-model="nickname" placeholder="昵称" @blur="validateNickname" />
    </div>
    <p v-if="errors.nickname" class="error">{{ errors.nickname }}</p>
    
    <!-- 手机号输入框 -->
    <div class="input-group">
      <label for="phoneNumber">手机号:</label>
      <input id="phoneNumber" type="text" v-model="phoneNumber" placeholder="手机号" @blur="validatePhoneNumber" />
    </div>
    <p v-if="errors.phoneNumber" class="error">{{ errors.phoneNumber }}</p>
    
    <!-- 密码输入框 -->
    <div class="input-group">
      <label for="password">密码:</label>
      <input id="password" type="password" v-model="password" placeholder="密码" @blur="validatePassword" />
    </div>
    <p v-if="errors.password" class="error">{{ errors.password }}</p>
    
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
      <button class="send-verification-code" @click="sendVerificationCode">获取</button>
    </div>
    <p v-if="errors.verificationCode" class="error">{{ errors.verificationCode }}</p>
    
    <!-- 注册按钮 -->
    <button class="register-button" @click="register">注册</button>
  </div>
</template>

<script>
import { register, sendVerificationCode } from '@/services/api'; // 导入注册和发送验证码 API

export default {
  data() {
    return {
      id: '',
      nickname: '',
      phoneNumber: '',
      password: '',
      confirmPassword: '',
      verificationCode: '',
      errors: {
        id: '',
        nickname: '',
        phoneNumber: '',
        password: '',
        confirmPassword: '',
        verificationCode: '',
      },
    };
  },
  
  methods: {
    // 验证 ID
    validateId() {
      if (!this.id) {
        this.errors.id = 'ID 不能为空';
      } else {
        this.errors.id = '';
      }
    },
    
    // 验证昵称
    validateNickname() {
      if (!this.nickname) {
        this.errors.nickname = '昵称不能为空';
      } else {
        this.errors.nickname = '';
      }
    },
    
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
    
    // 验证密码
    validatePassword() {
      if (!this.password) {
        this.errors.password = '密码不能为空';
      } else if (this.password.length < 6) {
        this.errors.password = '密码长度不能少于6位';
      } else {
        this.errors.password = '';
      }
    },
    
    // 验证确认密码
    validateConfirmPassword() {
      if (!this.confirmPassword) {
        this.errors.confirmPassword = '确认密码不能为空';
      } else if (this.confirmPassword !== this.password) {
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
    
    // 注册方法
    async register() {
      this.validateId();
      this.validateNickname();
      this.validatePhoneNumber();
      this.validatePassword();
      this.validateConfirmPassword();
      this.validateVerificationCode();
      
      if (Object.values(this.errors).some(error => error)) {
        return;
      }
      
      try {
        const response = await register({
          id: this.id,
          nickname: this.nickname,
          phoneNumber: this.phoneNumber,
          password: this.password,
          verificationCode: this.verificationCode,
        });
        
        if (response.success) {
          this.$router.push('/login');
        } else {
          alert(response.message || '注册失败');
        }
      } catch (error) {
        alert(error || '注册失败');
      }
    },
  },
};
</script>

<style scoped>
.register {
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
  margin-left: 80px;
  font-size: 24px;
  color: #333;
}

.input-group {
  display: flex;
  align-items: center;
  margin-bottom:15px;
  width: 100;
  max-width: 300px;
  box-sizing: border-box;
}

.input-group label {
  width: 100%;
  margin-right: 5px;
  font-size: 14px;
  color: #666;
  white-space: nowrap; 
  text-align: left;
}

.input-group input {
  flex: 1;
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 4px;
  margin-left: 20px;
  text-align: left;
  width:300px;
  height: 20px;
}

.send-verification-code {
  padding: 15px;
  font-size: 12px;
  color: #fff;
  background-color: #42b983;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-left:40px;
  height: 40px;
  width: 400px;
}

.error {
  color: red;
  font-size: 12px;
  margin-top: 5px;
  margin-left: 110px;
}

.register-button {
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

.register-button:hover {
  background-color: #369f6e;
}
</style>