<template>
  <div class="accountlogin">
    <div class = "all">
      <img  class="avatar" :src = "'user.avatar' || '@/assets/images/avatar.jpg'"/>
      
      <div class="input-group">
        <label for="account">账号:</label>
        <input id="account" type="text" v-model="account" placeholder="请输入账号" />
      </div>
      
      <div class="input-group">
        <label for="password">密码:</label>
        <input id="password" type="password" v-model="password" placeholder="请输入密码" />
      </div>

      <!-- 记住密码选项 -->
      <div class="remember-me">
      <input type="checkbox" id="rememberMe" v-model="rememberMe" />
      <label for="rememberMe">记住密码</label>
      </div>

      <button class="login-button" @click="testlogin">登录</button>
      
      <div class = "link">
        <!-- 注册链接 -->
        <router-link to = "/register" class = "register">注册</router-link>
        <router-link to = "/forgetpassword" class ="reset">忘记密码</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { login } from '@/services/loginth.js'; // 导入登录 API
import CryptoJS from 'crypto-js';
import { mapActions, mapGetters} from 'vuex';

export default {
  //数据的存储
  data() {
    return {
      account: '666666',
      password: '666666',
      rememberMe:false,
      encryptionKey: 'TalkHiveProject',
    };
  },

  computed: {
    ...mapGetters(['user']), // 从 Vuex 获取用户信息
  },

  
  methods: {
    ...mapActions(['login']), // 映射 Vuex 的 login 方法

    //对密码进行加密
    encryptPassword(password) {
      return CryptoJS.AES.encrypt(password, this.encryptionKey).toString();
    },

    // 解密密码
    decryptPassword(encryptedPassword) {
      const bytes = CryptoJS.AES.decrypt(encryptedPassword, this.encryptionKey);
      return bytes.toString(CryptoJS.enc.Utf8);
    },

    // //检查账号密码是否为空（与后端连接需要把测试登录testlogin函数直接删除即可）

    async testlogin(){
    //         // 检查账号和密码是否为空
    //   // 调用登录方法
       await this.login();
      this.$router.push('/home');
     },
    async login() {

      if (!this.account) {
        alert('账号不能为空');
        return;
        }
        else if(!this.password){
          alert('密码不能为空');
        }

      try {
        const response = await login(
          {
            account: this.account,
            password: this.password,
          });

        //如果记住密码
        if (this.rememberMe) {
          const encryptedPassword = this.encryptPassword(this.password);
          localStorage.setItem('rememberedAccount', this.account);
          localStorage.setItem('rememberedPassword', encryptedPassword);
        } else {
          // 如果未勾选“记住密码”，则清除之前保存的账号和密码
          localStorage.removeItem('rememberedAccount');
          localStorage.removeItem('rememberedPassword');
        }

        if (response.success) {
          //更新全局变量
          this.$store.commit('SET_USER', {
            username: response.nickname,
            id: response.account_id,
            avatar: response.avatar,
          });

          this.$router.push('/home');
        } else {
          alert(response.message || '登录失败');
        }
      } catch (error) {
        alert(error || '登录失败');
      }
    },
  },

  mounted() {
    // 页面加载时，检查是否有记住的账号和密码
    const rememberedAccount = localStorage.getItem('rememberedAccount');
    const encryptedPassword = localStorage.getItem('rememberedPassword');

    if (rememberedAccount && encryptedPassword) {
      this.phone = rememberedAccount;
      this.password = this.decryptPassword(encryptedPassword); // 解密密码
      this.rememberMe = true; // 自动勾选“记住密码”
    }
  },

};
</script>

<style scoped>
.accountlogin{
    /* 使用 Flexbox 布局，使内容居中 */
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 60vh; /* 设置高度为视口高度 */
    padding: 10px;
    box-sizing: border-box;
}

.all {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 350px;
  padding: 20px;
  box-sizing: border-box;
}


.avatar{
  width: 100px;
  height: 100px;
  margin-top: 10px;
  margin-bottom: 20px;
  margin-left: 45px;
  border-radius: 100%;
}

.input-group {
  display: flex;
  margin-left: 20px;
  margin-bottom: 20px;
  align-items: center;
  width: 100%;
  max-width: 300px;
  box-sizing: border-box;
}

.input-group label {
  margin-right: 15px;
  font-size: 14px;
  align-items: center;
  color: #666;
  white-space: nowrap; /* 防止标签换行 */
}

.input-group input {
  flex: 1; /* 使输入框占据剩余空间 */
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 4px;
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