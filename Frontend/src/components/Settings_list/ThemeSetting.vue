<template>
    <div class="theme-settings">
      <h3>主题</h3>
      <div class="theme-options">
        <div v-for="theme in themes" :key="theme.value" class="theme-option">
          <label :for="theme.value">
            <input
              type="radio"
              :id="theme.value"
              :value="theme.value"
              v-model="selectedTheme"
              @change="onThemeChange"
            />
            <span class="checkmark"></span>
            {{ theme.label }}
          </label>
        </div>
      </div>
    </div>
  </template>
  
  <script>
import { mapState,mapActions } from 'vuex';

  export default {
    data() {
      return {
        themes: [
          { label: '浅色模式', value: 'light' },
          { label: '深色模式', value: 'dark' },
          { label: '系统默认', value: 'system' },
        ],
      };
    },
    computed:{
        ...mapState(['setting']),
        selectedTheme:{
            get(){
                return this.setting?.theme;
            },
            set(value){
                this.setTheme(value);
            },
        },
    },
    watch: {
    selectedTheme(newVal) {
      // 手动更新视图
      this.$forceUpdate();
    },
  },
    methods: {
      ...mapActions(['setTheme']),
      onThemeChange() {
        // 在这里处理主题切换逻辑
        console.log('Selected theme:', this.selectedTheme);
        // 你可以在这里调用一个方法来应用选中的主题
      },
    },
  };
  </script>
  
  <style scoped>
  .theme-settings {
    padding: 20px;
  }
  
  h3 {
    margin-bottom: 20px;
  }
  
  .theme-options {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .theme-option {
    margin-bottom: 10px;
  }
  
  label {
    display: flex;
    align-items: center;
    cursor: pointer;
    font-size: 16px;
  }
  
  input[type="radio"] {
    position: absolute;
    opacity: 0;
    cursor: pointer;
  }
  
  .checkmark {
    position: relative;
    display: inline-block;
    width: 20px;
    height: 20px;
    margin-right: 10px;
    border: 2px solid #ccc;
    border-radius: 50%;
  }
  
  input[type="radio"]:checked ~ .checkmark {
    border-color: #42b983;
  }
  
  input[type="radio"]:checked ~ .checkmark::after {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 10px;
    height: 10px;
    background-color: #42b983;
    border-radius: 50%;
    transform: translate(-50%, -50%);
  }
  </style>