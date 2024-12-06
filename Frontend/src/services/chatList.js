import apiClient from '@/services/api';
// 获取聊天列表接口
export const getChatList = () => {
  return apiClient.get('/chatlist');
};
// 搜索聊天
export const searchChats = (keyword) => {
  return apiClient.get(`/chatlist/search/${keyword}`);
};

// 置顶聊天
export const pinChat = (tid) => {
  return apiClient.post(`/chatlist/pin`,{tid});
};
// 取消置顶聊天
export const unpinChat = (tid) => {
  return apiClient.post(`/chatlist/unpin`,{tid});
};
// 标记为已读消息
export const readMessages = (tid) => {
  return apiClient.post(`/messages/read`,{tid});
}
// 标记为未读消息
export const unreadMessages = (tid) => {
  return apiClient.post(`/messages/unread`,{tid});
}
// 删除聊天
export const deleteChat = (tid) => {
  return apiClient.delete(`/chatlist/${tid}`);
};
// 设置免打扰
export const setMute = (tid) => {
  return apiClient.post(`/chatlist/mute`,{tid});
};
// 取消免打扰
export const cancelMute = (tid) => {
  return apiClient.post(`/chatlist/unmute`,{tid});
};
// 屏蔽聊天
export const blockChat = (tid) => {
  return apiClient.post(`/chatlist/block`,{tid});
};
// 取消屏蔽聊天
export const unblockChat = (tid) => {
  return apiClient.post(`/chatlist/unblock`,{tid});
};

// 产生新聊天
export const generateNewChat = (tid) => {  
  // 创建一个新的聊天，后端需要返回新的chat，chat的格式同getChatList中的元素
  return apiClient.post(`/chatlist/createChat`,{tid});
};

// 获取聊天消息
export const getMessages = (tid) => {
  return apiClient.get(`/messages/${tid}`);
}
// 发送消息接口
export const sendMessage = (tid, content) => {
  return apiClient.post(`/messages/${tid}/send`, { content });
};
