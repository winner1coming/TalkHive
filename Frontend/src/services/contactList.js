import apiClient from '@/services/api';
export const getFriendRequests = () =>{
    return apiClient.get('/friendRequests');
  } 
  export const acceptFriendRequest = (requestId) => {
    return apiClient.post(`/friendRequests/accept`,{requestId});
  };
  export const rejectFriendRequest = (requestId) => {
    return apiClient.post(`/friendRequests/reject`,{requestId});
  }
  export const getGroupRequests = () =>{
    return apiClient.get('/groupRequests');
  }
  export const acceptGroupInvitationRequest = (accountId,groupId) => {
    return apiClient.post('/groupRequests/acceptInvite', { accountId, groupId });
  }
  export const rejectGroupInvitationRequest = (accountId,groupId) => {
    return apiClient.post('/groupRequests/rejectInvite', { accountId, groupId });
  }
  export const acceptGroupApplyRequest = (accountId,groupId) => {
    return apiClient.post('/groupRequests/acceptApply', { accountId, groupId });
  }
  export const rejectGroupApplyRequest = (accountId,groupId) => {
    return apiClient.post('/groupRequests/rejectApply', { accountId, groupId });
  }
  export const getBlackList = () => {
    return apiClient.get('/blackList');
  }
  export const removeFromBlackList = (accountId) => {
    return apiClient.post('/blackList/remove', { accountId });
  }
  
  // 获取好友列表接口
  export const getFriends = () => {
    return apiClient.get('/friends');
  };
  

  
  
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