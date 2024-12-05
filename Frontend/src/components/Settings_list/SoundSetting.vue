<template>
    <div class="notification-settings">
      <h2>消息通知</h2>
      <div class="setting-item">
        <span>开启消息通知</span>
        <button class="toggle-button" :class="{ 'status-on': notificationStatus === 'on' }" @click="toggleNotificationStatus">
          <span class="toggle-circle"></span>
        </button>
      </div>
      <div class="setting-item">
        <span>开启声音</span>
        <button class="toggle-button" :class="{ 'status-on': soundStatus === 'on' }" @click="toggleSoundStatus">
          <span class="toggle-circle"></span>
        </button>
      </div>
      <div class="setting-item">
        <span>接收群消息声音</span>
        <button class="toggle-button" :class="{ 'status-on': groupSoundStatus === 'on' }" @click="toggleGroupSoundStatus">
          <span class="toggle-circle"></span>
        </button>
      </div>
      <div class="setting-item">
        <span>选择提示音</span>
        <select v-model="selectedSound" @change="applySound">
          <option v-for="sound in sounds" :key="sound.value" :value="sound.value">{{ sound.label }}</option>
        </select>
      </div>
      <div class="setting-item">
        <span>上传音频文件</span>
        <input type="file" accept="audio/*" @change="onFileChange" />
      </div>
      <button @click="saveSettings">保存</button>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        notificationStatus: 'off', // 'on' or 'off'
        soundStatus: 'off', // 'on' or 'off'
        groupSoundStatus: 'off', // 'on' or 'off'
        selectedSound: '',
        sounds: [
          { label: '默认提示音', value: 'default.mp3' },
          { label: '提示音1', value: 'sound1.mp3' },
          { label: '提示音2', value: 'sound2.mp3' },
        ],
        newSound: null, // 新上传的音频文件
      };
    },
    methods: {
      toggleNotificationStatus() {
        this.notificationStatus = this.notificationStatus === 'on' ? 'off' : 'on';
      },
      toggleSoundStatus() {
        this.soundStatus = this.soundStatus === 'on' ? 'off' : 'on';
      },
      toggleGroupSoundStatus() {
        this.groupSoundStatus = this.groupSoundStatus === 'on' ? 'off' : 'on';
      },
      applySound() {
        // 你可以在这里处理应用提示音的逻辑
        console.log('Selected sound:', this.selectedSound);
      },
      onFileChange(event) {
        const file = event.target.files[0];
        if (file) {
          const reader = new FileReader();
          reader.onload = (e) => {
            this.newSound = e.target.result;
          };
          reader.readAsDataURL(file);
        }
      },
      saveSettings() {
        // 保存设置逻辑
        console.log('Settings saved:', {
          notificationStatus: this.notificationStatus,
          soundStatus: this.soundStatus,
          groupSoundStatus: this.groupSoundStatus,
          selectedSound: this.selectedSound,
          newSound: this.newSound,
        });
        // 你可以在这里调用一个方法来保存设置
      },
    },
  };
  </script>
  
  <style scoped>
  .notification-settings {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px;
    width: 100%;
    max-width: 400px;
    background-color: #f9f9f9;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  }
  
  h2 {
    margin-bottom: 20px;
    color: #000;
    background-color: rgb(173, 229, 210);
    padding: 10px;
    text-align: center;
    border-radius: 4px;
    width: 100%;
  }
  
  .setting-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 0;
    width: 100%;
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
  
  select {
    padding: 5px;
    border-radius: 4px;
    border: 1px solid #ccc;
  }
  
  input[type="file"] {
    padding: 5px;
    border-radius: 4px;
    border: 1px solid #ccc;
  }
  
  button {
    padding: 10px 20px;
    background-color: #42b983;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    margin-top: 20px;
  }
  
  button:hover {
    background-color: #369f6e;
  }
  </style>