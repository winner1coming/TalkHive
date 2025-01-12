import Mock from 'mockjs';
const baseURL = 'http://localhost:8080';

Mock.setup({
  timeout: '200-300', // 设置模拟延迟（可选）
});
const addCorsHeaders = (response) => {
  return {
    ...response,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization'
    }
  };
};

// 搜索陌生人
Mock.mock(`${baseURL}/stranger/search`, 'post', (options) => {
  const { key } = JSON.parse(options.body);
  const results = [];

  friends.friends.forEach(friend => {
    if (friend.account_id.includes(key) || friend.remark.includes(key)) {
      results.push({
        tid: friend.account_id,
        id: friend.account_id,
        nickname: friend.remark,
        avatar: friend.avatar,
        type: 'friend'
      });
    }
  });

  groups.groups.forEach(group => {
    if (group.account_id.includes(key) || group.remark.includes(key)) {
      results.push({
        tid: group.account_id,
        id: group.account_id,
        nickname: group.remark,
        avatar: group.avatar,
        type: 'group'
      });
    }
  });

  return {
    status: 200,
    data: results,
  };
});
// 搜索通讯录
Mock.mock(`${baseURL}/contactList/search`, 'post', (options) => {
  const { keyword } = JSON.parse(options.body);
  let results = [];
  friends.friends.forEach(friend => {
    if (friend.account_id.includes(keyword) || friend.remark.includes(keyword)) {
      results.push({
        type: 'friends',
        account_id: friend.account_id,
        avatar: friend.avatar,
        remark: friend.remark,
        status: friend.status,
        signature: friend.signature,
      });
    }
  });
  groups.groups.forEach(group => {
    if (group.account_id.includes(keyword) || group.remark.includes(keyword)) {
      results.push({
        type: 'groups',
        account_id: group.account_id,
        avatar: group.avatar,
        remark: group.remark,
        status: null,
        signature: group.signature,
      });
    }
  });
  return {
    status: 200,
    data: results,
  };

});
let friendRequests = Mock.mock({
  'requests|5-10': [{
    'apply_id|+1': 1,
    'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
    'name': '@name',
    'sender_id|1': /[0-9]{10}/,
    'receiver_id|1': /[0-9]{10}/,
    'reason': () => Mock.Random.csentence(3, 20),
    'status': '@pick(["pending", "accepted", "rejected"])',
    'time': '@datetime',
  }]
});

let groupRequests = Mock.mock({
  'requests|5-10': [{
    'apply_id|+1': 1,
    'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
    'group_name': '@name',
    'account_name': '@name',
    'sender_id|1-100': 1,
    'receiver_id|1-100': 1,
    'group_id|1-100': 1,
    'reason': () => Mock.Random.csentence(3, 20),
    'type': '@pick(["groupInvitation", "groupApply", "notification"])',
    'status': '@pick(["pending", "accepted", "rejected", "waiting", "notification"])',
    'time': '@datetime',
  }]
});

// 获取好友请求列表
Mock.mock(`${baseURL}/contactList/friendRequests`, 'get', () => {
  return { data: friendRequests.requests };
});
// 处理好友请求
Mock.mock(`${baseURL}/contactList/friendRequests/pend`, 'post', (options) => {
  const { account_id, accept } = JSON.parse(options.body);
  const request = friendRequests.requests.find(req => req.sender_id === account_id);
  if (request) {
    request.status = accept ? 'accepted' : 'rejected';
  }
  return {
    status: 200,
    data: request,
  };
});
// 添加好友
Mock.mock(`${baseURL}/contactList/friendRequests/addFriend`, 'post', (options) => {
  const { account_id, reason } = JSON.parse(options.body);
  const newFriend = Mock.mock({
    'friend': {
      'apply_id': friendRequests.requests.length + 1,
      'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
      'name': '@name',
      'sender_id': account_id,
      'receiver_id': '1000000000',
      'reason': reason,
      'status': 'pending',
      'time': '@datetime',
    }
  }).friend;
  friendRequests.requests.push(newFriend);
  return {
    status: 200,
    data: newFriend,
  };
});

