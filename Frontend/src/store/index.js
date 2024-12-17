import { createStore } from 'vuex';
import { EventBus } from '@/components/base/EventBus';

export default createStore({
  // 状态对象，包含应用的所有状态
  state: {
    // 用户信息
    user: {
      username: '', // 用户名
      id: '1', // 用户tID 
      avatar:'',
    },
    hasFloatComponent: null,   // 当前正在开启的悬浮组件
    currentChat: null, // 当前聊天对象
    // 系统设置
    settings: {
      theme: '', // 主题颜色
      fontSize: '', // 字体大小
    },
    socket: null,
  },
  
  // 同步修改状态的方法
  mutations: {
    SET_CHAT(state, chat) {
      state.currentChat = chat;
    },

    SET_THEME(state,theme){
      state.settings.theme = theme;
      localStorage.setItem('theme',theme);
    },

    // 设置用户信息
    SET_USER(state, user) {
      state.user = user;
    },

    // 设置消息列表
    SET_MESSAGES(state, messages) {
      state.messages = messages;
    },

    ADD_MESSAGE(state, message) {
      state.messages.push(message);
    },

    // 设置系统设置
    SET_SETTINGS(state, settings) {
      state.settings = settings;
    },

    SET_SOCKET(state, socket) {
      state.socket = socket;
    },
  },
  
  // 异步操作和提交 mutations 的方法
  actions: {
    // 设置聊天对象
    setChat({ commit }, chat) {
      commit('SET_CHAT', chat);
    },

    setTheme({commit},theme){
      commit('SET_THEME',theme);
    },

    // 登录操作
    login({ commit }, user) {
      // 登录逻辑
      commit('SET_USER', user); // 提交 SET_USER mutation
    },

    // 注册操作
    register({ commit }, user) {
      // 注册逻辑
      commit('SET_USER', { username, id: '12345' }); // 提交 SET_USER mutation
    },

    connectWebSocket({ commit, state }) {
      const socket = new WebSocket(`ws://your-websocket-url.com/${state.user.id}`);   // todo
      socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        // 除了对应内容外还需要type字段   todo todo
        if (data.type === 'message') {    // chatlist怎么办
          // 新增contact_id字段，原message内容被封装在message字段中
          if(data.contact_id===state.selectedChatID){
            EventBus.$emit('new-message', data.message);
          }
        } else if (data.type === 'notification') { // todo
          commit('ADD_NOTIFICATION', data);
        }
      };
      socket.onclose = () => {
        console.log('WebSocket connection closed');
      };
      commit('SET_SOCKET', socket);
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

  //获取用户信息
  getters:{
    user:(state) => state.user,
  },

});