<template>
    <div class="smslogin">
      <img  class="avatar" src = '@/assets/images/avatar.jpg'/>
      
      <div class="input-group">
        <label for="phoneNumber">手机号</label>
        <input id="phoneNumber" type="text" v-model="phoneNumber" placeholder="请输入手机号" />
      </div>
      
      <div class="verificate">
        <div class="input-group">
          <label for="smsCode">验证码</label>
          <input id="smsCode" type="text" v-model="smsCode" placeholder="请输入验证码" />
        </div>
        <button class="send-sms-button" @click="sendSmsCode">获取</button>
      </div>

      <button class="login-button" @click="smsLogin">登录</button>
      <div class = "link">
      <!-- 注册链接 -->
      <router-link to = "/register" class=" register">注册</router-link>
      <router-link to = "/forgetpassword" class=" reset">忘记密码</router-link>
      </div>
    </div>
  </template>
  
  <script>
  import { smsLogin, sendSmsCode } from '@/services/api'; // 导入登录和发送验证码 API
  
  export default {
    data() {
      return {
        phoneNumber: '',
        smsCode: '',
        code: '',
      };
    },
    
    methods: {
      // 登录方法
      async smsLogin() {
        this.validateCode();
        try {
          const response = await smsLogin({
            phoneNumber : this.phoneNumber,
          });
          if(response.success){
            this.$router.push('/');
          }
          else {
            alert(response.message || '登录失败');
          }
        } catch (error) {
          alert(error || '登录失败');
        }
      },
      
      async validateCode(){
        if(Code){
          if(Code !== this.smsCode){
              alert('验证码错误');
              return;
          }
        }
        else{
          alter('请先获取验证码！');
          return;
        }
      },
      // 发送验证码方法
      async sendSmsCode() {
        if (!this.validatePhoneNumber(this.phoneNumber)) {
          alert('请输入有效的手机号码');
          return;
        }
        
        try {
          const response = await sendSmsCode({
            command:'smsLogin',
            phoneNumber:this.phoneNumber,
          });

          if (response.success) {
            alert('验证码已发送');
            this.Code = response.code;
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

  .smslogin{
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 60vh; /* 设置高度为视口高度 */
    padding: 10px;
    box-sizing: border-box;
  }

  .avatar{
  width: 100px;
  height: 100px;
  margin-top: 10px;
  margin-bottom: 20px;
  margin-left:35px;
  border-radius: 100%;
}

  .input-group {
    display: flex;
    align-items: center;
    justify-items: center;
    margin-bottom: 20px;
    margin-left: 60px;
    width: 100%;
    max-width: 300px;
  }

  .verificate{
    display: flex;
    align-items: center;
    justify-self: start;
    justify-content: center;
  }
  
  .input-group label {
    margin-right: 10px;
    font-size: 14px;
    color: #666;
    white-space: nowrap; /* 防止标签换行 */
  }
  
  .input-group input {
    flex: 0; /* 使输入框占据剩余空间 */
    padding: 10px;
    font-size: 16px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  
  .send-sms-button {
    padding: 10px;
    font-size: 14px;
    width: 100%;
    color: #fff;
    background-color: #42b983;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    margin-right: 10px;
    margin-bottom: 20px;
    
  }
  
  .login-button {
  width: 100%;
  max-width: 100px;
  padding: 10px;
  font-size: 16px;
  color: #fff;
  background-color: #42b983;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  margin-top: 10px;
  margin-left: 20px;
}

.login-button:hover {
  background-color: #369f6e;
}

  .link {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  width: 100%;
  margin-top: 20px;
}

.register {
  text-decoration: none;
  justify-items: start;
  color: #42b983;
  font-size: 14px;
}

.reset {
  text-decoration: none;
  color: #42b983;
  font-size: 14px;
}

.register:hover, .reset:hover {
  text-decoration: underline;
}

  </style>