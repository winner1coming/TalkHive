<template>
  <div class="login">
    <h2>登录</h2>
    
    <!-- 手机号输入框 -->
    <div class="input-group">
      <label for="phoneNumber">手机号</label>
      <input id="phoneNumber" type="text" v-model="phoneNumber" placeholder="请输入手机号" />
    </div>
    
    <!-- 密码输入框 -->
    <div class="input-group">
      <label for="password">密码</label>
      <input id="password" :type="passwordFieldType" v-model="password" placeholder="请输入密码" />
      <button class="toggle-password" @click="togglePasswordVisibility">
        {{ passwordFieldType === 'password' ? '显示' : '隐藏' }}
      </button>
    </div>
    
    <!-- 记住密码复选框 -->
    <div class="input-group checkbox">
      <label>
        <input type="checkbox" v-model="rememberPassword" /> 记住密码
      </label>
    </div>
    
    <!-- 登录按钮 -->
    <button class="login-button" @click="login">登录</button>
    
    <!-- 找回密码链接 -->
    <p><router-link to="/forgot-password">忘记密码？</router-link></p>
    
    <!-- 注册链接 -->
    <p>没有账号？<router-link to="/register">注册</router-link></p>
  </div>
</template>

<script>
export default {
  data() {
    return {
      phoneNumber: '',
      password: '',
      rememberPassword: false,
      passwordFieldType: 'password',
    };
  },
  
  methods: {
    // 登录方法
    login() {
      if (!this.validatePhoneNumber(this.phoneNumber)) {
        alert('请输入有效的手机号码');
        return;
      }
      
      if (!this.password) {
        alert('请输入密码');
        return;
      }
      
      // 调用 Vuex store 中的 login 动作
      this.$store.dispatch('login', {
        phoneNumber: this.phoneNumber,
        password: this.password,
        rememberPassword: this.rememberPassword,
      }).then(() => {
        // 登录成功后的逻辑
        this.$router.push('/');
      }).catch(error => {
        alert(error.message);
      });
    },
    
    // 验证手机号
    validatePhoneNumber(phoneNumber) {
      // 简单的手机号验证，可以根据需求扩展
      return /^1[3-9]\d{9}$/.test(phoneNumber);
    },
    
    // 切换密码可见性
    togglePasswordVisibility() {
      this.passwordFieldType = this.passwordFieldType === 'password' ? 'text' : 'password';
    },
  },
};
</script>

<style scoped>
.login {
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
  flex-direction: column;
  margin-bottom: 15px;
  width: 100%;
  max-width: 300px;
}

.input-group label {
  margin-bottom: 5px;
  font-size: 14px;
  color: #666;
}

.input-group input {
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.input-group.checkbox {
  flex-direction: row;
  align-items: center;
}

.input-group.checkbox label {
  margin-left: 10px;
}

.toggle-password {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
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