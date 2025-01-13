<template>
  <div class="security-settings">
    <div class="left-panel">
      <div class="menu-item" :class="{ active: activeComponent === 'ChangeEmail' }" @click="setActiveComponent('ChangeEmail')">
        <img src="@/assets/icon/email.png" alt="ChangeEmail" class="icon"/>
        <span>邮箱</span>
        <span class="content">{{ user.email }}</span>
      </div>
      <div class="menu-item" :class="{ active: activeComponent === 'ChangePassword' }" @click="setActiveComponent('ChangePassword')">
        <img src="@/assets/icon/change_password.png" alt="ChangePassword" class="icon"/>
        <span>更改密码</span>
        <span class="content"></span>
      </div>
      <div class="menu-item" :class="{ active: activeComponent === 'FriendPermission' }" @click="setActiveComponent('FriendPermission')">
        <img src="@/assets/icon/permission.png" alt="FriendPermission" class="icon"/>
        <span>好友权限设置</span>
        <span class="content"></span>
      </div>
      <div class="menu-item" :class="{ active: activeComponent === 'DeactivateAccount' }" @click="showDeactivateConfirmation">
        <img src="@/assets/icon/deativate.png" alt="ChangeEmail" class="icon"/>
        <span>注销账号</span>
        <span class="content"></span>
      </div>
    </div>
    <div class="resizer" @mousedown="startResize"></div>
    <div class="right-panel">
      <component :is="activeComponent" :user="users" @updateUser="updateUser"></component>
      <img v-if="this.$store.state.settings.theme === 'light' && activeComponent === ''" src="@/assets/icon/light_wel.png" alt="Tip" class="icon" />
      <img v-if="this.$store.state.settings.theme === 'dark' && activeComponent === ''" src="@/assets/icon/dark_wel.png" alt="Tip" class="icon" />
    </div>
    <div v-if="showConfirmation" class="confirmation-modal">
      <div class="modal-content">
        <span class="close" @click="hideDeactivateConfirmation">&times;</span>
        <p>账号一旦注销，本用户信息将被销毁！</p>
        <p>是否选择注销账号？</p>
        <div class="modal-buttons">
          <button @click="confirmDeactivate">确认</button>
          <button @click="hideDeactivateConfirmation">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import ChangePassword from './ChangePassword.vue';
import FriendPermission from './FriendPermission.vue';
import ChangeEmail from './ChangeEmail.vue';
import { mapGetters } from 'vuex';
import { getUserInfo , confirmDeactivation } from '@/services/settingView.js';

export default {
  components: {
    ChangePassword,
    FriendPermission,
    ChangeEmail,
  },

  computed: {
    ...mapGetters(['user']),
  },

  data() {
    return {
      users: {
        ID:'',
        email:'',
        password:'',
        friend_permissionID:'off',
        friend_permissionNickname:'off',
      },
      activeComponent: '',
      showConfirmation:false,
    };
  },
  mounted(){
    this.fetchUserInfo();
    this.setActiveComponent('');
  },

  methods: {
    async fetchUserInfo(){
      try{
        const response = await getUserInfo();
        if(response.success){
          this.users.ID = response.data.id;
          this.users.email = response.data.email;
          this.users.password =  response.data.password;
          this.users.friend_permissionID = response.data.friend_permissionID ? 'on':'off';
          this.users.friend_permissionNickname = response.data.friend_permissionNickname?'on':'off';
        }
        else{
          alert(response.message || '获取用户邮箱失败');
        }
      }catch(error){
          console.error('获取信息失败:',error);
      }
    },

    setActiveComponent(component) {
      this.activeComponent = component;
    },
    updateUser(updatedUser) {
      this.users = { ...this.users, ...updatedUser };
      this.setActiveComponent('');  
    },
    showDeactivateConfirmation(){
      this.showConfirmation = true;
    },
    hideDeactivateConfirmation(){
      this.showConfirmation = false;
    },
    async confirmDeactivate(){
      // 注销账号的逻辑
      try{
        const response = await confirmDeactivation();
        if(response.success){
          alert('账号已注销');
          this.hideDeactivateConfirmation();
          this.$router.push('/loginth');
        }
        else{
          alert(response.message || '注销账号失败');
        }

      }catch (error){
        console.error("账号注销失败:",error)
      }

    },
    startResize(event) {
      this.isResizing = true;
      document.addEventListener('mousemove', this.resize);
      document.addEventListener('mouseup', this.stopResize);
    },
    resize(event) {
      if (this.isResizing) {
        this.chatListWidth = event.clientX - this.leftComponentWidth;
      }
    },
    stopResize() {
      this.isResizing = false;
      document.removeEventListener('mousemove', this.resize);
      document.removeEventListener('mouseup', this.stopResize);
    },
  },
};
</script>

<style scoped>
.security-settings {
  display: flex;
  height: 100%;
}

.left-panel {
  width: 14%;
  background-color: var(--background-color1);
}

.right-panel {
  width: 100%;
  padding:0;
  position: relative;
  background-color: var(--background-color);
}

.right-panel .icon{
  height: 100px;
  margin-top: 300px;
}

.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #ccc;
  cursor: pointer;
  width: 100%;
  height: 10vh;
}

.menu-item .icon{
  width: 25px;
  height: 25px;
  margin-left: 30px;
}

.menu-item.active {
  background-color: var(--select-background-color1);
  color: var(--select-text-color);
}

.menu-item:hover{
  background-color: var(--select-background-color);
  opacity: 70%;
}

.menu-item span {
  font-size: var(--font-size);
  margin-left: 5px;
}

.menu-item .content {
  font-size: var(--font-size-small);
  color: #666;
}

.menu-item.active .content {
  color: var(--select-text-color);
  font-weight: bold;
}

.menu-item.active,.menu-item:hover .span{
  font-weight: bold;
  color: var(--select-text-color);
}

.menu-item.active:hover {
  background-color: var(--select-background-color1); /* 保持点击后的背景颜色 */
  opacity: 100%; /* 如果需要，可以调整透明度 */
}

.confirmation-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background-color: var(--background-color);
  font-size: var(--font-size);
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
  color: var(--button-text-color);
  margin-top: 20px;
}

.modal-buttons button {
  margin-left: 10px;
  padding: 6px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.modal-buttons button:first-child {
  background-color: var(--button-background-color);
  color: var(--button-text-color);
}

.modal-buttons button:last-child {
  background-color: var(--background-color2);
  color:var(--button-text-color);
}

.modal-buttons button:hover{
  background-color: var(--button-background-color1);
}

.back-button {
  position: absolute;
  top: 10px;
  left: 10px;
  cursor: pointer;
  font-size: var(--font-size);
  color: var(--button-background-color1);
  display: flex;
  align-items: center;
}

.back-button i {
  margin-right: 5px;
}

.resizer {
  width: 3px;
  height: 100%;
  cursor: ew-resize;
  background-color: var(--background-color2);
}
</style>