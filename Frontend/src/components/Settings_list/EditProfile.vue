<template>
  <!-- 编辑资料页面容器 -->
  <div class="editprofile">
    <!-- 页面标题 -->
    <h2>编辑资料</h2>
      <div class ="avatar">
        <span >头像</span>
        <img :src="avatar" alt="avatar" class="avatar" @click="openFilePicker" />
        <input type="file" ref="fileInput" style="display: none;" @change="handleFileChange" />
      </div>
      <div class ="input_text">
        <label for = "username">用户名</label>
        <!-- 用户名输入框，使用 v-model 双向绑定到 data 中的 username 属性 -->
        <input id = "username" type="text" v-model="username" placeholder="用户名" />
      </div>
      <div class ="input_text">
        <!-- ID 输入框，使用 v-model 双向绑定到 data 中的 id 属性 -->
        <input type="text" v-model="id" placeholder="ID" />
      </div>

    
    <!-- 保存按钮，点击时触发 saveProfile 方法 -->
    <button @click="saveProfile">保存</button>
    
    <!-- 取消按钮，点击时触发 cancelEdit 方法 -->
    <button @click="cancelEdit">取消</button>
  </div>
</template>

<script>
import { getProfile, updateProfile } from '@/services/api';

export default {
  // 组件的 data 函数，返回一个对象，包含组件的响应式数据
  data() {
    return {
      avatar:'',
      // 用户名输入框的值
      username: '',
      // ID 输入框的值
      id: '',
      // 原始用户名，用于取消操作时恢复
      originalUsername: '',
      // 原始 ID，用于取消操作时恢复
      originalId: '',
    };
  },
  
  // 组件的生命周期钩子，在组件创建时调用
  created() {
    // 从数据库获取用户信息
    this.fetchProfile();
  },
  
  // 组件的方法定义
  methods: {
    // 从数据库获取用户信息
    async fetchProfile() {
      try {
        const profile = await getProfile();
        this.username = profile.username;
        this.id = profile.id;
        // 保存原始用户名和 ID
        this.originalUsername = profile.username;
        this.originalId = profile.id;
      } catch (error) {
        console.error('Failed to fetch profile:', error);
      }
    },
    
    // 保存资料方法，处理用户点击保存按钮时的逻辑
    async saveProfile() {
      try {
        await updateProfile({ username: this.username, id: this.id });
        // 更新原始用户名和 ID
        this.originalUsername = this.username;
        this.originalId = this.id;
      } catch (error) {
        console.error('Failed to save profile:', error);
      }
    },
    
    // 取消编辑方法，处理用户点击取消按钮时的逻辑
    cancelEdit() {
      // 恢复原始用户名和 ID
      this.username = this.originalUsername;
      this.id = this.originalId;
    },

    openFilePicker() {
      this.$refs.fileInput.click();
    },

    handleFileChange(event) {
      const file = event.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          this.avatar = e.target.result;
        };
        reader.readAsDataURL(file);
      }
    },
  },
};
</script>

<style scoped>
/* 编辑资料页面的样式 */
.editprofile {
  padding: 20px; /* 设置内边距 */
}
</style>