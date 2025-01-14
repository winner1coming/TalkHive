<template>
    <div class="search-bar">
      <img src="@/assets/icon/find.png" alt="find" class="find-icon"/>
        <input
        type="text"
        v-model="query"
        :placeholder="placeholder"
        @input="handleInput"
        @focus="handleFocus"
        @blur="handleBlur"
      />
    </div>
  </template>
  
  <script>
  export default {
    props: {
      placeholder: {
        type: String,
        default: "搜索...",
      },
    },
    data() {
      return {
        query: "", // 搜索关键词
        isFocused: false, // 是否聚焦
      };
    },
    methods: {
      // 处理输入事件
      handleInput() {
        this.$emit("input", this.query); // 向父组件发送输入事件
      },
      // 处理聚焦事件
      handleFocus() {
        this.isFocused = true;
        this.$emit("focus"); // 向父组件发送聚焦事件
      },
      // 处理失焦事件
      handleBlur() {
        this.isFocused = false;
        this.$emit("blur"); // 向父组件发送失焦事件
      },
    },
  };
  </script>
  
  <style scoped>
  .search-bar {
    display: flex;
    padding: 10px;
    margin-top: 10px;
  }
  
  .search-bar input {
    flex:1;
    padding: 5px;
    border-radius: 4px;
    padding-left: 30px;
    color: var(--text-color);
    background-color: var(--background-color);
  }

  .find-icon {
    width: 25px;
    height: 25px;
    position: absolute;
    top: 24px;
    left: 32px;
    font-size: var(--font-size);
    color: #666;
    cursor: pointer;
  }
  </style>