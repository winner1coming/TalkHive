<template>
  <div class="personal-homepage">
    <h2>个人主页</h2>
    <p>用户名: {{ user.username }}</p>
    <p>ID: {{ user.id }}</p>
    <router-link to="/edit-profile">编辑资料</router-link>
  </div>
</template>

<script>
export default {
  name: 'PersonalHomepage',
  computed: {
    user() {
      return this.$store.state.user;
    },
  },
};
</script>

<style scoped>
.personal-homepage {
  padding: 20px;
}
</style>