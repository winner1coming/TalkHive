<template>
    <div class="chat-background-settings">
      <h3>聊天背景</h3>
      <div class="current-background">
        <img :src="currentBackground" alt="Current Background" />
      </div>
      <div class="upload-container">
        <label class="custom-file-upload">
          <input type="file" accept="image/*" @change="onFileChange" />
          上传
        </label>
      </div>
      <button @click="saveBackground">保存</button>
    </div>
  </template>
  
  <script>
  import currentBackground from '@/assets/images/avatar.jpg';
  
  export default {
    data() {
      return {
        currentBackground, // 默认背景图片
        newBackground: null, // 新上传的背景图片
      };
    },
    methods: {
      onFileChange(event) {
        const file = event.target.files[0];
        if (file) {
          const reader = new FileReader();
          reader.onload = (e) => {
            this.currentBackground = e.target.result;
            this.newBackground = e.target.result;
          };
          reader.readAsDataURL(file);
        }
      },
      saveBackground() {
        // 保存背景图片逻辑
        if (this.newBackground) {
          console.log('Background saved:', this.currentBackground);
          // 你可以在这里调用一个方法来保存背景图片
        } else {
          console.log('No new background to save');
        }
      },
    },
  };
  </script>
  
  <style scoped>
  .chat-background-settings {
    padding: 20px;
    width: 100%;
    max-width: 400px;
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
    max-width: 300px;
    border-radius: 8px;
    overflow: hidden;
  }
  
  .current-background img {
    width: 100%;
    height: auto;
    display: block;
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