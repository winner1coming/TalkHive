import axios from 'axios';

// 创建 axios 实例
const apiClient = axios.create({
  baseURL: 'http://your-api-url.com', // 后端 API 的基础 URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// 登录接口
export const login = (username, password) => {
  return apiClient.post('/login', { username, password });
};

// 注册接口
export const register = (username, password) => {
  return apiClient.post('/register', { username, password });
};

// 获取消息列表接口
export const getMessages = () => {
  return apiClient.get('/messages');
};

// 发送消息接口
export const sendMessage = (content) => {
  return apiClient.post('/messages', { content });
};

// 通讯录部分

//
export const getFriendRequests = () =>{
  return apiClient.get('/friendRequests');
} 
export const acceptFriendRequest = (requestId) => {
  return apiClient.post(`/friendRequests/${requestId}/accept`);
};
export const rejectFriendRequest = (requestId) => {
  return apiClient.post(`/friendRequests/${requestId}/reject`);
}

export const getGroupRequests = () =>{
  return apiClient.get('/groupRequests');
}
export const acceptGroupRequest = (requestId) => {
  return apiClient.post(`/groupRequests/${requestId}/accept`);
}
export const rejectGroupRequest = (requestId) => {
  return apiClient.post(`/groupRequests/${requestId}/reject`);
}

// 获取群聊列表接口
export const getGroups = () => {
  return apiClient.get('/groups');
};

// 创建群聊接口
export const createGroup = (name) => {
  return apiClient.post('/groups', { name });
};

// 删除群聊接口
export const deleteGroup = (groupId) => {
  return apiClient.delete(`/groups/${groupId}`);
};

// 获取好友列表接口
export const getFriends = () => {
  return apiClient.get('/friends');
};

// 添加好友接口
export const addFriend = (friendId) => {
  return apiClient.post('/friends', { friendId });
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