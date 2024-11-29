<template>
  <!-- 应用的根元素 -->
  <div id="app">
      <!-- 路由视图，用于渲染当前路由对应的组件 -->
      <router-view></router-view>
  </div>
</template>

<script>
import { mapActions } from 'vuex';
export default {
  // 组件名称
  name: 'App',
  created() {
    // this.connectWebSocket();
  },
  methods: {
    ...mapActions(['connectWebSocket']),
  },
}

</script>

<style>
/* 应用的全局样式 */
#app {
  /* 设置字体 */
  font-family: Avenir, Helvetica, Arial, sans-serif;

  /* 抗锯齿，使字体在 Webkit 浏览器中更平滑 */
  -webkit-font-smoothing: antialiased;

  /* 抗锯齿，使字体在 Firefox 浏览器中更平滑 */
  -moz-osx-font-smoothing: grayscale;

  /* 文本居中对齐 */
  text-align: center;

  /* 设置文本颜色 */
  color: #2c3e50;

  position: absolute;
  width: 100%;
  height: 100%;
}
* {
  margin: 0;
  padding: 0;
}
</style>
