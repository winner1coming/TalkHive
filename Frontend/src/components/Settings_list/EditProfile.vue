<template>
  <!-- 编辑资料页面容器 -->
  <div class="edit">
    <div class="editprofile">
      <!-- 页面标题 -->
      <div class="avatar-container">
        <img :src="avatar" alt="avatar" class="headavatar" @click="showAvatarPreview" />
        <input type="file" ref="fileInput" style="display: none;" @change="handleFileChange" />
        <button v-if="isEditing" @click="openFilePicker">
          <img src="@/assets/icon/submit.png" alt="submit" class="icon"/>
        </button>
      </div>
      <div class="input_container">
        <div class="input_text" v-if="isEditing">
          <label for="id">ID:</label>
          <input id="id" type="text" v-model="id" :placeholder="id" :disabled="!isIdEditable"/>
          <span v-if="!isIdEditable" class="id-warning">一年只能更改一次-可更改时间{{ nextUpdateDate }}</span>
        </div>
        <div class="input_text" v-else>
          <label for="id">ID:</label>
          <span>{{ id }}</span>
        </div>

        <div class="input_text" v-if="isEditing">
          <label for="username">用户名:</label>
          <input id="username" type="text" v-model="username" :placeholder="username" />
        </div>
        <div class="input_text" v-else>
          <label for="username">用户名:</label>
          <span>{{ username }}</span>
        </div>

        <div class="input_text">
          <label>性别:</label>
          <span v-if="!isEditing">{{ gender == '男' ?'男':'女' }}</span>
          <div class="gender-options" v-if="isEditing">
            <label>
              <input type="radio" v-model="gender" value="男" />
              男
            </label>
            <label>
              <input type="radio" v-model="gender" value="女" />
              女
            </label>
          </div>
        </div>
        <div class="input_text" v-if="isEditing">
          <label>生日:</label>
          <input type="text" v-model="birthday" :placeholder="birthday" @click="showDatePicker" readonly />
          <div v-if="showDatePickerFlag" class="date-picker">
            <input type="date" v-model="birthday" @change="hideDatePicker" />
          </div>
        </div>
        <div class="input_text" v-else>
          <label>生日:</label>
          <span>{{ birthday }}</span>
        </div>

        <div class="input_text">
          <label>邮箱:</label>
          <span>{{ email }}</span>
        </div>

        <div class="input_text" v-if="isEditing">
          <label for="phone">手机号:</label>
          <input id="phone" type="text" v-model="phone" :placeholder="phone" />
        </div>
        <div class="input_text" v-else>
          <label for="phone">手机号:</label>
          <span>{{ phone }}</span>
        </div>

        <div class="input_sig" v-if="isEditing">
          <label>个性签名:</label>
          <textarea v-model="signature" :placeholder="signature" maxlength="100"></textarea>
          <span class="signature-count">{{ signature.length }}/100</span>
        </div>
        <div class="input_sig" v-else>
          <label>个性签名</label>
          <span>{{ signature }}</span>
        </div>
        <!-- 保存按钮，点击时触发 saveProfile 方法 -->
        <div class="button_container" v-if="isEditing">
        <button class="save_button" @click="saveProfile">保存</button>
        <!-- 取消按钮，点击时触发 cancelEdit 方法 -->
        <button class="cancle_button" @click="cancelEdit">取消</button>
        </div>
        <div class="botton_container" v-else>
          <button class="edit_button" @click="toggleEdit">编辑</button>
        </div>
      </div>

      <!-- 头像预览弹窗 -->
      <div v-if="showPreview" class="avatar-preview">
        <img :src="avatar" alt="avatar" class="avatar-large" />
        <button @click="hideAvatarPreview">
          <img src="@/assets/icon/cancel.png" alt="hide" class="icon"/>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { showProfile, saveEdit } from '@/services/settingView.js';
import { mapGetters } from 'vuex';

