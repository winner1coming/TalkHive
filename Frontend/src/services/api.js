import axios from 'axios';
import store from '@/store'; 
// 创建 axios 实例
const apiClient = axios.create({
  baseURL: 'http://your-api-url.com', // 后端 API 的基础 URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// 添加请求拦截器（请求头中有id，后端可通过headers['User-ID']来获取id）
apiClient.interceptors.request.use(config => {
  const userId = store.state.user.id; // 从 Vuex 存储中获取用户 ID
  if (userId) {
    config.headers['User-ID'] = userId; // 在请求头中添加用户 ID
  }
  return config;
}, error => {
  return Promise.reject(error);
});

export default apiClient;

// 登录接口
export const login = async (payload) => {
  try {
    const response = await apiClient.post('/login', payload);
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : error.message;
  }
};
//短信接收码的接口
export const sendSmsCode = async (phoneNumber) => {
  try {
    const response = await axios.post('/sendSmsCode', {
      command: 'SendSmsCode',
      phoneNumber,
    });
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : error.message;
  }
};

export const sendVerificationCode = (phoneNumber) => {
  return apiClient.post('/sendVerificationCode', { phoneNumber })
    .then(response => {
      return response.data;
    })
    .catch(error => {
      throw error.response?.data || error.message;
    });
};

// 注册接口
export const register = (username, password) => {
  return apiClient.post('/register', { username, password });
};


// 重置密码接口
export const resetPassword = (data) => {
  return apiClient.post('/resetPassword', data)
    .then(response => {
      return response.data;
    })
    .catch(error => {
      throw error.response?.data || error.message;
    });
};

// 获取消息列表接口
export const getMessages = () => {
  return apiClient.get('/messages');

// chat和contact
// 搜索好友/群聊（key可能是id或者昵称）
export const searchFriendGroup = (key) => {
  return apiClient.get('/search/Stranger', { key });  
};
// 添加好友/群聊
export const addFriendGroup = (id) => {
  return apiClient.post('/add/Stranger', { id });
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
