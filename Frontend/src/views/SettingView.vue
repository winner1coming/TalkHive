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
  import {logout} from '@/services/settingView.js';
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
      setActiveSystemSetting(component) {
        this.activeComponent = component;
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
    justify-content: space-around; /* 均匀分布按钮 */
    align-items: center; /* 垂直居中 */
    padding: 10px;
    background-color: var(--sidebar-background-color1); /* 使用 CSS 变量 */
    border-bottom: 1px solid #ccc;
    border-radius: 10px;
    width: 100%;
    height: 35px; /* 固定高度 */
    flex-wrap: nowrap;
    opacity: 100%;
  }
  
  .top-panel button {
    width: fit-content; /* 固定按钮宽度 */
    height: 40px; /* 固定按钮高度 */
    padding: 8px;
    background-color: var(--button-background-color); /* 使用 CSS 变量 */
    color: var(--button-text-color); /* 使用 CSS 变量 */
    border: none;
    border-radius: 4px;
    cursor: pointer;
    text-align: center; /* 文字居中 */
    font-size: var(--font-size); /* 使用 CSS 变量设置字体大小 */
  }
  
  .top-panel button.active {
    background-color: var(--button-background-color2);
    color: var(--button-text-color);
    font-weight: bold;
  }
  
  .top-panel button:hover {
    background-color:var(--button-background-color1) ;
    color: var(--button-text-color);
    opacity: 80%;
  }

  .top-panel button.active:hover {
    background-color: var(--button-background-color2); /* 保持点击后的背景颜色 */
    opacity: 100%; /* 如果需要，可以调整透明度 */
    font-weight: bold;
  }

  .top-panel button.active ::content{
    color: var(--button-text-color);
    color: #76cea7d4;
  }
  
  .bottom-panel {
    flex: 1;
    overflow-y: auto;
    background-color: var(--background-color);
  }
  
  .confirmation-modal {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%); /* 确保模态框居中 */
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
  }
  
  .modal-content {
    background-color: var(--background-color);
    color: var(--sidebar-text-color);
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
    margin-top: 30px;
    font-size: var(--font-size);
  }
  
  .modal-buttons button {
    margin-left: 10px;
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    color: var(--button-text-color);
  }
  
  
  .modal-buttons button:first-child {
    background-color: var(--button-background-color2);
    color: var(--button-text-color);
  }
  
  .modal-buttons button:last-child {
    background-color:var(--background-color2);
    color: var(--button-text-color);
  }

  .modal-buttons button:hover{
    background-color: var(--button-background-color);
  }

  .modal-buttons button:active{
    background-color: var(--button-background-color2);
  }

  .modal-content p{
    color: var(--text-color);
    font-size: var(--font-size);
  }
  </style>