export default {
   // 从 Vuex 获取用户信息
  computed: {
    ...mapGetters(['user']),
  },
  // 组件的 data 函数，返回一个对象，包含组件的响应式数据
  data() {
    return {
      avatar:'',
      username: '',
      id: '',
      gender: '',
      birthday: '',
      signature: '',
      email:'',
      phone:'',
      lastUpdate:'',
      originalAvatar:'',
      originalUsername: '',
      originalId: '',
      originalGender: '',
      originalBirthday: '',
      originalSignature: '',
      originalPhone:'',
      showPreview: false,
      showDatePickerFlag: false,
      isEditing:false,
      isIdEditable: true, // 新增：控制账号是否可编辑
      nextUpdateDate: '',
    };
  },

  // 组件的生命周期钩子，在组件创建时调用
  created() {
    // 从数据库获取用户信息
    this.fetchProfile();
  },

  // 组件的方法定义
  methods: {
    // 从父组件取用户信息
    async fetchProfile() {
      try {
        const profile = await showProfile();
        if(profile.success){
          this.id = profile.data.id;
          this.avatar = this.user.avatar;
          this.username = profile.data.nickname;
          this.gender = profile.data.gender;
          this.birthday = profile.data.birthday;
          this.signature = profile.data.signature;
          this.email = profile.data.email;
          this.phone = profile.data.phone;
          this.lastUpdate = profile.data.lastUpdateID;

          // 保存原始用户名和 ID
          this.originalUsername = profile.data.nickname;
          this.originalId = this.user.id;
          this.originalGender = profile.data.gender;
          this.originalBirthday = profile.data.birthday;
          this.originalSignature = profile.data.signature;
          this.originalPhone = profile.data.phone;
          this.originalAvatar = this.user.avatar;

          //检查ID是否可以修改
          const currentDate = new Date();
          const timeDifference = currentDate - this.lastUpdate;
          const timegap = 365*24*60*60*1000;
          
          if(timeDifference > timegap){
            this.isIdEditable = false;
            const nextUpdateDate = new Date(this.lastUpdate.getTime() + timegap);
            this.nextUpdateDate = nextUpdateDate.toLocaleDateString();
          }
          else{
            this.isIdEditable = true;
            this.nextUpdateDate = '';
          }
        }
        else{
          alert('获取个人主页数据失败'+response.message);
        }
      } catch (error) {
        console.error('加载数据失败:', error);
      }
    },

    // 保存资料方法，处理用户点击保存按钮时的逻辑
    async saveProfile() {
      try {
        //检查账号是否改变
        if(this.id !== this.originalId){
          const currentDate =new Date();
          this.lastUpdate = currentDate.toISOString();
        }

        await saveEdit({
          nickname: this.username,
          id: this.id,
          avatar:this.avatar,
          gender: this.gender,
          birthday: this.birthday,
          signature: this.signature,
          phone:this.phone,
          lastUpdateID: this.lastUpdate,
        });
        // 更新原始用户名和 ID
        this.originalUsername = this.username;
        this.originalId = this.id;
        this.originalGender = this.gender;
        this.originalBirthday = this.birthday;
        this.originalSignature = this.signature;
        this.originalPhone = this.phone,
        this.originalAvatar = this.avatar;

        //更新全局变量
        this.$store.commit('SET_USER', {
            username:this.username,
            avatar : this.avatar,
        });

        //更新缓存的内容
        let users = JSON.parse(localStorage.getItem('users')) || [];
        const userInfo = {
          account: this.id,
          avatar : this.avatar,
        }

        const index = users.findIndex(user => user.account === this.id);
        if(index !== -1){
          users[index].account = userInfo.account;
          users[index].avatar = userInfo.avatar;
        }else{
          users.push(userInfo);
        }

        localStorage.setItem('users',JSON.stringify(users));

        this.isEditing =false;
      } catch (error) {
        console.error('Failed to save profile:', error);
      }
    },

    // 取消编辑方法，处理用户点击取消按钮时的逻辑
    cancelEdit() {
      // 恢复原始用户名和 ID
      this.username = this.originalUsername;
      this.id = this.originalId;
      this.avatar = this.originalAvatar;
      this.gender = this.originalGender;
      this.birthday = this.originalBirthday;
      this.signature = this.originalSignature;
      this.phone = this.originalPhone;
      this.isEditing=false;
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

    //切换编辑模式
    toggleEdit(){
      this.isEditing = !this.isEditing;
    }
  },
};
</script>

<style scoped>

.edit{
  display: flex;
  justify-content: center;
  height: 100%;
  background-color: var(--background-color);
}

/* 编辑资料页面的样式 */
.editprofile {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
  align-items: center;
  align-content: center;
  justify-items: center;
  margin-top: 60px;
  width: 100%;
  height: 70%;
  max-width: 600px;
  border-radius: 8px;
  box-shadow: 0 2px 10px var(--background-color1);
  background-color: var(--background-color);
}

h2{
  align-items: center;
  margin-top: 10px;
}

.avatar-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 10px;
}

.avatar-container button{
  background-color: var(--button-background-color);
  width: 40px;
  height: 30px;
  margin-top: 0px;
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
  margin-top: 30px;
  padding: 6px 10px;
  font-size: var(--font-size);
  color: var(--button-text-color);
  background-color: var(--button-background-color);
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.input_container {
  width: 300px;
  margin-bottom: 10px;
}

.input_text {
  display: flex;
  align-items: center;
  width: 100%;
  max-width: 300px;
  margin-bottom: 20px;
  color: var(--text-color);
  font-size: var(--font-size);
}

.input_text label {
  margin-right: 60px;
  width: 50px;
  white-space: nowrap;
  text-align: left;
}

.input_text input[type="text"],
.input_text textarea {
  border: none;
  border-bottom: 1px solid #ccc;
  padding: 10px 0;
  font-size: var(--font-size);
  color: var(--text-color);
  outline: none;
  text-align: center;
  background-color: var(--background-color);
}


.input_sig {
  display: flex;
  align-items: flex-start;
  width: 100%;
  max-width: 300px;
  margin-bottom: 20px;
  color: var(--text-color);
  font-size:var(--font-size);
}

.input_sig label {
  margin-right: 10px;
  width: 100px;
  white-space: nowrap;
  text-align: left;
}

.input_sig textarea {
  border-bottom: 1px solid #ccc;
  padding: 8px 8px;
  font-size: var(--font-size);
  outline: none;
  text-align: left;
  resize: none;
  height: 100px;
  width: 400px;
  margin-right: -55px;
  color: var(--text-color);
  background-color: var(--background-color);
}

.input_text textarea {
  resize: none;
  height: 100px;
  color: var(--text-color);
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
  font-size: var(--font-size-small);
  color: var(--text-color);
}

.button_container{
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: var(--button-text-color);
  font-size: var(--font-size);
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
  padding: 6px 10px;
  font-size: var(--font-size);
  color:var(--button-text-color);
  background-color: var(--button-background-color);
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-bottom: 20px;
  margin-top: 10px;
}

button:hover ,.edit_button:hover{
  background-color: var(--button-background-color1);
}


.icon{
  width:30px;
  height: 30px;
  margin-left:-5px;
  margin-top: -10px;
}
</style>
