import { createStore } from 'vuex';
import { EventBus } from '@/components/base/EventBus';

export default createStore({
  // 状态对象，包含应用的所有状态
  state: sessionStorage.getItem("state") ? JSON.parse(sessionStorage.getItem("state")):{
    // 用户信息
    user: {
      username: 'hh', // 用户名
      id: '1', // 用户tID 
      avatar:'', //默认头像
    },

    hasFloatComponent: null,   // 当前正在开启的悬浮组件
    currentChat: null, // 当前聊天对象
    creatingChat: false, // 是否正在创建聊天
    newChat: null, // 新创建的聊天的参数

    // 系统设置
    settings: {
      theme: 'light', // 主题颜色
      fontSize: '16px', // 字体大小
      fontStyle: 'Microsoft YaHei', // 字体样式
      sound:'dingdo.mp3',
      background:'',
      isNotice:true,
      isNoticeGroup:true,
    },
    links:[],

    socket: null,
    // 要编辑的笔记
    currentNote: {
      note_id: null,
      filename: '',
      category: ''
    },
    // 笔记的标签
    notesCategories: [],
    // 要编辑的代码
    currentCode: {
      code_id: null,
      filename: '',
    },
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
    SET_CREATING_CHAT(state, creatingChat) {
      state.creatingChat = creatingChat;
    },
    SET_NEW_CHAT(state, newChat) {
      state.newChat = newChat;
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

    SET_NOTICE(state, isNotice){
      state.settings.isNotice = isNotice;
    },

    SET_GROUPNOTICE(state,isNoticeGroup){
      state.settings.isNoticeGroup = isNoticeGroup;
    },

    SET_BACKGROUND(state,background){
      state.settings.background = background;
    },

    SET_LINKS(state,links){
      state.links = links;
    },

    SET_SOCKET(state, socket) {
      state.socket = socket;
    },

    //设置当前笔记
    setCurrentNote(state, note) {
      state.currentNote = { ...note };
    },
    setNotesCategories(state, categories) {
      state.notesCategories = categories;
    },
    setCurrentCode(state, code) {
      state.currentCode = { ...code };
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
    setCreatingChat({ commit }, creatingChat) {
      commit('SET_CREATING_CHAT', creatingChat);
    },
    setNewChat({ commit }, newChat) {
      commit('SET_NEW_CHAT', newChat);
    },

    // 注册操作
    register({ commit }, user) {
      // 注册逻辑
      commit('SET_USER', { username, id: '12345' }); // 提交 SET_USER mutation
    },

    connectWebSocket({ commit, state }) {
      const url = `https://localhost:8080/ws/websocketMessage/` + state.user.id.toString();
      const socket = new WebSocket(url);  
      socket.onmessage = (event) => {
        const type = JSON.parse(event.data.type);
        const data = JSON.parse(event.data.data);
        // 播放提示音
        if(state.settings.isNotice){
          const audio = new Audio(require(`@/assets/sound/${state.settings.sound}`));
          audio.play();
        }
        // 除了对应内容外还需要type字段   todo todo
        if (true || type === 'message') {   
          if(data.send_account_id === state.currentChat.id){
            const message ={
              message_id: data.message_id, 
              send_account_id: data.send_id, 
              content: data.content,
              sender: data.sender, 
              create_time: data.create_time, 
              avatar: data.avatar,
              type: data.type, 
            }
            EventBus.emit('new-message', message);
          }
        } else if (type === 'notification') { 
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
      commit('SET_NOTICE',settings.isNotice);
      commit('SET_GROUPNOTICE',settings.isNoticeGroup);
      commit('SET_BACKGROUND',settings.background);
    },

    updateLinks({commit}, links){
      commit('SET_LINKS',links);
    },

    // 笔记
    updateCurrentNote({ commit }, note) {
      commit('setCurrentNote', note);
    },
    updateCategories({ commit }, categories) {
      commit('setNotesCategories', categories);
    },
    updateCurrentCode({ commit }, code) {
      commit('setCurrentCode', code);
    },
  },

  //获取用户信息.计算属性
  getters:{
    user:(state) => state.user,
    settings:(state)=>state.settings,
    links:(state)=>state.links,
    getCurrentNote: (state) => state.currentNote,
    getCategories: (state) => state.notesCategories,
    getCurrentCode: (state) => state.currentCode,
  },

});