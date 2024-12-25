<template>
    <div class="smslogin">
      <img  class="avatar" src = '@/assets/images/avatar.jpg'/>
      
      <div class="input-group">
        <label for="email">邮箱</label>
        <input id="email" type="text" v-model="email" placeholder="请输入邮箱" />
        <p v-if="!isEmailValid && email" class="error-message">请输入有效的邮箱地址</p>
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
  import { smsLogin, sendSmsCode } from '@/services/loginth.js'; // 导入登录和发送验证码 API
  import { mapActions, mapGetters} from 'vuex';
  
  export default {
    computed: {
      ...mapGetters(['user']), // 从 Vuex 获取用户信息
      isEmailValid() {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(this.email);
      },
    },

    data() {
      return {
        email: '',
        smsCode: '',
        Code: '',
      };
    },
    
    methods: {
      ...mapActions(['login']), // 映射 Vuex 的 login 方法


      // 登录方法
      async smsLogin() {
        if(!this.email){
          alert("邮箱不能为空！");
        }
        else if(!this.smsCode){
          alert("验证码不能为空！");
        }

        this.validateCode();  
          
        try {
          const response = await smsLogin({
            email : this.email,
          });

          if(response.success){
            this.$store.commit('SET_USER', {
            username: response.nickname,
            id: response.account_id,
            avatar: response.avatar,
            });

            this.$router.push('/home');
          }
          else {
            alert(response.message);
          }
        } catch (error) {
          alert(error || '登录失败');
        }
      },
      
      async validateCode(){
        if(this.Code){
          if(this.Code !== this.smsCode){
              alert('验证码错误');
              return;
          }
        }
        else{
          alert('请先获取验证码！');
          return;
        }
      },
      // 发送验证码方法
      async sendSmsCode() {
        
        try {
          const response = await sendSmsCode({
            command:'smsLogin',
            email:this.email,
          });

          if (response.success) {
            alert('验证码已发送');
            this.Code = response.code;
          } else {
            alert(response.message);
          }
        } catch (error) {
          alert(error || '发送验证码失败');
        }
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

.error-message {
  color: red;
  font-size: 12px;
  margin-top: 5px;
  flex-direction: column;
}
  </style>