<template>
  <div class="home">
    <h1>欢迎来到聊天室</h1>
    <router-link to="/messages">消息列表</router-link>
    <router-link to="/friends">好友列表</router-link>
    <router-link to="/groups">群聊列表</router-link>
    <router-link to="/notes">笔记</router-link>
    <router-link to="/favorites">收藏</router-link>
    <router-link to="/ddls">DDL</router-link>
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