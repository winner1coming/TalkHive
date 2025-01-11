<template>
  <div class="home">
    <!-- 左侧导航栏 -->
    <aside class="sidebar">
      <div class="user-info">
        <img :src="avatar" alt="Avatar" class="avatar" />
      </div>
      <ul class="nav-links">
        <li><router-link to="/chat" title="聊天">
          <img src="@/assets/icon/chat-icon.png" alt="Chat" class="icon"/>
        </router-link></li>
        <li><router-link to="/contact" title="通讯录">
          <img src="@/assets/icon/contact.png" alt="Contact" class="icon"/>
        </router-link></li>
        <li><router-link to="/setlist" title="设置">
          <img src="@/assets/icon/setting-icon.png" alt="Settings" class="icon"/>
        </router-link></li>
        <li><router-link to="/workspace" title="工作空间">
          <img src="@/assets/icon/workspace-icon.png" alt="Workspace" class="icon"/>
        </router-link></li>
        <li><Link/></li>
      </ul>
      <ul class="logout">
        <li>
          <li><a href="#" @click="showLogoutConfirmation">
            <img src="@/assets/icon/logout.png" alt="Logout" class="icon" title="退出登录"/>
          </a></li>
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
  width: 60px;
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
  gap:10px;
  margin-bottom: 10px;
}

.avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
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

.nav-links a .icon {
  width: 40px; /* 图标大小 */
  height: 40px; /* 图标大小 */
  margin-right: 5px; /* 图标和文字之间的间距 */
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
  margin-top: auto; /* 将退出按钮推到最底部 */
}

.logout li {
  margin: 10px 0;
}

.logout a .icon{
  width: 25px;
  height: 25px;
  margin-right: 5px;
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
    z-index: 1000;
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