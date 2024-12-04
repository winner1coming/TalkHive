<template>
  <div class="home">

    <h1>欢迎来到聊天室</h1>
    <router-link to = "/loginth">登录</router-link>
    <router-link to = "/security">设置</router-link>
    <router-link to = "/profile">主页</router-link>
    <router-link to="/chat">聊天</router-link>
    <router-link to="/contact">通讯录</router-link>
    <router-link to="/groups">群聊列表</router-link>
    <router-link to="/notes">笔记</router-link>
    <router-link to="/favorites">收藏</router-link>
    <router-link to="/ddls">DDL</router-link>
  </div>
  <div>
      <component :is="currentTabComponent" />
  </div>
</template>

<script>
export default {
  name: 'Home',
};
</script>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}
</style>