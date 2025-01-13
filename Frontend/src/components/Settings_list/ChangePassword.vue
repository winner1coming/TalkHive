<template>
    <div class="change-password">
      <h3>更改密码</h3>

      <!--原密码的输入框-->
      <div class="text_group">
        <label for="oldPassword">原密码:</label>
        <input type="password" v-model="oldPassword" @blur="validatePassword(oldPassword,'oldp')" :placeholder="oldPassword" />
      </div>
      <p v-if="errors.oldp" class="error">{{ errors.oldp }}</p>

      <!--新密码输入框-->
      <div class="text_group">
        <label for="newPassword">新密码:</label>
        <div class="password-input">
          <input type="password" v-model="newPassword" @blur="validatePassword(newPassword,'newp')" :placeholder="newPassword" />
        </div>
      </div>
      <p v-if="errors.newp" class="error">{{ errors.newp }}</p>

      <!--验证密码的输入框-->
      <div class="text_group">
        <label for="confirmPassword">确认密码:</label>
        <div class="password-input">
          <input :type="showConfirmPassword ? 'text' : 'password'" v-model="confirmPassword" @input="validateChange" :placeholder="confirmPassword" />
          <span class="eye-icon" @click="toggleConfirmPasswordVisibility">
            <i :class="showConfirmPassword ? 'fas fa-eye' : 'fas fa-eye-slash'"></i>
          </span>
        </div>
      </div>
      <p v-if="errors.validatep" class="error">{{ errors.validatep }}</p>

      <!--按钮设置-->
      <div class="button_container">
        <button @click="savepassword">保存</button>
        <button @click="cancel">取消</button>
      </div>

      <!--提供忘记密码的链接-->
      <div class="forgot-password">
        <router-link to="/forgetpassword">忘记密码？</router-link>
      </div>
    </div>
  </template>
  
  <script>
  import { savePassword } from '@/services/settingView.js';

  export default {
    props:{
      user: {
        type: Object,
        required: true,
      },
    },

    data() {
      return {
        oldPassword: '',
        newPassword: '',
        confirmPassword: '',
        showConfirmPassword: false,
        errors:{
            newp:'',
            oldp:'',
            validatep:'',
        },
      };
    },

    methods: {
      async savepassword() {
        if(!this.validate_comfirmPassword()){
          return;
        }

        try{
          const response = await savePassword({
            newpassword:this.newPassword,
          });
          if(response.success){
            alert('密码已更新，请重新登录')
            this.$router.push('/loginth');
          }else{
            alert(response.message || '密码更改失败');
          }

        }catch(error){
          console.error("更改密码失败:",error);
        }

      },

      cancel() {
        // 取消逻辑
        this.oldPassword = '';
        this.newPassword = '';
        this.confirmPassword = '';
        this.errors.newp ='';
        this.errors.oldp ='';
        this.errors.validatep = '';
      },


      validate_comfirmPassword() {
        const password = this.user.password;
        if (this.oldPassword !== password) {
          alert("原密码输入错误");
          return false;
        } else if(this.oldPassword === this.newPassword){
          alert("新密码与原密码相同！");
          return false;
        }
        return true;
      },

      validatePassword(password , key){
          if(!password){
            this.errors[key] = '请输入密码！';
          }else if(password.length < 6){
            this.errors[key] = '密码长度不能小于6位！';
          }else if(password.length > 20)
          {
            this.errors[key] = '密码长度不能超过20位';
          }else if(/\s/.test(password)){
            this.errors[key] = '密码不能包含空格！';
          }else if(!/^[a-zA-Z0-9]+$/.test(password)){
            this.errors[key] = '密码只能包含字母和数字！';
          }else{
            this.errors[key] = '';
          }
      },

      validateChange(){
        if(!this.confirmPassword){
          this.errors.validatep = '请输入确认密码！';
        }else if(this.newPassword !== this.confirmPassword){
          this.errors.validatep = '两次输入的密码不一致！';
        }
        else{
          this.errors.validatep = '';
        }
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
    margin-top: 60px;
    font-size: var(--font-size);
    color: var(--text-color);
  }
  
  h3 {
    margin-top: 0;
    font-size: var(--font-size-large);
  }
  
  
  .text_group {
    display: grid;
    grid-template-columns: 80px 1fr;
    align-items: center;
    gap: 10px;
  }
  
  .text_group label {
    text-align: right;
  }
  
  .text_group input {
    padding: 6px;
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
    font-size: var(--font-size);
    color: var(--button-text-color);
  }
  
  .button_container button {
    padding: 8px 10px;
    background-color:var(--button-background-color);
    color: var(--button-text-color);
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .button_container button:hover {
    background-color: var(--button-background-color1);
  }

  .button_container button.active, .button_container button.active:hover{
    background-color: var(--button-background-color2);
  }
  
  .forgot-password {
    margin-top: 10px;
    text-align: center;
    width: 100%;
    max-width: 300px;
  }
  
  .forgot-password a {
    color:#62ca69;
    text-decoration: var(--button-background-color2);
  }
  
  .forgot-password a:hover {
    text-decoration: underline;
  }

  p{
    color: crimson;
    width: 200px;
    font-size: var(--font-size-small);
  }
  </style>