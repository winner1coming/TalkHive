<<<<<<< HEAD
<template>
  <div class="system-settings">
    <div class="left-panel">
      <div class="menu-item" :class="{ active: activeComponent === 'ThemeSetting' }" @click="setActiveComponent('ThemeSetting')">
        <img src="@/assets/icon/theme.png" alt="ThemeSetting" class="icon"/>
        <span>主题</span>
        <span class="content">{{ theme }}</span>
      </div>
      <div class="menu-item" :class="{ active: activeComponent === 'FontSize' }" @click="setActiveComponent('FontSize')">
        <img src="@/assets/icon/font.png" alt="FontSize" class="icon"/>
        <span>字体大小</span>
        <span class="content">{{ fontsize }}</span>
      </div>
      <div class="menu-item" :class="{ active: activeComponent === 'SoundSetting' }" @click="setActiveComponent('SoundSetting')">
        <img src="@/assets/icon/notice.png" alt="SoundSetting" class="icon"/>
        <span>消息通知</span>
        <span class="content"></span>
      </div>
      <div class="menu-item" :class="{ active: activeComponent === 'ChatBackground' }" @click="setActiveComponent('ChatBackground')">
        <img src="@/assets/icon/background.png" alt="ChatBackground" class="icon"/>
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
      theme:this.$store.state.settings.theme,
      fontsize:this.$store.state.settings.fontSize,
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
  height: 100%;
}

.left-panel {
  width: 15%;
  background-color: var(--background-color1);
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
  width: 100%;
  height: 10vh;
  padding: 0;
}

.menu-item.active {
  background-color: var(--select-background-color1);
  color: var(--select-text-color);
  font-weight: bold;
}

.menu-item:hover{
  background-color: var(--select-background-color);
}

.menu-item.active:hover{
  background-color: var(--select-background-color1);
}

.menu-item span {
  font-size: var(--font-size);
  margin-left: 10px;
  flex-wrap: nowrap;
  width: fit-content;
}

.menu-item .content {
  font-size: var(--font-size-small);
  color: var(--text-color);
}

.menu-item.active .content {
  color: var(--select-text-color);
}

.menu-item .icon{
  width: 25px;
  height: 25px;
  margin-left: 30px;
}

=======
<template>
  <div class="system-settings">
    <div class="left-panel">
      <div class="menu-item" :class="{ active: activeComponent === 'ThemeSetting' }" @click="setActiveComponent('ThemeSetting')">
        <img src="@/assets/icon/theme.png" alt="ThemeSetting" class="icon"/>
        <span>主题</span>
        <span class="content">{{ theme }}</span>
      </div>
      <div class="menu-item" :class="{ active: activeComponent === 'FontSize' }" @click="setActiveComponent('FontSize')">
        <img src="@/assets/icon/font.png" alt="FontSize" class="icon"/>
        <span>字体大小</span>
        <span class="content">{{ fontsize }}</span>
      </div>
      <div class="menu-item" :class="{ active: activeComponent === 'SoundSetting' }" @click="setActiveComponent('SoundSetting')">
        <img src="@/assets/icon/notice.png" alt="SoundSetting" class="icon"/>
        <span>消息通知</span>
        <span class="content"></span>
      </div>
      <div class="menu-item" :class="{ active: activeComponent === 'ChatBackground' }" @click="setActiveComponent('ChatBackground')">
        <img src="@/assets/icon/background.png" alt="ChatBackground" class="icon"/>
        <span>聊天背景</span>
        <span class="content"></span>
      </div>
    </div>
    <div class="resizer" @mousedown="startResize"></div>
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
      theme:this.$store.state.settings.theme,
      fontsize:this.$store.state.settings.fontSize,
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

    startResize(event) {
      this.isResizing = true;
      document.addEventListener('mousemove', this.resize);
      document.addEventListener('mouseup', this.stopResize);
    },
    resize(event) {
      if (this.isResizing) {
        this.chatListWidth = event.clientX - this.leftComponentWidth;
      }
    },
    stopResize() {
      this.isResizing = false;
      document.removeEventListener('mousemove', this.resize);
      document.removeEventListener('mouseup', this.stopResize);
    },
  },
};
</script>

<style scoped>
.system-settings {
  display: flex;
  height: 100%;
}

.left-panel {
  width: 15%;
  background-color: var(--background-color1);
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
  width: 100%;
  height: 10vh;
  padding: 0;
}

.menu-item.active {
  background-color: var(--select-background-color1);
  color: var(--select-text-color);
  font-weight: bold;
}

.menu-item:hover{
  background-color: var(--select-background-color);
}

.menu-item.active:hover{
  background-color: var(--select-background-color1);
}

.menu-item span {
  font-size: var(--font-size);
  margin-left: 10px;
  flex-wrap: nowrap;
  width: fit-content;
}

.menu-item .content {
  font-size: var(--font-size-small);
  color: var(--text-color);
}

.menu-item.active .content {
  color: var(--select-text-color);
}

.menu-item .icon{
  width: 25px;
  height: 25px;
  margin-left: 30px;
}

.resizer {
  width: 3px;
  height: 100%;
  cursor: ew-resize;
  background-color: #ccc;
}

>>>>>>> back
</style>