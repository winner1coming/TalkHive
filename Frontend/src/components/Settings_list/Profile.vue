<template>
  <!-- 个人信息页面容器 -->
  <div class="profile">
    <!-- 页面标题 -->
    <h2>个人信息</h2>
    
    <!-- 头像 -->
    <div class="avatar">
      <img :src="user.avatar" alt="" />
    </div>
    
    <!-- 昵称 -->
    <p class="nickname">{{ user.nickname }}</p>
    
    <!-- 性别标识 -->
    <p class="gender">
      <span v-if="user.gender === 'male'" class="gender-icon male"></span>
      <span v-else-if="user.gender === 'female'" class="gender-icon female"></span>
    </p>
    
    <!-- 账号 -->
    <p>账号: {{ user.account }}</p>
    
    <!-- 电话号码 -->
    <p>电话: {{ user.phone }}</p>
    
    <!-- 生日 -->
    <p>生日: {{ user.birthday }}</p>
    
    <!-- 个性签名 -->
    <p class="signature">个性签名部分{{ user.signature }}</p>
  </div>
</template>

<script>
export default {
  // 计算属性，用于从 Vuex store 中获取当前用户信息
  computed: {
    user() {
      // 返回 Vuex store 中的 user 状态
      return this.$store.state.user;
    },
  },
};
</script>

<style scoped>
/* 个人信息页面的样式 */
.profile {
  padding: 20px; /* 设置内边距 */
  text-align: center; /* 居中对齐 */
}

h2 {
  margin-bottom: 20px; /* 标题下边距 */
}

.avatar {
  margin-bottom: 20px; /* 头像下边距 */
}

.avatar img {
  width: 100px; /* 头像宽度 */
  height: 100px; /* 头像高度 */
  border-radius: 50%; /* 圆形头像 */
}

.nickname {
  font-size: 24px; /* 昵称字体大小 */
  margin-bottom: 10px; /* 昵称下边距 */
}

.gender {
  margin-bottom: 20px; /* 性别标识下边距 */
}

.gender-icon {
  display: inline-block;
  width: 20px; /* 性别图标宽度 */
  height: 20px; /* 性别图标高度 */
  border-radius: 50%; /* 圆形性别图标 */
}

.gender-icon.male {
  background-color: blue; /* 男性图标颜色 */
}

.gender-icon.female {
  background-color: pink; /* 女性图标颜色 */
}

p {
  margin-bottom: 10px; /* 段落下边距 */
}

.signature {
  margin-top: 20px; /* 个性签名上边距 */
  font-style: italic; /* 斜体 */
}
</style>