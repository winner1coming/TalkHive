<template>
  <!-- 编辑资料页面容器 -->
  <div class="editprofile">
    <!-- 页面标题 -->
    <h2>编辑资料</h2>
    <div class="avatar-container">
      <img :src="avatar" alt="avatar" class="headavatar" @click="showAvatarPreview" />
      <input type="file" ref="fileInput" style="display: none;" @change="handleFileChange" />
      <button @click="openFilePicker">上传</button>
    </div>
    <div class="input_text">
      <label>账号:</label>
      <span>{{ id }}</span>
    </div>
    <div class="input_text">
      <label for="username">用户名:</label>
      <input id="username" type="text" v-model="username" :placeholder="username" />
    </div>
    <div class="input_text">
      <label>性别:</label>
      <div class="gender-options">
        <label>
          <input type="radio" v-model="gender" value="male" />
          男
        </label>
        <label>
          <input type="radio" v-model="gender" value="female" />
          女
        </label>
      </div>
    </div>
    <div class="input_text">
      <label>生日:</label>
      <input type="text" v-model="birthday" :placeholder="birthday" @click="showDatePicker" readonly />
      <div v-if="showDatePickerFlag" class="date-picker">
        <input type="date" v-model="birthday" @change="hideDatePicker" />
      </div>
    </div>
    <div class="input_sig">
      <label>个性签名:</label>
      <textarea v-model="signature" :placeholder="signature" maxlength="100"></textarea>
      <span class="signature-count">{{ signature.length }}/100</span>
    </div>
    <!-- 保存按钮，点击时触发 saveProfile 方法 -->
    <div class="button_container">
    <button class="save_button" @click="saveProfile">保存</button>
    <!-- 取消按钮，点击时触发 cancelEdit 方法 -->
    <button class="cancle_button" @click="cancelEdit">取消</button>
    </div>
    <!-- 头像预览弹窗 -->
    <div v-if="showPreview" class="avatar-preview">
      <img :src="avatar" alt="avatar" class="avatar-large" />
      <button @click="hideAvatarPreview">×</button>
    </div>
  </div>
</template>

<script>
import { getProfile, updateProfile } from '@/services/api';

export default {
  // 组件的 data 函数，返回一个对象，包含组件的响应式数据
  data() {
    return {
      avatar: '',
      // 用户名输入框的值
      username: '',
      // ID 输入框的值
      id: '',
      // 性别
      gender: '',
      // 生日
      birthday: '',
      // 个性签名
      signature: '',
      // 原始用户名，用于取消操作时恢复
      originalUsername: '',
      // 原始 ID，用于取消操作时恢复
      originalId: '',
      // 原始性别，用于取消操作时恢复
      originalGender: '',
      // 原始生日，用于取消操作时恢复
      originalBirthday: '',
      // 原始个性签名，用于取消操作时恢复
      originalSignature: '',
      // 是否显示头像预览
      showPreview: false,
      // 是否显示日期选择器
      showDatePickerFlag: false,
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
        this.gender = profile.gender;
        this.birthday = profile.birthday;
        this.signature = profile.signature;
        // 保存原始用户名和 ID
        this.originalUsername = profile.username;
        this.originalId = profile.id;
        this.originalGender = profile.gender;
        this.originalBirthday = profile.birthday;
        this.originalSignature = profile.signature;
      } catch (error) {
        console.error('Failed to fetch profile:', error);
      }
    },

    // 保存资料方法，处理用户点击保存按钮时的逻辑
    async saveProfile() {
      try {
        await updateProfile({
          username: this.username,
          id: this.id,
          gender: this.gender,
          birthday: this.birthday,
          signature: this.signature,
        });
        // 更新原始用户名和 ID
        this.originalUsername = this.username;
        this.originalId = this.id;
        this.originalGender = this.gender;
        this.originalBirthday = this.birthday;
        this.originalSignature = this.signature;
      } catch (error) {
        console.error('Failed to save profile:', error);
      }
    },

    // 取消编辑方法，处理用户点击取消按钮时的逻辑
    cancelEdit() {
      // 恢复原始用户名和 ID
      this.username = this.originalUsername;
      this.id = this.originalId;
      this.gender = this.originalGender;
      this.birthday = this.originalBirthday;
      this.signature = this.originalSignature;
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

    showAvatarPreview() {
      this.showPreview = true;
    },

    hideAvatarPreview() {
      this.showPreview = false;
    },

    showDatePicker() {
      this.showDatePickerFlag = true;
    },

    hideDatePicker() {
      this.showDatePickerFlag = false;
    },
  },
};
</script>

<style scoped>
/* 编辑资料页面的样式 */
.editprofile {
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
  padding: 20px; /* 设置内边距 */
  align-items: center;
  justify-items: center;
}

h2{
  align-items: center;
  margin-top: 30px;
}

.avatar-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.headavatar{
  width: 100px;
  height: 100px;
  margin-top: 10px;
  margin-bottom: 20px;
  border-radius: 100%;
}

.avatar-large {
  width: 300px;
  height: 300px;
  border-radius: 50%;
}

.avatar-preview {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
}

.avatar-preview button {
  margin-top: 20px;
  padding: 10px 20px;
  font-size: 16px;
  color: #fff;
  background-color: #42b983;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.input_text {
  display: flex;
  align-items: center;
  width: 100%;
  max-width: 300px;
  margin-bottom: 20px;
}

.input_text label {
  margin-right: 20px;
  width: 50px;
  white-space: nowrap;
  text-align: left;
}

.input_text input[type="text"],
.input_text textarea {
  border: none;
  border-bottom: 1px solid #ccc;
  padding: 10px 0;
  font-size: 16px;
  outline: none;
  text-align: center;
}


.input_sig {
  display: flex;
  align-items: flex-start;
  width: 100%;
  max-width: 300px;
  margin-bottom: 20px;
}

.input_sig label {
  margin-right: 10px;
  width: 100px;
  white-space: nowrap;
  text-align: left;
}

.input_sig textarea {
  border: none;
  border-bottom: 1px solid #ccc;
  padding: 10px 0;
  font-size: 16px;
  outline: none;
  text-align: left;
  resize: none;
  height: 100px;
  margin-right: 20px;
}

.input_text textarea {
  resize: none;
  height: 100px;
}

.gender-options {
  display: flex;
  gap: 20px;
}

.gender-options label {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.gender-options input[type="radio"] {
  margin-right: 5px;
}

.date-picker {
  margin-top: 10px;
}

.signature-count {
  align-self: flex-end;
  font-size: 12px;
  color: #666;
}

.button_container{
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  max-width: 300px;
}

.save_button{
  margin-right: auto;
}

.cancle_button{
  margin-left:auto;
}

button {
  padding: 10px 10px;
  font-size: 16px;
  color: #fff;
  background-color: #42b983;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-bottom: 20px;
}

button:hover {
  background-color: #369f6d;
}

</style>