<template>
    <div class="smslogin">
      <img  class="avatar" :src ="avatar" />
      
      <div class="input-group">
        <label for="email">邮箱</label>
        <input id="email" type="text" v-model="email" @blur="hideDropdown" @input="handleInput" placeholder="请输入邮箱" />
        <ul v-if="showDropdown" class="dropdown">
          <li
            v-for="matchedAccount in matchedAccounts"
            :key="matchedAccount"
            @mousedown="selectAccount(matchedAccount)"
          >
            {{ matchedAccount }}
          </li>
        </ul>
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
  import img from '@/assets/images/avatar.jpg';
  
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
        avatar:img,
        users:[],
        matchedAccounts: [],
        showDropdown:false,
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

        if(this.validateCode()){
          return;
        } 
          
        try {
          const response = await smsLogin({
            email : this.email,
          });

          if(response.success){
            this.avatar = `data:${response.mimeType};base64,${response.avatar}`;
            this.$store.commit('SET_USER', {
            username: response.nickname,
            id: response.account_id,
            avatar: this.avatar,
            });
            this.$store.commit('SET_LINKS',response.links);

            let users = JSON.parse(localStorage.getItem('users')) || [];
            //本地缓存的处理
            const userInfo = {
              account : '',
              avatar: this.avatar,
              email: this.email,
              password:'',
            };
            //是否已经存在账号
            const index = users.findIndex(user => user.email === this.email);
            if(index !== -1){
              users[index] = userInfo;
            }else{
              users.push(userInfo);
            }
            localStorage.setItem('users', JSON.stringify(users));

            alert(response.message);
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
              return false;
          }
          return true;
        }
        else{
          alert('请先获取验证码！');
          return false;
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

        handleInput() {
        if (this.email) {
          // 模糊匹配账号（最左匹配）
          this.matchedAccounts = this.users
            .map(user => user.email)
            .filter(email => email.startsWith(this.email));
          this.showDropdown = this.matchedAccounts.length > 0; // 如果有匹配的账号，显示下拉框
          const matchedUser = this.users.find(user => user.email === this.email);
          if (matchedUser) {
            this.avatar = matchedUser.avatar; // 如果账号存在，设置头像
          }
          else{
            this.avatar = img;
          }
        } else {
          this.matchedAccounts = [];
          this.showDropdown = false;
          this.avatar=img;
        }
      },

      // 选择账号
      selectAccount(email) {
        this.email = email; // 填充账号到输入框
        this.showDropdown = false; // 隐藏下拉框

        // 直接从缓存中获取对应的头像
        const matchedUser = this.users.find(user => user.email === email);
        if(matchedUser){
          this.avatar = matchedUser.avatar; // 设置头像       
        }
      },

      // 隐藏下拉框
      hideDropdown() {
        setTimeout(() => {
            this.showDropdown = false;
        }, 200); // 延迟隐藏，避免点击下拉框时立即隐藏
      },
    },

    mounted() {
    this.users = JSON.parse(localStorage.getItem('users')) || [];
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
    position: relative;
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

.dropdown {
  position: absolute;
  width: 85%;
  border: 1px solid #ccc;
  border-radius: 4px;
  background-color: #fff;
  list-style: none;
  padding: 0;
  margin: 0;
  max-height: 100px;
  overflow-y: auto;
  z-index: 1000;
  top:100%;
  left: 45px;
}

.dropdown li {
  padding: 8px 20px;
  cursor: pointer;
}

.dropdown li:hover {
  background-color: #f0f0f0;
}

.dropdown::-webkit-scrollbar {
  width: 8px; /* 滚动条宽度 */
}

.dropdown::-webkit-scrollbar-thumb {
  background-color: #ccc; /* 滚动条颜色 */
  border-radius: 4px; /* 滚动条圆角 */
}

.dropdown::-webkit-scrollbar-track {
  background-color: #f0f0f0; /* 滚动条轨道颜色 */
}
  </style>