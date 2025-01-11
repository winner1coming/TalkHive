<template>
  <div class="home">
    <!-- 左侧导航栏 -->
    <aside class="sidebar">
      <div class="user-info">
        <img :src="avatar" alt="Avatar" class="avatar" />
        <span class="nickname">{{nickname}}</span>
      </div>
      <ul class="nav-links">
        <li><router-link to="/chat">
          <img src="@/assets/icon/chat-icon.png" alt="Chat" class="icon"/>
        </router-link></li>
        <li><router-link to="/contact">
          通讯录
        </router-link></li>
        <li><router-link to="/setlist">设置</router-link></li>
        <li><router-link to="/workspace">工作区</router-link></li>
        <li><Link/></li>
      </ul>
      <ul class="logout">
        <li>
          <li><a href="#" @click="showLogoutConfirmation">Logout</a></li>
        </li>
      </ul>
      <div v-if="showConfirmation" class="confirmation-modal">
        <div class="modal-content">
          <span class="close" @click="hideLogoutConfirmation">&times;</span>
          <p>是否确认退出登录？</p>
          <div class="modal-buttons">
            <button @click="logout">确认</button>
            <button @click="hideLogoutConfirmation">取消</button>
          </div>
        </div>
      </div>
    </aside>

    <!-- 右侧视图区域 -->
    <main class="content">
      <router-view></router-view>
    </main>
  </div>
</template>

<script>
import Link from './Link.vue';
import { logout } from '@/services/settingView';
export default {
  name: 'Home',
  data() {
    return {
      showDropdown: false,
      showConfirmation:false,
    };
  },
  computed: {
  avatar() {
    return this.$store.state.user.avatar;
  },
  nickname() {
    return this.$store.state.user.username;
  },
  },
  components:{
    Link,
  },
  methods: {
    toggleDropdown() {
      this.showDropdown = !this.showDropdown;
    },
    showLogoutConfirmation() {
        this.showConfirmation = true;
      },
    hideLogoutConfirmation() {
      this.showConfirmation = false;
    },
    async logout() {
      try{
          const response = await logout();
          if(response.success){
            alert('已退出登录~');
            localStorage.removeItem('isLoggedIn');
            this.$router.push('/loginth');
          // 你可以在这里添加退出登录的逻辑
          }else{
            alert("退出登录失败，请重试！");
          }
        }catch(error){
          console.error("退出登录失败！");
      }
    },
  },
};
</script>

<style scoped>
.home {
  display: flex;
  height: 100vh;
}

/* 左侧导航栏样式 */
.sidebar {
  width: 120px;
  background-color: var(--sidebar-background-color);
  color: var(--sidebar-text-color);
  padding: 10px;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  position: relative;
  font-size:var(--font-size-medium);
}

.user-info {
  display: flex;
  align-items: center;
  flex-direction: column;
  margin-bottom: 20px;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
}

.nickname {
  font-weight: bold;
  color: var(--sidebar-text-color);
  font-size: var(--font-size-medium);
}

.nav-links {
  list-style: none;
  padding: 0;
  width: 100%;
}

.nav-links li {
  margin: 10px 0;
  position: relative;
}

.nav-links a {
  color: var(--sidebar-background-color);
  text-decoration: none;
  display: block;
  padding: 5px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.nav-links a:hover {
  background-color: var(--sidebar-background-color);
}

.dropdown {
  background-color: #fff;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-top: 5px;
  padding: 10px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  width: 100%;
}

.dropdown li {
  margin: 0;
}

.dropdown a {
  color: #333;
  display: block;
  padding: 5px 0;
}

/* 右侧视图区域样式 */
.content {
  flex: 1;
  overflow-y: auto;
  background-color: var(--background-color);
}

/* 退出登录按钮样式 */
.logout {
  list-style: none;
  padding: 0;
  width: 100%;
  margin-top: auto; /* 将退出按钮推到最底部 */
}

.logout li {
  margin: 10px 0;
}

.logout a {
  color: var(--sidebar-text-color);
  text-decoration: none;
  display: block;
  padding: 10px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.nav-links a .icon {
  width: 50px; /* 图标大小 */
  height: 50px; /* 图标大小 */
  margin-right: 8px; /* 图标和文字之间的间距 */
}

.logout a:hover {
  background-color: var(--sidebar-text-color);
}

.confirmation-modal {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%); /* 确保模态框居中 */
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
  }
  
  .modal-content {
    background-color: white;
    padding: 20px;
    border-radius: 8px;
    position: relative;
    width: 300px;
  }
  
  .close {
    position: absolute;
    top: 10px;
    right: 10px;
    cursor: pointer;
    font-size: 20px;
  }
  
  .modal-buttons {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
  
  .modal-buttons button {
    margin-left: 10px;
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .modal-buttons button:first-child {
    background-color: #42b983;
    color: white;
  }
  
  .modal-buttons button:last-child {
    background-color: #ccc;
    color: black;
  }
</style>