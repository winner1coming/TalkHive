import axios from 'axios';
import store from '@/store'; 
// 创建 axios 实例
const apiClient = axios.create({
  baseURL: 'http://localhost:8080', // 后端 API 的基础 URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// 添加请求拦截器（请求头中有id，后端可通过headers['tid']来获取id）
apiClient.interceptors.request.use(config => {
  const userId = store.state.user.id; // 从 Vuex 存储中获取用户 ID
  if (userId) {
    config.headers['tid'] = userId; // 在请求头中添加用户 ID
  }
  return config;
}, error => {
  return Promise.reject(error);
});

export default apiClient;


  // 获取消息列表接口
  export const getMessages = () => {
    return apiClient.get('/messages');
  };


  // chat和contact
  // 资料卡片
  export const getProfileCard = (tid, group_id=null) => {
    return apiClient.get(`/profileCard/${tid}`, {group_id});
  };
  // 搜索好友/群聊（key可能是id或者昵称）
  export const searchFriendGroup = (key) => {
    return apiClient.get('/Stranger/search', { key });  
  };
  // 添加好友/群聊（id为tid，若为群聊，则为群号）
  export const addFriendGroup = (id) => {
    return apiClient.post('/Stranger/add', { id });
  };
  // 新建群聊(tids为成员id列表，其中没有用户自己的)
  export const createGroup = (tids) => {
    return apiClient.post('/GroupList/create', { tids });
  };


  // 获取笔记列表接口
  export const getNotes = () => {
    return apiClient.get('/notes');
  };

  // 创建笔记接口
  export const createNote = (title, content) => {
    return apiClient.post('/notes', { title, content });
  };

  // 获取收藏列表接口
  export const getFavorites = () => {
    return apiClient.get('/favorites');
  };

  // 添加收藏接口
  export const addFavorite = (itemId) => {
    return apiClient.post('/favorites', { itemId });
  };

  // 获取 DDL 列表接口
  export const getDDLs = () => {
    return apiClient.get('/ddls');
  };

  // 创建 DDL 接口
  export const createDDL = (title, deadline) => {
    return apiClient.post('/ddls', { title, deadline });
  };
