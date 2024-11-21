import { createStore } from 'vuex';

export default createStore({
  // 状态对象，包含应用的所有状态
  state: {
    // 用户信息
    user: {
      username: '', // 用户名
      id: '', // 用户ID
    },
    // 消息列表
    messages: [],
    // 系统设置
    settings: {
      theme: '', // 主题颜色
      fontSize: '', // 字体大小
    },
  },
  
  // 同步修改状态的方法
  mutations: {
    // 设置用户信息
    SET_USER(state, user) {
      state.user = user;
    },
    // 设置消息列表
    SET_MESSAGES(state, messages) {
      state.messages = messages;
    },
    // 设置系统设置
    SET_SETTINGS(state, settings) {
      state.settings = settings;
    },
  },
  
  // 异步操作和提交 mutations 的方法
  actions: {
    // 登录操作
    login({ commit }, { username, password }) {
      // 登录逻辑
      commit('SET_USER', { username, id: '12345' }); // 提交 SET_USER mutation
    },
    
    // 注册操作
    register({ commit }, { username, password }) {
      // 注册逻辑
      commit('SET_USER', { username, id: '12345' }); // 提交 SET_USER mutation
    },
    
    // 发送消息操作
    sendMessage({ commit, state }, { content }) {
      // 发送消息逻辑
      commit('SET_MESSAGES', [...state.messages, { id: state.messages.length + 1, content }]); // 提交 SET_MESSAGES mutation
    },
    
    // 发送群消息操作
    sendGroupMessage({ commit, state }, { content }) {
      // 发送群消息逻辑
      commit('SET_MESSAGES', [...state.messages, { id: state.messages.length + 1, content }]); // 提交 SET_MESSAGES mutation
    },
    
    // 更新资料操作
    updateProfile({ commit }, { username, id }) {
      // 更新资料逻辑
      commit('SET_USER', { username, id }); // 提交 SET_USER mutation
    },
    
    // 修改密码操作
    changePassword({ commit }, { oldPassword, newPassword }) {
      // 修改密码逻辑
    },
    
    // 保存设置操作
    saveSettings({ commit }, { theme, fontSize }) {
      // 保存设置逻辑
      commit('SET_SETTINGS', { theme, fontSize }); // 提交 SET_SETTINGS mutation
    },
  },
});