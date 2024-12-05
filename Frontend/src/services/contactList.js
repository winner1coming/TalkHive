import apiClient from '@/services/api';

// 获取好友请求列表接口
export const getFriendRequests = () => {
  return apiClient.get('/contactList/friendRequests');
};

// 接受好友请求接口
export const acceptFriendRequest = (account_id) => {
  return apiClient.post('/contactList/friendRequests/accept', { account_id });
};

// 拒绝好友请求接口
export const rejectFriendRequest = (account_id) => {
  return apiClient.post('/contactList/friendRequests/reject', { account_id });
};

// 获取群聊请求列表接口
export const getGroupRequests = () => {
  return apiClient.get('/contactList/groupRequests');
};

// 接受群聊邀请请求接口
export const acceptGroupInvitationRequest = (account_id, group_id) => {
  return apiClient.post('/contactList/groupRequests/acceptInvite', { account_id, group_id });
};

// 拒绝群聊邀请请求接口
export const rejectGroupInvitationRequest = (account_id, group_id) => {
  return apiClient.post('/contactList/groupRequests/rejectInvite', { account_id, group_id });
};

// 接受群聊申请请求接口
export const acceptGroupApplyRequest = (account_id, group_id) => {
  return apiClient.post('/contactList/groupRequests/acceptApply', { account_id, group_id });
};

// 拒绝群聊申请请求接口
export const rejectGroupApplyRequest = (account_id, group_id) => {
  return apiClient.post('/contactList/groupRequests/rejectApply', { account_id, group_id });
};

// 获取黑名单列表接口
export const getBlackList = () => {
  return apiClient.get('/contactList/blackList');
};

// 从黑名单中移除用户接口
export const removeFromBlackList = (account_id) => {
  return apiClient.post('/contactList/blackList/remove', { account_id });
};

// 获取好友列表接口
export const getFriends = () => {
  return apiClient.get('/contactList/friends');
};

// 获取群聊列表接口
export const getGroups = () => {
  return apiClient.get('/contactList/groups');
};

// 创建群聊接口
export const createGroup = (name) => {
  return apiClient.post('/contactList/groups', { name });
};

// 删除群聊接口
export const deleteGroup = (group_id) => {
  return apiClient.delete(`/contactList/groups/${group_id}`);
};