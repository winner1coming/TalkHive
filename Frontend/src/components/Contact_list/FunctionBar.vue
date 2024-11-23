<template>
    <div class="function-bar">
      <SearchBar @search="handleSearch" />
      <ul class="menu">
        <li
          v-for="item in menuItems"
          :key="item.name"
          :class="{ active: currentTab === item.name }"
          @click="selectTab(item.name)"
        >
          {{ item.label }}
        </li>
      </ul>
    </div>
  </template>
  
  <script>
  import SearchBar from "../base/SearchBar.vue";
  
  export default {
    props: {
      currentTab: String, // 当前选中的标签
    },
    components: {
      SearchBar,
    },
    data() {
      return {
        menuItems: [
          { name: "friendRequests", label: "好友申请" },
          { name: "groupNotifications", label: "群聊通知" },
          { name: "friendList", label: "好友列表" },
          { name: "groupList", label: "群组列表" },
          { name: "blacklist", label: "黑名单" },
        ],
      };
    },
    methods: {
      selectTab(tabName) {
        this.$emit("tab-selected", tabName); // 触发父组件事件
      },
      handleSearch(query) {
        console.log("搜索内容:", query);
        // 搜索逻辑在父组件实现
      },
    },
  };
  </script>
  
  <style scoped>
  .function-bar {
    width: 250px;
    background-color: #f8f9fa;
    border-right: 1px solid #ddd;
  }
  .menu {
    list-style: none;
    padding: 0;
  }
  .menu li {
    padding: 15px;
    cursor: pointer;
    border-bottom: 1px solid #ddd;
  }
  .menu li:hover {
    background-color: #e9ecef;
  }
  .menu li.active {
    background-color: #007bff;
    color: #fff;
  }
  </style>
  