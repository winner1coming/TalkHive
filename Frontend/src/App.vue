<template>
  <!-- 应用的根元素 -->
  <div id="app">
      <!-- 路由视图，用于渲染当前路由对应的组件 -->
      <router-view></router-view>
  </div>
</template>

<script>
import { mapActions } from 'vuex';
import { EventBus } from '@/components/base/EventBus';
export default {
  // 组件名称
  name: 'App',
  
  methods: {
    ...mapActions(['connectWebSocket']),
    hideClick() {
      EventBus.emit('hide-float-component'); // 通知其他组件
    },
    hideContext(event) {
      event.preventDefault();
      EventBus.emit('hide-float-component'); // 通知其他组件
    },
  },
  created() {
    // this.connectWebSocket();
    window.addEventListener('click', this.hideClick, true); // 使用 capture 选项
    window.addEventListener('contextmenu', this.hideContext, true); // 使用 capture 选项
  },
  beforeUnmount() {
    window.removeEventListener('click', this.hideClick, true); 
    window.removeEventListener('contextmenu', this.hideContext, true); 
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
