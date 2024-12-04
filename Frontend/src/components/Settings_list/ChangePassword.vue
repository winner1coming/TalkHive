<template>
    <div class="change-password">
      <h3>更改密码</h3>
      <div class="text_group">
        <label for="oldPassword">原密码:</label>
        <input type="password" v-model="oldPassword" :placeholder="oldPassword" />
      </div>
      <div class="text_group">
        <label for="newPassword">新密码:</label>
        <div class="password-input">
          <input :type="showNewPassword ? 'text' : 'password'" v-model="newPassword" :placeholder="newPassword" />
          <span class="eye-icon" @click="toggleNewPasswordVisibility">
            <i :class="showNewPassword ? 'fas fa-eye' : 'fas fa-eye-slash'"></i>
          </span>
        </div>
      </div>
      <div class="text_group">
        <label for="confirmPassword">确认密码:</label>
        <div class="password-input">
          <input :type="showConfirmPassword ? 'text' : 'password'" v-model="confirmPassword" :placeholder="confirmPassword" />
          <span class="eye-icon" @click="toggleConfirmPasswordVisibility">
            <i :class="showConfirmPassword ? 'fas fa-eye' : 'fas fa-eye-slash'"></i>
          </span>
        </div>
      </div>
      <div class="button_container">
        <button @click="savePassword">保存</button>
        <button @click="cancel">取消</button>
      </div>
      <div class="forgot-password">
        <router-link to="/forgetpassword">忘记密码？</router-link>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        oldPassword: '',
        newPassword: '',
        confirmPassword: '',
        showNewPassword: false,
        showConfirmPassword: false,
      };
    },
    methods: {
      savePassword() {
        if (this.newPassword !== this.confirmPassword) {
          alert('新密码与确认密码不一致');
          return;
        }
        // 调用后端API更新密码
        alert('密码已更新，请重新登录');
        // 跳转到登录界面
        this.$router.push('/login');
      },
      cancel() {
        // 取消逻辑
        this.oldPassword = '';
        this.newPassword = '';
        this.confirmPassword = '';
      },
      toggleNewPasswordVisibility() {
        this.showNewPassword = !this.showNewPassword;
      },
      toggleConfirmPasswordVisibility() {
        this.showConfirmPassword = !this.showConfirmPassword;
      },
    },
  };
  </script>
  
  <style scoped>
  .change-password {
    padding: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
  }
  
  h3 {
    margin-top: 0;
  }
  
  
  .text_group {
    display: grid;
    grid-template-columns: 80px 1fr;
    align-items: center;
    gap: 10px;
  }
  
  .text_group label {
    text-align: right;
    margin-right: 10px;
  }
  
  .text_group input {
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  
  .password-input {
    display: flex;
    align-items: center;
    position: relative;
  }
  
  .password-input input {
    flex: 1;
    padding-right: 10px; /* 为眼睛图标留出空间 */
  }
  
  .eye-icon {
    position: absolute;
    right: 10px;
    cursor: pointer;
    color: #666;
  }
  
  .button_container {
    display: flex;
    justify-content: space-between;
    gap: 20px;
    width: 100%;
    max-width: 200px;
  }
  
  .button_container button {
    padding: 10px 20px;
    background-color: #42b983;
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .button_container button:hover {
    background-color: #369f6d;
  }
  
  .forgot-password {
    margin-top: 10px;
    text-align: center;
    width: 100%;
    max-width: 300px;
  }
  
  .forgot-password a {
    color: #42b983;
    text-decoration: none;
  }
  
  .forgot-password a:hover {
    text-decoration: underline;
  }
  </style>