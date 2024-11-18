<template>
  <!-- 系统设置页面容器 -->
  <div class="system-settings">
    <!-- 页面标题 -->
    <h2>系统设置</h2>
    
    <!-- 主题颜色输入框，使用 v-model 双向绑定到 data 中的 theme 属性 -->
    <input type="text" v-model="theme" placeholder="主题颜色" />
    
    <!-- 字体大小输入框，使用 v-model 双向绑定到 data 中的 fontSize 属性 -->
    <input type="text" v-model="fontSize" placeholder="字体大小" />
    
    <!-- 保存设置按钮，点击时触发 saveSettings 方法 -->
    <button @click="saveSettings">保存设置</button>
  </div>
</template>

<script>
export default {
  // 组件的 data 函数，返回一个对象，包含组件的响应式数据
  data() {
    return {
      // 主题颜色输入框的值
      theme: '',
      // 字体大小输入框的值
      fontSize: '',
    };
  },
  
  // 组件的方法定义
  methods: {
    // 保存设置方法，处理用户点击保存设置按钮时的逻辑
    saveSettings() {
      // 调用 Vuex store 中的 saveSettings 动作，传递主题颜色和字体大小
      this.$store.dispatch('saveSettings', { theme: this.theme, fontSize: this.fontSize });
    },
  },
};
</script>

<style scoped>
/* 系统设置页面的样式 */
.system-settings {
  padding: 20px; /* 设置内边距 */
}
</style>