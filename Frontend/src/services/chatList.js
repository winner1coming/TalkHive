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
export const pinChat = (friendId) => {
  return apiClient.post(`/chatlist/pin`,{friendId});
};
// 取消置顶聊天
export const unpinChat = (friendId) => {
  return apiClient.post(`/chatlist/unpin`,{friendId});
};
// 标记为已读消息
export const readMessages = (friendId) => {
  return apiClient.post(`/messages/read`,{friendId});
}
// 标记为未读消息
export const unreadMessages = (friendId) => {
  return apiClient.post(`/messages/unread`,{friendId});
}
// 删除聊天
export const deleteChat = (friendId) => {
  return apiClient.delete(`/chatlist/${friendId}`);
};
// 设置免打扰
export const setMute = (friendId) => {
  return apiClient.post(`/chatlist/mute`,{friendId});
};
// 取消免打扰
export const cancelMute = (friendId) => {
  return apiClient.post(`/chatlist/unmute`,{friendId});
};
// 屏蔽聊天
export const blockChat = (friendId) => {
  return apiClient.post(`/chatlist/block`,{friendId});
};
// 取消屏蔽聊天
export const unblockChat = (friendId) => {
  return apiClient.post(`/chatlist/unblock`,{friendId});
};


// 获取聊天消息
export const getMessages = (friendId) => {
  return apiClient.get(`/messages/${friendId}`);
}
// 发送消息接口
export const sendMessage = (content) => {
  return apiClient.post('/messages', { content });
};
