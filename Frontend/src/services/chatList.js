import apiClient from '@/services/api';
// 获取聊天列表接口
export const getChatList = () => {
  return apiClient.get('/chatlist');
};
// 获取单个新聊天
export const getChat = (tid) => {  
  // 创建一个新的聊天，后端需要返回新的chat，chat的格式同getChatList中的元素
  return apiClient.post(`/chatlist/createChat`,{tid});
};

// 搜索聊天
export const searchChats = (keyword) => {
  return apiClient.get(`/chatlist/search/${keyword}`);
};

// 置顶或取消置顶聊天
export const pinChat = (tid, is_pinned) => {
  return apiClient.post(`/chatlist/pin`,{tid, is_pinned});
};

// 标记为已读或未读消息
export const readMessages = (tid, is_read) => {
  return apiClient.post(`/messages/read`,{tid, is_read});
}
// 删除聊天
export const deleteChat = (tid) => {
  return apiClient.delete(`/chatlist/${tid}`);
};
// 设置免打扰或取消
export const setMute = (tid, is_muted) => {
  return apiClient.post(`/chatlist/mute`,{tid, is_muted});
};
// 屏蔽聊天或取消
export const blockChat = (tid, is_blocked) => {
  return apiClient.post(`/chatlist/block`,{tid, is_blocked});
};

// 获取聊天消息
export const getMessages = (tid) => {
  return apiClient.get(`/messages/${tid}`);
}
// 发送消息接口
export const sendMessage = (tid, content) => {
  return apiClient.post(`/messages/${tid}/send`, { content });

};
