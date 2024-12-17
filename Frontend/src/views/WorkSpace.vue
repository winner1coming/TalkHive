<template>
    <div class="workspace">
      <!-- 左侧工具栏 -->
      <aside class="toolbar">
        <ul>
          <li v-for="item in tools" :key="item.name" @click="navigateTo(item.route)" :class="{ active: currentRoute === item.route }">
            {{ item.name }}
          </li>
        </ul>
      </aside>
  
      <!-- 右侧内容区域 -->
      <main class="content">
        <router-view />
      </main>
    </div>
  </template>
  
  <script>
  export default {
    name: "WorkSpace",
    data() {
      return {
        tools: [
          { name: "我的收藏", route: "/workspace/favorites" },
          { name: "我的笔记", route: "/workspace/notes" },
          { name: "我的代码", route: "/workspace/code" },
          { name: "ddl", route: "/workspace/ddl" },
          { name: "回收站", route: "/workspace/recycle" },
        ],
        currentRoute: "",
      };
    },
    watch: {
      $route(to) {
        this.currentRoute = to.path;
      },
    },
    created() {
      this.currentRoute = this.$route.path;
    },
    methods: {
      navigateTo(route) {
        this.$router.push(route);
      },
    },
  };
  </script>
  
  <style scoped>
  .workspace {
    display: flex;
    height: 100vh;
  }
  
  .toolbar {
    width: 200px;
    background-color: #f5f5f5;
    border-right: 1px solid #ddd;
    padding: 10px;
    box-sizing: border-box;
  }
  
  .toolbar ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  
  .toolbar li {
    padding: 10px;
    cursor: pointer;
    transition: background-color 0.3s;
  }
  
  .toolbar li:hover {
    background-color: #eaeaea;
  }
  
  .toolbar li.active {
    font-weight: bold;
    background-color: #d0e8ff;
    color: #007bff;
  }
  
  .content {
    flex: 1;
    padding: 20px;
    box-sizing: border-box;
    overflow-y: auto;
  }
  </style>
  