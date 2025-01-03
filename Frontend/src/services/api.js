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
    config.headers['User-ID'] = userId; // 在请求头中添加用户 ID
  }
  console.log('请求拦截器:', config); // 打印请求信息
  return config;
}, error => {
  return Promise.reject(error);
});

// 添加响应拦截器
apiClient.interceptors.response.use(response => {
  console.log('响应拦截器:', response); // 打印响应信息
  return response;
}, error => {
  console.error('响应错误拦截器:', error); // 打印响应错误信息
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

  // 添加好友/群聊（id为tid，若为群聊，则为群号）
  export const addStranger = (tid) => {
    return apiClient.post('/stranger/add', { tid });
  };




