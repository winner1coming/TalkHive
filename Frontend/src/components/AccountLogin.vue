<template>
  <div class="accountlogin">
    <div class = "all">
      <img  class="avatar" :src = "avatar"/>
      
      <div class="input-group">
        <label for="account">账号:</label>
        <input 
          id="account" 
          type="text" 
          v-model="account" 
          @blur="hideDropdown" 
          @input="handleInput" 
          placeholder="请输入账号" />
        <ul v-if="showDropdown" class="dropdown">
        <li
          v-for="matchedAccount in matchedAccounts"
          :key="matchedAccount"
          @mousedown="selectAccount(matchedAccount)"
        >
            {{ matchedAccount }}
          </li>
        </ul>
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

      <button class="login-button" @click="login">登录</button>
      
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
import img from '@/assets/images/avatar.jpg';

export default {
  //数据的存储
  data() {
    return {
      account: '',
      password: '',
      rememberMe:false,
      encryptionKey: 'TalkHiveProject',
      avatar:img,
      users:[],
      matchedAccounts: [],
      showDropdown:false,
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

        if (response.success) {
          this.avatar = `data:${response.mimeType};base64,${response.avatar}`;
          //更新全局变量
          this.$store.commit('SET_USER', {
            username: response.nickname,
            id: response.account_id,
            avatar: this.avatar,
          });
          this.$store.commit('SET_LINKS',response.links);

          let users = JSON.parse(localStorage.getItem('users')) || [];
          //本地缓存的处理
          const userInfo = {
            account : this.account,
            avatar: this.avatar,
            email: response.email || '',
            password: this.rememberMe ? this.encryptPassword(this.password) : '',
          };
          //是否已经存在账号
          const index = users.findIndex(user => user.account === this.account);
          if(index !== -1){
            users[index] = userInfo;
          }else{
            users.push(userInfo);
          }
          localStorage.setItem('users', JSON.stringify(users));

          alert(response.message);
          this.$router.push('/home');
        } else {
          alert(response.message);
        }
      } catch (error) {
        alert(error || '登录失败');
      }
    },

    handleInput() {
      if (this.account) {
        // 模糊匹配账号（最左匹配）
        this.matchedAccounts = this.users
          .map(user => user.account)
          .filter(account => account.startsWith(this.account));
        this.showDropdown = this.matchedAccounts.length > 0; // 如果有匹配的账号，显示下拉框

        const matchedUser = this.users.find(user => user.account === this.account);
        if (matchedUser) {
          this.avatar = matchedUser.avatar; // 如果账号存在，设置头像
          if (matchedUser.password) {
            this.password = this.decryptPassword(matchedUser.password);
          }else{
            this.password = '';
          }
        }
        else{
          this.avatar = img;
          this.password='';
        }
      } else {
        this.matchedAccounts = [];
        this.showDropdown = false;
        this.avatar=img;
        this.password = '';
      }
    },

    // 选择账号
    selectAccount(account) {
      this.account = account; // 填充账号到输入框
      this.showDropdown = false; // 隐藏下拉框

      // 直接从缓存中获取对应的头像
      const matchedUser = this.users.find(user => user.account === account);
      if(matchedUser){
        this.avatar = matchedUser.avatar; // 设置头像  
        // 如果该账号存储了加密的密码，则解密并填充到密码输入框
        if (matchedUser.password) {
          this.password = this.decryptPassword(matchedUser.password);
        }else{
          this.password = '';
        }      
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
  position: relative;
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