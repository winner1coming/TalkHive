<template>
    <div class="font-size-settings">
      <div class="font-view">
      <h3>字体大小</h3>
        <div class="preview">
            <div class="message-box">
            <p :style="{ fontSize: `${fontSize}px` }">预览字体大小</p>
            </div>
            <div class="avatar">
                <img :src="avatar">
            </div>
        </div>
        <div class="slider-container">
            <input
            type="range"
            min="12"
            max="36"
            v-model="fontSize"
            @input="onFontSizeChange"
            />
            <span>{{ fontSize }}px</span>
        </div>
        <div class="description">
            <p>拖动下面的滑块，可设置字体大小</p>
        </div>
        <button @click="saveFontSize">完成</button>
      </div>
    </div>
  </template>

  
  <script>

  import avatar from '@/assets/images/avatar.jpg';
  export default {
    data() {
      return {
        avatar,
        fontSize: 16, // 默认字体大小
      };
    },
    methods: {
      onFontSizeChange() {
        // 处理字体大小变化逻辑
        console.log('Font size changed:', this.fontSize);
      },
      saveFontSize() {
        // 保存字体大小逻辑
        console.log('Font size saved:', this.fontSize);
        // 你可以在这里调用一个方法来保存字体大小
      },
    },
  };
  </script>
  
  <style scoped>
  .font-size-settings {
    padding: 20px;
    height: 50vh;
    width: 100%;
    width: 400px;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .font-view {
  padding: 20px;
  width: 100%;
  max-width: 400px; /* 设置最大宽度 */
  background-color: #f9f9f9; /* 添加背景色 */
  border-radius: 8px; /* 添加圆角 */
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1); /* 添加阴影 */
  display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
}
  
  h3 {
    margin-bottom: 20px;
    color: #000;
    background-color: rgb(173, 229, 210);
    width: 300px;
    border-radius: 8px; /* 添加圆角 */
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1); /* 添加阴影 */
  }
  
  .preview {
    margin-bottom: 20px;
    text-align: center;
    align-items: flex-end;
    display: flex;
    margin-left: 80px;
  }

  .message-box {
  background-color: #e0f7fa; /* 浅蓝色背景 */
  border-radius: 8px; /* 圆角 */
  padding: 10px;
  margin: 0 auto; /* 居中 */
  max-width: 100%; /*最大宽度 */
  margin-right: 10px;
}

.message-box p {
  font-size: 16px; /* 默认字体大小 */
  margin: 0;
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 100%; /* 圆形头像 */
  overflow: hidden; /* 隐藏超出部分 */
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover; /* 保持图片比例 */
}
  
  .preview p {
    font-size: 16px; /* 默认字体大小 */
  }
  
  .slider-container {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
  }
  
  input[type="range"] {
    flex: 1;
    margin-right: 10px;
  }
  
  .description {
    margin-bottom: 20px;
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