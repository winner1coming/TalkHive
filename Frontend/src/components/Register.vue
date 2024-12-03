<template>
  <!-- 注册页面容器 -->
  <div class="register">
    <div class="container">
    <!-- 页面标题 -->
    <h2>注册</h2>

    <div class="avatar-container">
      <img class="headavatar" :src="avatar" @click="triggerFileInput" />
      <input type="file" ref="fileInput" style="display: none;" @change="handleFileChange" accept="image/*" />
    </div>

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
      <button class="send-verification-code" @click="sendSmsCode">获取</button>
    </div>
    <p v-if="errors.verificationCode" class="error">{{ errors.verificationCode }}</p>
    
    <!-- 注册按钮 -->
    <button class="register-button" @click="register">注册</button>
    </div>
 </div>
</template>

<script>
import { Register, sendSmsCode } from '@/services/api'; // 导入注册和发送验证码 API

export default {
  data() {
    return {
      avatar:'@/assets/images/avatar.jpg',
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
    // 触发文件输入
    triggerFileInput() {
      this.$refs.fileInput.click();
    },

    // 处理文件选择
    handleFileChange(event) {
      const file = event.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          this.avatar = e.target.result;
        };
        reader.readAsDataURL(file);
      }
    },

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
    async sendSmsCode() {
      if (!this.validatePhoneNumber()) {
        return;
      }
      
      try {
        const response = await sendSmsCode({
          command:'register',
          phoneNumber: this.phoneNumber,
        });
        if (response.success) {
          alert('验证码已发送');
        } else {
          alert(response.message || '发送验证码失败');
        }
      } catch (error) {
        alert(error || '发送验证码失败');
      }
    },

    //验证发送过来的验证码是否正确

    
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
        const response = await Register({
          avatar :this.avatar,
          id: this.id,
          nickname: this.nickname,
          phoneNumber: this.phoneNumber,
          password: this.password,
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
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
  align-items: center;
  justify-items: center;
  height: 100vh;
  background-color: #f9f9f9;
  padding: 20px;
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
    margin-top: 40px;
    width: 400px;
  }

h2 {
  margin-bottom: 10px;
  align-items: center;
  font-size: 24px;
  color: #333;
}

.avatar-container {
  position: relative;
  cursor: pointer;
}

.headavatar {
  width: 100px;
  height: 100px;
  margin-top: 10px;
  margin-bottom: 20px;
  margin-left: 10px;
  border-radius: 100%;
  cursor: pointer;
}

.input-group {
  display: flex;
  grid-template-columns: 100px 1fr;
  gap: 10px;
  align-items:center;
  width: 100%;
  max-width: 400px;
  margin-bottom: 20px;
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
  flex:1;
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