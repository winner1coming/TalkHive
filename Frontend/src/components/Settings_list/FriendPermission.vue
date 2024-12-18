<template>
    <div class="permission">
      <h2>添加我的方式</h2>
      <div class="user-info">
        <span class="user-id">ID: {{ ID }}</span>
        <button class="toggle-button" :class="{ 'status-on': idStatus === 'on' }" @click="toggleIdStatus">
          <span class="toggle-circle"></span>
        </button>
      </div>
      <div class="user-info">
        <span class="user-phone">手机号: {{ phone }}</span>
        <button class="toggle-button" :class="{ 'status-on': phoneStatus === 'on' }" @click="togglePhoneStatus">
          <span class="toggle-circle"></span>
        </button>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        ID: '阳光小青年',
        phone: '18682139298',
        idStatus: 'off', // 'on' or 'off'
        phoneStatus: 'off', // 'on' or 'off'
      };
    },
    methods: {
      toggleIdStatus() {
        this.idStatus = this.idStatus === 'on' ? 'off' : 'on';
      },
      togglePhoneStatus() {
        this.phoneStatus = this.phoneStatus === 'on' ? 'off' : 'on';
      },
    },
  };
  </script>
  
  <style scoped>
  .permission {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .user-info {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 20px;
  }
  
  .user-id, .user-phone {
    font-size: 16px;
  }
  
  .toggle-button {
    width: 60px;
    height: 30px;
    border-radius: 15px; /* 调整为圆角矩形 */
    background-color: gray; /* 默认关闭状态 */
    border: none;
    cursor: pointer;
    position: relative;
    transition: background-color 0.3s;
  }
  
  .toggle-button.status-on {
    background-color: green; /* 开启状态 */
  }
  
  .toggle-circle {
    width: 26px;
    height: 26px;
    border-radius: 50%;
    background-color: white;
    position: absolute;
    top: 2px;
    left: 2px;
    transition: transform 0.3s;
  }
  
  .toggle-button.status-on .toggle-circle {
    transform: translateX(30px); /* 滑动到右边 */
  }
  </style>