// 获取群聊请求列表
Mock.mock(`${baseURL}/contactList/groupRequests`, 'get', () => {
  return { data: groupRequests.requests };
});
// 处理群聊邀请请求
Mock.mock(`${baseURL}/contactList/groupRequests/invitationPend`, 'post', (options) => {
  const { account_id, group_id, accept } = JSON.parse(options.body);
  const request = groupRequests.requests.find(req => req.sender_id === account_id && req.group_id === group_id);
  if (request) {
    request.status = accept ? 'accepted' : 'rejected';
  }
  return {
    status: 200,
    data: request,
  };
});
// 处理群聊申请请求
Mock.mock(`${baseURL}/contactList/groupRequests/applyPend`, 'post', (options) => {
  const { account_id, group_id, accept } = JSON.parse(options.body);
  const request = groupRequests.requests.find(req => req.sender_id === account_id && req.group_id === group_id);
  if (request) {
    request.status = accept ? 'accepted' : 'rejected';
  }
  return {
    status: 200,
    data: request,
  };
});
// 发送群聊申请请求
Mock.mock(`${baseURL}/contactList/groupRequests/addGroup`, 'post', (options) => {
  const { group_id, reason } = JSON.parse(options.body);
  const newRequest = Mock.mock({
    'request': {
      'apply_id': groupRequests.requests.length + 1,
      'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
      'group_name': '@name',
      'account_name': '@name',
      'sender_id': '1000000000',
      'receiver_id': '1000000000',
      'group_id': group_id,
      'reason': reason,
      'type': 'groupApply',
      'status': 'pending',
      'time': '@datetime',
    }
  }).request;
  groupRequests.requests.push(newRequest);
  return {
    status: 200,
    data: newRequest,
  };
});

// 黑名单
let blackList = Mock.mock({
  'blackList|5-10': [{
    'account_id|1': /[0-9]{10}/,
    'name': '@name',
    'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
    'signature': '爱拼才会赢',
  }]
});
Mock.mock(`${baseURL}/contactList/blackList`, 'get', () => {
  return { data: blackList.blackList };
});
Mock.mock(`${baseURL}/contactList/blackList/remove`, 'post', (options) => {
  const { account_id } = JSON.parse(options.body);
  let blacks = blackList.blackList;
  const index = blacks.findIndex(item => item.account_id === account_id);
  if (index !== -1) {
    blacks.splice(index, 1);
  }
  return {
    status: 200,
    data: blacks,
  };
});
Mock.mock(`${baseURL}/contactList/blackList/add`, 'post', (options) => {
  const { account_id } = JSON.parse(options.body);
  const newBlack = friends.friends.find(friend => friend.account_id === account_id);
  blackList.blackList.push(newBlack);
  return {
    status: 200,
    data: blackList.blackList,
  };
});

// 好友和群聊
// changeRemark
Mock.mock(`${baseURL}/contactList/remark`, 'post', (options) => {
  const { id, remark } = JSON.parse(options.body);
  const friend = friends.friends.find(friend => friend.account_id === id);
  if (friend) {
    friend.remark = remark;
  }
  return {
    status: 200,
    data: friend,
  };
});

