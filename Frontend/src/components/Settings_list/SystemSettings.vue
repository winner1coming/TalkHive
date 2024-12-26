<template>
  <div class="system-settings">
    <div class="left-panel">
      <div class="menu-item" :class="{ active: activeComponent === 'ThemeSetting' }" @click="setActiveComponent('ThemeSetting')">
        <span>主题</span>
        <span class="content">{{ theme }}</span>
      </div>
      <div class="menu-item" :class="{ active: activeComponent === 'FontSize' }" @click="setActiveComponent('FontSize')">
        <span>字体大小</span>
        <span class="content">{{ fontsize }}</span>
      </div>
      <div class="menu-item" :class="{ active: activeComponent === 'SoundSetting' }" @click="setActiveComponent('SoundSetting')">
        <span>消息通知</span>
        <span class="content"></span>
      </div>
      <div class="menu-item" :class="{ active: activeComponent === 'ChatBackground' }" @click="setActiveComponent('ChatBackground')">
        <span>聊天背景</span>
        <span class="content"></span>
      </div>
      </div>
    <div class="right-panel">
      <component :is="activeComponent" @updateUser="updateUser" @cancel="setActiveComponent('')"></component>
    </div>
  </div>
</template>

<script>
import ThemeSetting from './ThemeSetting.vue';
import FontSize from './FontSize.vue';
import SoundSetting from './SoundSetting.vue';
import ChatBackground from './ChatBackground.vue';
import { mapGetters } from 'vuex';

export default {
  components: {
   ThemeSetting,
   FontSize,
   SoundSetting,
   ChatBackground,
  },
  data() {
    return {
      theme:'浅色',
      fontsize:'16px',
      activeComponent: '',
    };
  },
  computed:{
    ...mapGetters(['settings']),
  },
  methods: {
    setActiveComponent(component) {
      this.activeComponent = component;
    },
    updateUser(updateUser) {
      // 处理子组件传递的更新数据
      if (updateUser.fontsize) {
        this.fontsize = updateUser.fontsize;
      }
      if (updateUser.theme) {
        this.theme = this.getThemeLabel(updateUser.theme);
      }
      // 关闭当前组件
      this.setActiveComponent('');
    },

    getThemeLabel(themeValue) {
      // 根据 theme 的值返回对应的标签
      const themeMap = {
        light: '浅色',
        dark: '深色',
        system: '系统默认',
      };
      return themeMap[themeValue] || '浅色';
    },
  },
};
</script>

<style scoped>
.system-settings {
  display: flex;
  height: 100vh;
}

.left-panel {
  width: 20%;
  background-color: #f0f0f0;
}

.right-panel {
  width: 80%;
  padding: 20px;
  position: relative;
}

.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #ccc;
  cursor: pointer;
  width: 80%;
  height: 10vh;
  padding: auto;
}

.menu-item.active {
  background-color: #42b983;
  color: white;
}

.menu-item span {
  font-size: 16px;
  margin-left: 40px;
}

.menu-item .content {
  font-size: 14px;
  color: #666;
}

.menu-item.active .content {
  color: white;
}

</style>