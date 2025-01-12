<template>
    <div class="chat-background-settings">
      <div class="background-view">
      <h3>聊天背景</h3>
      <div class="current-background">
        <div v-if="!currentBackground" class="placeholder">
          <span>默认纯白背景</span>
        </div>
        <img v-else :src="currentBackground" alt="Current Background" />
      </div>
      <div class="upload-container">
        <label class="custom-file-upload">
          <input type="file" accept="image/*" @change="onFileChange" />
          上传
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
    height: 50vh;
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .background-view{
    padding: 20px;
    width: 100%;
    max-width: 500px;
    height: 400px;
    background-color: #f9f9f9;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  h3 {
    margin-bottom: 20px;
    color: #000;
    background-color: rgb(173, 229, 210);
    padding: 10px;
    text-align: center;
    border-radius: 4px;
    width: 100%;
  }
  
  .current-background {
    margin-bottom: 20px;
    width: 100%;
    max-width: 400px;
    border-radius: 8px;
    overflow: hidden;
    border: 2px dashed #ccc; /* 添加虚线边框 */
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 200px; /* 设置最小高度 */
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); /* 添加阴影效果 */
    background-color: #fff; /* 设置背景颜色 */
  }
  
  .current-background img {
    width: 100%;
    height: auto;
    display: block;
  }

  .placeholder {
  text-align: center;
  color: #888; /* 占位文字颜色 */
  font-size: var(--font-size-small);
  }
  
  .upload-container {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
  }
  
  .custom-file-upload {
    padding: 10px 20px;
    background-color: #42b983;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    display: inline-block;
  }
  
  .custom-file-upload:hover {
    background-color: #369f6e;
  }
  
  input[type="file"] {
    display: none; /* 隐藏原生的文件选择按钮 */
  }
  
  button {
    padding: 10px 20px;
    background-color: #42b983;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  button:hover {
    background-color: #369f6e;
  }
  </style>