let divideList = Mock.mock({
  'divides': {
    'divides': ["家人", "朋友", "同事"]
  }
});
Mock.mock(new RegExp(`${baseURL}/contactList/\\w+/divides`), 'get', (options) => {
  const type = options.url.match(/\/contactList\/(.*?)\/divides/)[1];
  return divideList.divides ;
});
Mock.mock(new RegExp(`${baseURL}/contactList/\\w+/divides/create`), 'post', (options) => {
  const { fd_name } = JSON.parse(options.body);
  divideList.divides.divides.push(fd_name);
  return {
    status: 200,
    data: divideList.divides,
  };
});
Mock.mock(new RegExp(`${baseURL}/contactList/\\w+/divides/delete/\\w+`), 'delete', (options) => {
  const fd_name = decodeURIComponent(options.url.match(/\/delete\/(.*)$/)[1]);
  divideList.divides.divides = divideList.divides.divides.filter(divide => divide !== fd_name);
  return addCorsHeaders({
    status: 200,
    data: divideList.divides,
  });
});
Mock.mock(new RegExp(`${baseURL}/contactList/\\w+/divides/rename`), 'post', (options) => {
  const { old_fd_name, new_fd_name } = JSON.parse(options.body);
  const index = divideList.divides.divides.indexOf(old_fd_name);
  if (index !== -1) {
    divideList.divides.divides.splice(index, 1, new_fd_name);
  }
  return {
    status: 200,
    data: divideList.divides,
  };
});
Mock.mock(new RegExp(`${baseURL}/contactList/\\w+/divides/moveIn`), 'post', (options) => {
  const { tid, divide } = JSON.parse(options.body);
  friends.friends.forEach(friend => {
    if (friend.account_id === tid) {
      friend.tag = divide;
    }
  });
  return {
    status: 200,
    data: { message: `好友 ${tid} 已移动到分组 ${divide}` },
  };
});

let friends = Mock.mock({
  'friends|5-10': [{
    'account_id|1': /[0-9]{10}/,
    'remark': '@name',
    'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
    'status': '@pick(["online", "offline"])',
    'signature': () => Mock.Random.csentence(3, 20),
    'tag': '@pick(["家人", "朋友", "同事"])',
  }]
});
Mock.mock(`${baseURL}/contactList/friends`, 'get', () => {
  return { data: friends.friends };
});

let groups = Mock.mock({
  'groups|5-10': [{
    'account_id|1': /[0-9]{10}/,
    'remark': '@name',
    'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
    'status': '@pick(["online", "offline"])',
    'signature': () => Mock.Random.csentence(3, 20),
    'tag': '@pick(["家人", "朋友", "同事"])',
  }]
});
Mock.mock(`${baseURL}/contactList/groups`, 'get', () => {
  return { data: groups.groups };
});

Mock.mock(`${baseURL}/contactList/groups/create`, 'post', (options) => {
  const { name } = JSON.parse(options.body);
  const group = Mock.mock({
    'group': {
      'account_id|1': /[0-9]{10}/,
      'remark': name,
      'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
      'signature': () => Mock.Random.csentence(3, 20),
      'tag': '@pick(["家人", "朋友", "同事"])',
    }
  }).group;
  return {
    status: 200,
    data: group,
  };
});

Mock.mock(new RegExp(`${baseURL}/contactList/groups/\\d+`), 'delete', (options) => {
  const groupId = parseInt(options.url.split('/').pop());
  const groupList = groups.groups;
  const index = groupList.findIndex(group => group.id === groupId);
  if (index !== -1) {
    groupList.splice(index, 1);
  }
  return {
    status: 200,
    data: groupList,
  };
});

let groupInfo = Mock.mock({
  'groupInfo': {
    'group_name': () => Mock.Random.csentence(3, 10),
    'group_owner': '@id',  // 随机生成群主tid
    'introduction': () => Mock.Random.csentence(3, 20),  // 随机生成一句话作为介绍
    'my_group_nickname': '@name',   // 随机生成一个名字作为群昵称
    'members|50': [
      {
        'account_id': '@id',
        'avatar': '@image("200x200", "#50B347", "#FFF", "Avatar")',
        'group_role': 'group_owner',
        'group_nickname': () => Mock.Random.csentence(3, 10),
        'id': '@id',
        'remark': '@name',
        'nickname': '@name',
        'is_banned': true,
      },
    ],
    'my_group_role': '@pick(["group_owner", "group_manager"])',
  }
});
// 模拟 getGroupInfo 接口
Mock.mock(new RegExp(`${baseURL}/contactList/groups/groupInfo/\\d+`), 'get', () => {
  return { data: groupInfo.groupInfo };
});
// changeGroupNickname
Mock.mock(`${baseURL}/contactList/groups/changeNickname`, 'post', (options) => {
  const { group_id, group_nickname } = JSON.parse(options.body);
  return {
    status: 200,
    data: {},
  };
});
