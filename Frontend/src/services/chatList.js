import apiClient from '@/services/api';
// 获取聊天列表接口
export const getChatList = () => {
  return apiClient.get('/chatlist');
};
// 获取单个新聊天
export const getChat = (tid, is_group) => {  
  // 创建一个新的聊天，后端需要返回新的chat，chat的格式同getChatList中的元素
  return apiClient.post(`/chatlist/createChat`,{tid, is_group});
};
// 搜索聊天
export const searchChats = (keyword) => {
  return apiClient.get(`/chatlist/search/${keyword}`);
};

// 置顶或取消置顶聊天
export const pinChat = (tid, is_pinned, is_group) => {
  return apiClient.post(`/chatlist/pin`,{tid, is_pinned, is_group});
};

// 标记为已读或未读消息
export const readMessages = (tid, is_read, is_group) => {
  return apiClient.post(`/messages/read`,{tid, is_read, is_group});
}
// 删除聊天
export const deleteChat = (tid, is_group) => {
  return apiClient.post(`/chatlist/delete`,{tid, is_group});
};
// 设置免打扰或取消
export const setMute = (tid, is_muted, is_group) => {
  return apiClient.post(`/chatlist/mute`,{tid, is_muted, is_group});
};
// 屏蔽聊天或取消
export const blockChat = (tid, is_blocked, is_group) => {
  return apiClient.post(`/chatlist/block`,{tid, is_blocked, is_group});
};

// 获取聊天消息
export const getMessages = (tid, is_group) => {
  return apiClient.post(`/messages`, {tid, is_group});
}
// 发送消息接口
export const sendMessage = (tid, content, type, is_group) => {
  return apiClient.post(`/messages/send`, { tid, content, type, is_group });
};
// 收藏消息
export const collectMessage = (table_name, message_id) => {
  return apiClient.post(`/messages/collect`, { table_name, message_id });
};
// 删除消息
export const deleteMessage = (message_id) => {
  return apiClient.post(`/messages/delete`, { message_id });
};
// 获取历史记录
export const getHistory = (tid) => {
  return apiClient.post(`/messages/history`,{tid});
};
