import apiClient from '@/services/api';

// 获取好友请求列表接口
export const getFriendRequests = () => {
  return apiClient.get('/contactList/friendRequests');
};

// 处理好友请求接口
export const friendRequestPend = (account_id, accept) => {
  return apiClient.post('/contactList/friendRequests/pend', { account_id, accept});
};

// 获取群聊请求列表接口
export const getGroupRequests = () => {
  return apiClient.get('/contactList/groupRequests');
};

// 处理群聊邀请请求接口
export const groupInvitationRequestPend = (account_id, group_id, accept) => {
  return apiClient.post('/contactList/groupRequests/invitationPend', { account_id, group_id, accept });
};


// 处理群聊申请请求接口
export const groupApplyRequestPend = (account_id, group_id, accept) => {
  return apiClient.post('/contactList/groupRequests/applyPend', { account_id, group_id, accept });
};

// *好友列表
// 获取好友列表接口
export const getFriends = () => {
  return apiClient.get('/contactList/friends');
};
// 获取分组名
export const getDevides = (type) => {   // type: 'friends' or 'groups'
  return apiClient.get(`/contactList/${type}/devides`);
};
// 新建分组
export const createDevide = (type, fd_name) => {   // type: 'friends' or 'groups'
  return apiClient.post(`/contactList/${type}/devides/create`, { fd_name });
};
// 删除分组
export const deleteDevide = (type, fd_name) => {   // type: 'friends' or 'groups'
  return apiClient.delete(`/contactList/${type}/devides/delete/${fd_name}`);
};
// 修改分组名称
export const renameDevide = (type, old_fd_name, new_fd_name) => {   // type: 'friends' or 'groups'
  return apiClient.post(`/contactList/${type}/devides/rename`, {old_fd_name, new_fd_name });
};
// 移动好友到分组
export const moveInDevide = (type, tid, divide) => {
  return apiClient.post(`/contactList/${type}/devides/moveIn`, { tid, divide });
};


// 获取黑名单列表接口
export const getBlackList = () => {
  return apiClient.get('/contactList/blackList');
};

// 从黑名单中移除用户接口
export const removeFromBlackList = (account_id) => {
  return apiClient.post('/contactList/blackList/remove', { account_id });
};
// 拉黑用户
export const addToBlackList = (account_id) => {
  return apiClient.post('/contactList/blackList/add', { account_id });
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
