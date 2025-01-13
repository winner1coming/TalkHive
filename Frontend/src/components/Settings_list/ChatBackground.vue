<template>
    <div class="chat-background-settings">
      <div class="background-view">
      <h3>聊天背景</h3>
      <div class="current-background">
        <div v-if="!currentBackground" class="placeholder">
          <span>默认背景</span>
        </div>
        <img v-else :src="currentBackground" alt="Current Background" />
      </div>
      <div class="upload-container">
        <label class="custom-file-upload">
          <input type="file" accept="image/*" @change="onFileChange" />
          <img src="@/assets/icon/submit.png" alt="Submit" class="icon"/>
        </label>
      </div>
      <button @click="saveBackground">保存</button>
      <Windows 
        :visible="showModal"
        :message="modalMessage"
        @close="showModal = false"
      />
      </div>
    </div>
  </template>
  
  <script>
  import Windows from '@/components/base/Windows.vue';
  import {changeBackground} from '@/services/settingView.js';
  export default {
    components:{
      Windows,
    },
    data() {
      return {
        currentBackground:this.$store.state.settings.background, // 默认背景图片
        newBackground: null, // 新上传的背景图片
        background:null,
        showModal:false,
        modalMessage:'',
      };
    },
    methods: {
      onFileChange(event) {
        const file = event.target.files[0];
        if (file) {
          const reader = new FileReader();
          reader.onload = (e) => {
            this.background = `data:${file.type};base64,${e.target.result.split(',')[1]}`;
            this.currentBackground = e.target.result;
            this.newBackground = this.currentBackground;
          };
          reader.readAsDataURL(file);
        }
      },
      async saveBackground() {
        if(!this.newBackground){
          this.modalMessage = '请先上传图片!';
          this.showModal = true;
          return;
        }
        // 保存背景图片逻辑
        try{
          const response = await changeBackground({background:this.newBackground});
          if(response.success){
            this.$store.commit('SET_BACKGROUND',this.background);
            this.modalMessage = response.message;
            this.showModal = true;
          }else{
            this.modalMessage = response.message;
            this.showModal = true;
          }
        }catch(error){
            this.modalMessage = error;
            this.showModal = true;
        }
      },
    },
  };
  </script>
  
  <style scoped>
  .chat-background-settings {
    padding: 20px;
    margin-top: 60px;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .background-view{
    padding: 20px;
    width: 100%;
    max-width: 500px;
    height: 400px;
    background-color: var(--background-color);
    border-radius: 8px;
    box-shadow: 0 2px 10px var(--background-color2);
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  h3 {
    margin-bottom: 20px;
    color: #000;
    background-color: var(--sidebar-background-color1);
    color: var(--sidebar-text-color);
    font-size: var(--font-size-large);
    padding: 10px;
    text-align: center;
    border-radius: 4px;
    width: 100%;
  }
  
  .current-background {
    margin-bottom: 10px;
    width: 100%;
    max-width: 400px;
    border-radius: 8px;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 200px; /* 设置最小高度 */
    box-shadow: 0 4px 8px var(--background-color2); /* 添加阴影效果 */
    background-color: var(--background-color); /* 设置背景颜色 */
    color: var(--text-color);
  }
  
  .current-background img {
    width: 100%;
    height: auto;
    display: block;
  }

  .placeholder {
  text-align: center;
  color: var(--text-color); /* 占位文字颜色 */
  opacity: 60%;
  font-size: var(--font-size-small);
  }
  
  .upload-container {
    display: flex;
    align-items: center;
  }
  
  .custom-file-upload {
    padding: 4px 6px;
    color: var(--background-color);
    border: none;
    border-radius: 0px;
    cursor: pointer;
    display: inline-block;
    margin-bottom: 10px;
  }
  
  .custom-file-upload:hover {
    background-color: var(--button-background-color1);
  }
  
  input[type="file"] {
    display: none; /* 隐藏原生的文件选择按钮 */
  }

  .icon{
    width: 35px;
    height: 35px;
  }
  
  button {
    padding: 10px 20px;
    background-color: var(--button-background-color);
    color:var(--button-text-color);
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  button:hover {
    background-color: var(--button-background-color1);
  }
  </style>