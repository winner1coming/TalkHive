<template>
  <div class="home">
    <!-- 左侧导航栏 -->
    <aside class="sidebar">
      <div class="user-info">
        <img :src="avatar" alt="Avatar" class="avatar" @click="toggleProfile"/>
        <PersonProfileCard ref="profileCard" />
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
import { getPersonProfileCard } from '@/services/api';
import Link from './Link.vue';
import PersonProfileCard from '@/components/base/PersonProfileCard.vue';
import {logout, getSystemSetting} from '@/services/settingView.js';


export default {
  name: 'Home',
  data() {
    return {
      showDropdown: false,
      showConfirmation:false,
      showProfile:false,
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
    PersonProfileCard,
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
            this.$router.push('/loginth');
          // 你可以在这里添加退出登录的逻辑
          }else{
            alert("退出登录失败，请重试！");
          }
        }catch(error){
          console.error("退出登录失败！");
      }
    },
    async toggleProfile(event) {
      try{
        const response = await getPersonProfileCard(this.$store.state.user.id); 
        if(response.status !== 200){
          this.$root.notify(response.data.message, 'error');
          return;
        }
        this.showProfile= true;
        const boundD = '50px';
        const boundR = '50px';
        const profile = response.data.data;
        this.$refs.profileCard.show(event, profile, boundD,boundR);
      }catch(err){
        console.log(err);
      }
    },
    hideProfileCard() {
      this.showProfile = false;
      document.removeEventListener('click', this.handleClickOutside);
    },
    handleClickOutside(event) {
      if (this.$refs.profileCard && !this.$refs.profileCard.$el.contains(event.target)) {
        this.hideProfileCard();
      }
    },

    async fetchSystemSettings() {
      try {
        const response = await getSystemSetting();
        if (response.success) {
          let BackGround = '';
          if (response.data.background !== '') {
            BackGround = `data:${response.data.mimeType};base64,${response.data.background}`;
          }
          this.$store.commit('SET_SETTINGS', {
            theme: response.data.theme,
            fontSize: response.data.fontSize,
            fontStyle: response.data.fontStyle,
            sound: response.data.sound,
            isNotice: response.data.notice,
            isNoticeGroup: response.data.noticeGroup,
            background: BackGround,
          });
        } else {
          alert(response.message);
        }
      } catch (error) {
        alert(error, '获取系统设置失败，请检查网络');
      }
    },
  },
  mounted() {
    this.fetchSystemSettings();
    document.addEventListener('click', this.handleClickOutside);
    this.$store.dispatch('connectWebSocket'); // 连接 WebSocket
  },
  beforeUnmount() {
    document.removeEventListener('click', this.handleClickOutside);
  },
};
</script>

<style scoped>
.home {
  display: flex;
  height: 100%;
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
  font-size:var(--font-size);
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

.nav-links a.router-link-active {
  background-color: var(--sidebar-background-color1);
}

.dropdown {
  background-color: var(--background-color);
  border: 1px solid var(--background-color2);
  border-radius: 4px;
  margin-top: 5px;
  padding: 10px;
  box-shadow: 0 2px 4px var(--background-color2);
  width: 100%;
}

.dropdown li {
  margin: 0;
}

.dropdown a {
  color: var(--text-color);
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
    font-size: var(--font-size-large);
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
    background-color: var(--background-color);
    color: var(--sidebar-text-color);
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
    font-size: var(--font-size-large);
  }
  
  .modal-buttons {
    display: flex;
    justify-content: space-evenly;
    margin-top: 30px;
    font-size: var(--font-size);
  }
  
  .modal-buttons button {
    margin-left: 10px;
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    color: var(--button-text-color);
  }
  
  
  .modal-buttons button:first-child {
    background-color: var(--button-background-color2);
    color: var(--button-text-color);
  }
  
  .modal-buttons button:last-child {
    background-color:var(--background-color2);
    color: var(--button-text-color);
  }

  .modal-buttons button:hover{
    background-color: var(--button-background-color);
  }

  .modal-buttons button:active{
    background-color: var(--button-background-color2);
  }

  .modal-content p{
    color: var(--text-color);
    font-size: var(--font-size);
  }
</style>