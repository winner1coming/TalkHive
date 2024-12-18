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
export const sendSmsCode = async (data) => {
  try {
    const response = await apiClient.post('/sendSmsCode',data);
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : error.message;
  }
};

export const smsLogin = async(email) => {
  try{
    const response =await apiClient.post('/smslogin',email);
      return response.data;
  }catch(error){
      throw error.response?.data || error.message;
    };
};

// 注册接口
export const Register =async (data) => {
  try{
    const response = await apiClient.post('/register', data);
    return response.data;
  } 
  catch(error) {
    throw error.response?.data || error.message;
  };
};


// 重置密码接口
export const resetPassword =async (msg) => {
  try{
    const response = await apiClient.post('/resetPassword', msg);
    return response.data;
  } 
  catch(error) {
    throw error.response?.data || error.message;
  };
};

// 获取用户信息接口
export const showProfile = async (id) => {
  try {
    const response = await apiClient.get(`/Settings/profile/${id}`);
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

// 更新用户信息接口
export const saveEdit = async (data) => {
  try {
    const response = await apiClient.post('/Settings/saveEdit', data);
    return response.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
};

//获取用户的邮箱（安全设置）
export const getUserInfo = async(id)=>{
  try{
    const response = await apiClient.get('/Settings/getInfo',{id});
    return response.data;

  }catch(error){
    throw error.response?.data ||error.message;
  }
};

//更换邮箱时获取验证码
export const getCode = async(data)=>{
  try{
    const response = await apiClient.post('/Settings/getCode',data);
    return response.data;
  }catch(error){
    throw error.response?.data||error.message;
  }
};

  //注销账号
  export const confirmDeactivation = async(id)=>{
    try{
      const response = await apiClient.post('/Settings/deactivate',{id});
      return response.data;
    }catch(error){
      throw error.response?.data ||error.message;
    }
  };

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
