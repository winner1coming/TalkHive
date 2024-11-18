// 从 'vue' 模块中导入 Vue 构造函数
import Vue from 'vue';

// 从 './App.vue' 模块中导入 App 组件
import App from './App.vue';

// 从 './router' 模块中导入路由配置
import router from './router';

// 从 './store' 模块中导入 Vuex 状态管理配置
import store from './store';

// 关闭 Vue 的生产提示
Vue.config.productionTip = false;

// 创建一个新的 Vue 实例
new Vue({
  // 注入路由配置
  router,

  // 注入 Vuex 状态管理配置
  store,

  // 使用 render 函数渲染 App 组件
  render: h => h(App),
}).$mount('#app'); // 将 Vue 实例挂载到 id 为 'app' 的 DOM 元素上