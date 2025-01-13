<template>
  <!-- 应用的根元素 -->
  <div id="app" :class="themeClass" :style="rootStyle">
      <Notification ref="notification" />
      <!-- 路由视图，用于渲染当前路由对应的组件 -->
      <router-view></router-view>
  </div>
</template>

<script>
import { mapActions } from 'vuex';
import { EventBus } from '@/components/base/EventBus';
import Notification from '@/components/base/Notification.vue';
import { mapGetters } from 'vuex/dist/vuex.cjs.js';

export default {
  name: 'App',

  components: {
    Notification,
  },

  mounted() {
    window.addEventListener('beforeunload', this.saveState);
  },

  beforeDestroy() {
    window.removeEventListener('beforeunload', this.saveState);
  },

  computed:{
    ...mapGetters(['settings']),
    rootStyle(){
      const baseFontSize = parseInt(this.settings.fontSize, 10);
      return{
        '--font-size':this.settings.fontSize,
        '--font-size-small-small':`${baseFontSize - 3}px`,
        '--font-size-small':`${baseFontSize - 2}px`,
        '--font-size-mlarge':`${baseFontSize + 2}px`,
        '--font-size-large': `${baseFontSize + 4}px`,
        '--font-family':this.settings.fontStyle,
        '--background-color': this.getBackgroundColor(),   //从浅到深 白
        '--background-color1': this.getBackgroundColor1(), 
        '--background-color2': this.getBackgroundColor2(),
        '--text-color': this.getTextColor(),
        '--button-background-color': this.getButtonBackgroundColor(),   //普通状态  light 绿色
        '--button-background-color1': this.getButtonBackgroundColor1(), //悬浮
        '--button-background-color2': this.getButtonBackgroundColor2(),  //点击
        '--button-text-color': this.getButtonTextColor(),
        '--sidebar-background-color': this.getSidebarBackgroundColor(),   //普通状态栏色调（浅  light紫
        '--sidebar-background-color1': this.getSidebarBackgroundColor1(), //鼠标悬浮状态栏色调
        '--sidebar-background-color2': this.getSidebarBackgroundColor2(), //鼠标点击状态栏色调
        '--sidebar-text-color': this.getSidebarTextColor(),
        '--select-background-color' : this.selectBackgroundColor(),   //新增可选择的颜色色调（light 蓝色——自行搭配一下，用不上也没关系
        '--select-background-color1' : this.selectBackgroundColor1(),
        '--select-background-color2' : this.selectBackgroundColor2(),
        '--select-text-color' : this.selectTextColor(),
      }
    },
    themeClass(){
      return this.settings.theme;
    },
  },

  methods: {
    ...mapActions(['connectWebSocket']),
    
    saveState: function() {
      this.$store.dispatch('setChat',null);
      sessionStorage.setItem("state", JSON.stringify(this.$store.state));
      this.$store.state.user,id;
    },

    getBackgroundColor() {
      switch (this.settings.theme) {
        case 'light':
          return '#ffffff';   //纯白色
        case 'dark':
          return '#333333';
        case 'system':
          return '#f0f0f0';
        default:
          return '#ffffff';
      }
    },

    getBackgroundColor1() {
      switch (this.settings.theme) {
        case 'light':
          return '#eaeaea';   //更深一点的灰色
        case 'dark':
          return '#333333';
        case 'system':
          return '#f0f0f0';
        default:
          return '#ffffff';
      }
    },

    getBackgroundColor2() {
      switch (this.settings.theme) {
        case 'light':
          return '#d5d2d2';   //最深的灰色
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
          return '#000000';   //黑色字体
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
          return '#8ae2ba';     //浅绿色
        case 'dark':
          return '#666666';
        case 'system':
          return '#42b983';
        default:
          return '#42b983';
      }
    },
    
    getButtonBackgroundColor1() {
      switch (this.settings.theme) {
        case 'light':
          return '#6bc76e';   //更深的绿色
        case 'dark':
          return '#666666';
        case 'system':
          return '#42b983';
        default:
          return '#42b983';
      }
    },
    
    getButtonBackgroundColor2() {
      switch (this.settings.theme) {
        case 'light':
          return '#2b9d19';   //最深的绿色
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
          return '#333333';   //黑色字体
        case 'dark':
          return '#ffffff';
        case 'system':
          return '#ffffff';
        default:
          return '#ffffff';
      }
    },

    getSidebarBackgroundColor() {
      switch (this.settings.theme) {
        case 'light':
          return '#af8edd95'; // 浅色主题的侧边栏背景颜色  浅紫色
        case 'dark':
          return '#2c3e50'; // 深色主题的侧边栏背景颜色
        case 'system':
          return '#6dc79fb1'; // 根据系统主题设置
        default:
          return '#6dc79fb1';
      }
    },

    getSidebarBackgroundColor1() {
      switch (this.settings.theme) {
        case 'light':
          return '#7445b666'; // 更深的紫色
        case 'dark':
          return '#2c3e50'; // 深色主题的侧边栏背景颜色
        case 'system':
          return '#6dc79fb1'; // 根据系统主题设置
        default:
          return '#6dc79fb1';
      }
    },

    getSidebarBackgroundColor2() {
      switch (this.settings.theme) {
        case 'light':
          return '#8757c9a1'; // 更深的紫色
        case 'dark':
          return '#2c3e50'; // 深色主题的侧边栏背景颜色
        case 'system':
          return '#6dc79fb1'; // 根据系统主题设置
        default:
          return '#6dc79fb1';
      }
    },

    getSidebarTextColor() {
      switch (this.settings.theme) {
        case 'light':
          return '#000000'; // 黑色
        case 'dark':
          return '#ffffff'; // 深色主题的侧边栏文字颜色
        case 'system':
          return '#000000'; // 根据系统主题设置
        default:
          return '#000000';
      }
    },

    selectBackgroundColor() {
      switch (this.settings.theme) {
        case 'light':
          return '#bcdfff88'; // 蓝色设置
        case 'dark':
          return '#2c3e50'; // 深色主题的侧边栏背景颜色
        case 'system':
          return '#6dc79fb1'; // 根据系统主题设置
        default:
          return '#6dc79fb1';
      }
    },

    selectBackgroundColor1() {
      switch (this.settings.theme) {
        case 'light':
          return '#d0e8ff'; // 更深一点的蓝色
        case 'dark':
          return '#2c3e50'; // 深色主题的侧边栏背景颜色
        case 'system':
          return '#6dc79fb1'; // 根据系统主题设置
        default:
          return '#6dc79fb1';
      }
    },

    selectBackgroundColor2() {
      switch (this.settings.theme) {
        case 'light':
          return 'rgba(26, 216, 226, 0.879)'; // 最深的蓝色
        case 'dark':
          return '#2c3e50'; // 深色主题的侧边栏背景颜色
        case 'system':
          return '#6dc79fb1'; // 根据系统主题设置
        default:
          return '#6dc79fb1';
      }
    },

    selectTextColor() {
      switch (this.settings.theme) {
        case 'light':
          return '#007bff'; // 蓝色
        case 'dark':
          return '#ffffff'; // 深色主题的侧边栏文字颜色
        case 'system':
          return '#000000'; // 根据系统主题设置
        default:
          return '#000000';
      }
    },



    // 通知
    notify(message, type) {
      this.$refs.notification.show(message, type);
    },

    // 点击事件
    hideClick(event) {
      const clickedElement = event.target;
      if(this.$store.hasFloatComponent){
        EventBus.emit('close-float-component', clickedElement); // 通知其他组件
      }
    },

    hideContext(event) {
      event.preventDefault();
      const clickedElement = event.target;
      if(this.$store.hasFloatComponent){
        EventBus.emit('close-float-component', clickedElement); // 通知其他组件
      }
    },

  },
  created() {
    //恢复vuex状态
    const savedState = sessionStorage.getItem("state");
    if (savedState) {
      this.$store.replaceState(JSON.parse(savedState));
    }

    // 全局监视器
    //this.$store.dispatch('connectWebSocket');
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

button{
  font-family: var(--font-family);
  font-size:var(--font-size);
}

</style>
