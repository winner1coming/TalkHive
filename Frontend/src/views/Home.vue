<template>
  <div class="home">
    <!-- 左侧导航栏 -->
    <aside class="sidebar">
      <div class="user-info">
        <img :src="avatar" alt="Avatar" class="avatar" />
        <span class="nickname">用户昵称</span>
      </div>
      <ul class="nav-links">
        <li><router-link to="/chat">聊天</router-link></li>
        <li><router-link to="/contact">通讯录</router-link></li>
        <li><router-link to="/setlist">设置</router-link></li>
        <li><router-link to="/workspace">工作区</router-link></li>
        <li><Link/></li>
      </ul>
      <ul class="logout">
        <li>
          <li><a href="#" @click="logout">Logout</a></li>
        </li>
      </ul>
    </aside>

    <!-- 右侧视图区域 -->
    <main class="content">
      <router-view></router-view>
    </main>
  </div>
</template>

<script>
import Link from './Link.vue';
export default {
  name: 'Home',
  data() {
    return {
      showDropdown: false,
      avatar: this.$store.state.user.avatar,
    };
  },
  components:{
    Link,
  },
  methods: {
    toggleDropdown() {
      this.showDropdown = !this.showDropdown;
    },
    logout() {
      localStorage.removeItem('isLoggedIn');
      this.$router.push('/login');
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
  background-color: #6dc79fb1;
  color: white;
  padding: 10px;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  position: relative;
}

.user-info {
  display: flex;
  align-items: center;
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
  color: white;
  text-decoration: none;
  display: block;
  padding: 10px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.nav-links a:hover {
  background-color: rgba(255, 255, 255, 0.1);
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
  background-color: #f5f5f5;
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
  color: white;
  text-decoration: none;
  display: block;
  padding: 10px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.logout a:hover {
  background-color: rgba(255, 255, 255, 0.1);
}
</style>