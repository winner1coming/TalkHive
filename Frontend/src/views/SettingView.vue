<template>
    <div class="setting">
      <div class="top-panel">
        <button
          :class="{ active: activeComponent === 'EditProfile' }"
          @click="setActiveComponent('EditProfile')"
        >
          个人主页
        </button>
        <button
          :class="{ active: activeComponent === 'SecuritySettings' }"
          @click="setActiveComponent('SecuritySettings')"
        >
          安全设置
        </button>
        <button
          :class="{ active: activeComponent === 'SystemSettings' }"
          @click="setActiveSystemSetting('SystemSettings')"
        >
          系统设置
        </button>
        <button @click="showLogoutConfirmation">退出登录</button>
      </div>
      <div class="bottom-panel">
        <component :is="activeComponent"></component>
      </div>
      <div v-if="showConfirmation" class="confirmation-modal">
        <div class="modal-content">
          <span class="close" @click="hideLogoutConfirmation">&times;</span>
          <p>是否确认退出登录？</p>
          <div class="modal-buttons">
            <button @click="confirmLogout">确认</button>
            <button @click="hideLogoutConfirmation">取消</button>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import EditProfile from '@/components/Settings_list/EditProfile.vue';
  import SecuritySettings from '@/components/Settings_list/SecuritySettings.vue';
  import SystemSettings from '@/components/Settings_list/SystemSettings.vue';
  import {logout, getSystemSetting} from '@/services/settingView.js';
  import { mapGetters } from 'vuex';
  
  export default {
    components: {
      EditProfile,
      SecuritySettings,
      SystemSettings,
    },

    computed: {
      ...mapGetters(['user']),
    },

    data() {
      return {
        activeComponent: 'EditProfile', // 默认显示个人主页
        showConfirmation: false,
      };
    },

    methods: {
      setActiveComponent(component){
        this.activeComponent = component;
      },
      async setActiveSystemSetting(component) {
        this.activeComponent = component;
        try{
          const response = await getSystemSetting();
          if(response.success){
            this.store.commit('SET_SETTINGS',{
              theme:response.data.theme,
              fontSize:response.data.fontSize,
              fontStyle:fontStyle,
              sound:response.data.sound,
              isNotice: response.data.notice,
              isNoticeGroup: response.data.noticeGroup,
              background:response.data.background,
            });
          }else{
            alert(response.message);
          }
        }catch(error){
          alert(error,'获取系统失败，请检查网络')
        }
      },
      showLogoutConfirmation() {
        this.showConfirmation = true;
      },
      hideLogoutConfirmation() {
        this.showConfirmation = false;
      },
      async confirmLogout() {
        try{
          const response = await logout();
          if(response.success){
            alert('已退出登录~');
            this.$router.push('/loginth');
            this.hideLogoutConfirmation();
          // 你可以在这里添加退出登录的逻辑
          }else{
            alert("退出登录失败，请重试！");
          }
        }catch(error){
          console.error("退出登录失败！");
        }
      },
    },
  };
  </script>
  
  <style scoped>
  .setting {
    display: flex;
    flex-direction: column;
    height: 100vh;
  }
  
  .top-panel {
    display: flex;
    justify-content: space-around;
    align-items: center;
    padding: 10px;
    background-color: #f0f0f0;
    border-bottom: 1px solid #ccc;
  }
  
  .top-panel button {
    padding: 10px 20px;
    background-color: #ccc;
    color: black;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .top-panel button.active {
    background-color: #42b983;
    color: white;
  }
  
  .top-panel button:hover {
    background-color: #369f6e;
    color: white;
  }
  
  .bottom-panel {
    flex: 1;
    overflow-y: auto;
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
    font-size: 20px;
  }
  
  .modal-buttons {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
  
  .modal-buttons button {
    margin-left: 10px;
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .modal-buttons button:first-child {
    background-color: #42b983;
    color: white;
  }
  
  .modal-buttons button:last-child {
    background-color: #ccc;
    color: black;
  }
  </style>