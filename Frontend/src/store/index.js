import { createStore } from 'vuex';
import { EventBus } from '@/components/base/EventBus';

export default createStore({
  // 状态对象，包含应用的所有状态
  state: sessionStorage.getItem("state") ? JSON.parse(sessionStorage.getItem("state")):{
    // 用户信息
    user: {
      username: 'hh', // 用户名
      id: '1', // 用户tID 
      avatar:'',
    },

    hasFloatComponent: null,   // 当前正在开启的悬浮组件
    currentChat: null, // 当前聊天对象

    // 系统设置
    settings: {
      theme: 'light', // 主题颜色
      fontSize: '16px', // 字体大小
      fontStyle: 'Arial,sans-serif', // 字体样式
      sound:'',
    },
    socket: null,
  },
  
  // 同步修改状态的方法
  mutations: {
    // 设置用户信息
    SET_USER_ID(state, id) {
      state.user.id = id;
    },

    SET_USER_NAME(state,username){
      state.user.username = username;
    },

    SET_USER_AVATAR(state,avatar){
      state.user.avatar = avatar;
    },

    //设置整个用户——用于初始化和整体更新
    SET_USER(state,user){
      state.user = {...state.user,...user};
    },
    
    SET_CHAT(state, chat) {
      state.currentChat = chat;
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
      state.settings = {...state.settings,...settings};
    },

    SET_THEME(state, theme) {
      state.settings.theme = theme;
    },

    SET_FONTSIZE(state,fontSize){
      state.settings.fontSize = fontSize;
    },

    SET_FONTSTYLE(state,fontStyle){
      state.settings.fontStyle = fontStyle;
    },

    SET_SOUND(state,sound){
      state.settings.sound = sound;
    },

    SET_SOCKET(state, socket) {
      state.socket = socket;
    },
  },
  
  // 异步操作和提交 mutations 的方法
  actions: {
    // 登录操作
    setUser({ commit }, user) {
      // 登录逻辑
      commit('SET_USER_ID', user.id); // 提交 SET_USER mutation
      commit('SET_USER_NAME',user.username);
      commit('SET_USER_AVATAR',user.avatar);
    },

    // 设置聊天对象
    setChat({ commit }, chat) {
      commit('SET_CHAT', chat);
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
        
    // 保存设置操作
    saveSettings({ commit }, settings) {
      // 保存设置逻辑
      commit('SET_THEME', settings.theme); // 提交 SET_SETTINGS mutation
      commit('SET_FONTSIZE', settings.fontSize);
      commit('SET_FONTSTYLE',settings.fontStyle);
      commit('SET_SOUND',settings.sound);
    },
  },

  //获取用户信息.计算属性
  getters:{
    user:(state) => state.user,
    settings:(state)=>state.settings,
  },

});