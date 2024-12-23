<template>
  <!-- 应用的根元素 -->
  <div id="app" :class="themeClass" :style="rootStyle">
      <!-- 路由视图，用于渲染当前路由对应的组件 -->
      <router-view></router-view>
  </div>
</template>

<script>
import { mapActions } from 'vuex';
import { EventBus } from '@/components/base/EventBus';
import { mapGetters } from 'vuex/dist/vuex.cjs.js';
export default {
  // 组件名称
  name: 'App',
  mounted() {
    window.addEventListener('beforeunload', this.saveState);
  },

  beforeDestroy() {
    window.removeEventListener('beforeunload', this.saveState);
  },

  computed:{
    ...mapGetters(['settings']),
    rootStyle(){
      return{
        '--font-size':this.settings.fontSize,
        '--font-family':this.settings.fontStyle,
        '--background-color': this.getBackgroundColor(),
        '--text-color': this.getTextColor(),
        '--button-background-color': this.getButtonBackgroundColor(),
        '--button-color': this.getButtonTextColor(),
      }
    },
    themeClass(){
      return this.settings.theme;
    }
  },

  methods: {
    ...mapActions(['connectWebSocket']),
    
    saveState: function() {
      sessionStorage.setItem("state", JSON.stringify(this.$store.state));
      this.$store.state.user,id;
    },
    hideClick(component=null) {
      if(this.$store.hasFloatCompoent){
        EventBus.emit('hide-float-component', component); // 通知其他组件
        if(component===null){
          this.$store.hasFloatCompoent = false;
        }
      }
    },
    hideContext(event) {
      event.preventDefault();
      const clickedElement = event.target;
      if(this.$store.hasFloatComponent){
        EventBus.emit('close-float-component', clickedElement); // 通知其他组件
      }
    },
    getBackgroundColor() {
      switch (this.settings.theme) {
        case 'light':
          return '#ffffff';
        case 'dark':
          return '#333333';
        case 'system':
          return '#f0f0f0';
        default:
          return '#ffffff';
      }
    },
    getTextColor() {
      switch (this.settings.theme) {
        case 'light':
          return '#000000';
        case 'dark':
          return '#ffffff';
        case 'system':
          return '#000000';
        default:
          return '#000000';
      }
    },
    getButtonBackgroundColor() {
      switch (this.settings.theme) {
        case 'light':
          return '#42b983';
        case 'dark':
          return '#666666';
        case 'system':
          return '#42b983';
        default:
          return '#42b983';
      }
    },
    getButtonTextColor() {
      switch (this.settings.theme) {
        case 'light':
          return '#ffffff';
        case 'dark':
          return '#ffffff';
        case 'system':
          return '#ffffff';
        default:
          return '#ffffff';
      }
    },
  },
  created() {
    //恢复vuex状态
    const savedState = sessionStorage.getItem("state");
    if (savedState) {
      this.$store.replaceState(JSON.parse(savedState));
    }
    // this.connectWebSocket();
    window.addEventListener('click', this.hideClick, true); // 使用 capture 选项
    window.addEventListener('contextmenu', this.hideContext, true); // 使用 capture 选项
    EventBus.on('float-component-open', (component) => {
      if(this.$store.hasFloatComponent){
        EventBus.emit('other-float-component', component); // 通知其他组件
      }
      this.$store.hasFloatComponent = true;
    });
    EventBus.on('hide-float-component', () => {
      this.$store.hasFloatComponent = false;
    });
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
  font-family: var(--font-family);
  font-size:var(--font-size);
  /*font-family: Avenir, Helvetica, Arial, sans-serif;*/

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

.light button{
  background-color: aliceblue;
  color: white;
}

.dark{
  background-color: #333333;
  color: #ffffff;
}

.dark button{
  background-color: #666666;
  color:#ffffff
}

.system{
  background-color: #f0f0f0;
  color:white;
}

.system button{
  background-color: #42b983;
  color:white;
}
</style>
