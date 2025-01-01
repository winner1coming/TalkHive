<template>
    <div v-if="isImmidiate" class="search-bar" >
      <input
        type="text"
        v-model="query"
        placeholder="搜索..."
        @compositionstart="isComposing = true"
        @compositionend="isComposing = false;triggerSearch()"
        @input="triggerSearch"
      />
      <button @click="buttonClick">+</button>
    </div>
    <div v-else class="search-bar">
      <input
        type="text"
        v-model="query"
        placeholder="搜索..."
        @compositionstart="isComposing = true"
        @compositionend="isComposing = false;triggerSearch()"
        @keydown.enter="triggerSearch"
      />
      <button v-show="showButton" @click="triggerSearch">搜索</button>
    </div>

  </template>
  
  <script>
  export default {
    props:{
      isImmidiate:{
        type:Boolean,
        default:true
      },
      showButton:{
        type:Boolean,
        default:true,
      }
    },
    data() {
      return {
        query: "", // 搜索关键词
        isComposing: false, // 是否正在使用输入法输入，防止频繁触发搜索
      };
    },
    methods: {
      triggerSearch() {
        if (this.isComposing) return; // 正在输入中，不触发搜索
        console.log("searching...");
        this.$emit("search", this.query); // 向父组件发送搜索事件
      },
      buttonClick(event){
        this.$emit("button-click", event);
      },
      clear(){
        this.query = "";
      }
    },
  };
  </script>
  
  <style scoped>
  .search-bar {
    display: flex;
    padding: 10px;
  }
  .search-bar input {
    flex: 1;
    padding: 5px;
    border: 1px solid #ddd;
    border-radius: 4px;
  }
  .search-bar button {
    margin-left: 5px;
    padding: 5px 10px;
    border: none;
    background-color: #007bff;
    color: #fff;
    border-radius: 4px;
    cursor: pointer;
  }
  </style>
  