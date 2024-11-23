<template>
    <div class="search-bar">
      <input
        type="text"
        v-model="query"
        placeholder="搜索..."
        @keyup.enter="triggerSearch"
      />
      <button @click="triggerSearch">搜索</button>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        query: "", // 搜索关键词
      };
    },
    methods: {
      triggerSearch() {
        this.$emit("search", this.query); // 向父组件发送搜索事件
      },
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
  