<template>
  <div class="sound_set">
    <div class="notification-settings">
      <h2>消息通知</h2>
      <div class="setting-item">
        <span>开启消息通知</span>
        <button class="toggle-button" :class="{ 'status-on': notificationStatus}" @click="toggleNotificationStatus">
          <span class="toggle-circle"></span>
        </button>
      </div>
      <div class="setting-item">
        <span>接收群消息声音</span>
        <button class="toggle-button" :class="{ 'status-on': groupSoundStatus}" @click="toggleGroupSoundStatus">
          <span class="toggle-circle"></span>
        </button>
      </div>
      <div class="setting-item">
        <span>选择提示音</span>
        <select v-model="selectedSound" @change="applySound">
          <option v-for="sound in sounds" :key="sound.value" :value="sound.value">{{ sound.label }}</option>
        </select>
        <button @click="saveSound">更改</button>
      </div>
      <Windows 
      :visible="showModal"
      :message="modalMessage"
      @close="showModal = false"
    />
    </div>
  </div>
  </template>
  
  <script>
  import bird from '@/assets/sounds/bird.mp3';
  import dingdo from '@/assets/sounds/dingdo.mp3';
  import huawei from '@/assets/sounds/huawei.mp3';
  import skype from '@/assets/sounds/skype.mp3';
  import wechat from '@/assets/sounds/wechat.mp3';
  import { isNotice,isNoticeGroup,changeSound } from '@/services/settingView.js';
  import Windows from '@/components/base/Windows.vue';

  export default {
    components:{
      Windows,
    },
    data() {
      return {
        selectedSound: this.$store.state.settings.sound,
        sounds: [
          { label: '默认提示音', value: 'dingdo.mp3', path:dingdo },
          { label: '清新鸟鸣', value: 'bird.mp3', path:bird },
          { label: '华为提示音', value: 'huawei.mp3',path:huawei },
          { label: '微信提示音', value: 'wechat.mp3',path:wechat },
          { label: '气泡提示音', value: 'skype.mp3',path:skype },
        ],
        audio: null, // 新上传的音频文件
        showModal:false,
        modalMessage:'',
        notificationStatus:this.$store.state.settings.isNotice,
        groupSoundStatus:this.$store.state.settings.isNoticeGroup,
      };
    },
    methods: {
      async toggleNotificationStatus() {
        //开启消息通知
        try{
          const temp = !this.notificationStatus;
          const response = await isNotice({notice:temp});
          if(response.success){
            this.notificationStatus = temp;
            this.$store.commit('SET_NOTICE',this.notificationStatus);
          }else{
            this.showModal = true;
            this.modalMessage = response.message;
          }

        }catch(error){
          this.showModal = true;
          this.modalMessage = '设置消息通知失败，请检查网络';
        }
      },
      async toggleGroupSoundStatus() {
        try{
          const temp = !this.groupSoundStatus;
          const response = await isNoticeGroup({noticeGroup:temp});
          if(response.success){
            this.groupSoundStatus = temp;
            this.$store.commit('SET_GROUPNOTICE',this.groupSoundStatus);
          }else{
            this.showModal = true;
            this.modalMessage = response.message;
          }
        }catch(error){
          this.showModal = true;
          this.modalMessage = '设置群聊消息通知失败,请检查网络';
        }
      },
      async applySound() {
        // 你可以在这里处理应用提示音的逻辑
        if(this.audio){
          this.audio.pause();
        }

        const select = this.sounds.find(sound=>sound.value === this.selectedSound);
        if(select){
          this.audio = new Audio(select.path);
          this.audio.play();

          setTimeout(()=>{
            if(this.audio){
              this.audio.pause();
            }
          },8000);
        }
      },
      async saveSound() {
        // 保存设置逻辑
        try{
          const response = await changeSound({sound:this.selectedSound});
          this.showModal = true;
          if(response.success){
            this.$store.commit('SET_SOUND',this.selectedSound);
          }
          this.modalMessage = response.message;
        }catch(error){
          this.showModal = true;
          this.modalMessage = '通知音保持失败，请检查网络';
        }
        
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
    align-items: flex-start;
    gap: 10px;
    padding: 10px 0;
    width: 100%;
    justify-content: space-between;
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

  .sound_set{
    padding: 20px;
    height: 50vh;
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  </style>