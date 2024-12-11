import Mock from 'mockjs';
const baseURL = 'http://your-api-url.com';

const friendRequests = Mock.mock({
  'requests|5-10': [{
  'apply_id|+1': 1,
  'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
  'name': '@name',
  'sender_id|1-100': 1,
  'receiver_id|1-100': 1,
  'reason': '@sentence',
  'status': '@pick(["pending", "accepted", "rejected"])',
  'time': '@datetime',
  }]
});

const groupRequests = Mock.mock({
  'requests|5-10': [{
  'apply_id|+1': 1,
  'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
  'group_name': '@name',
  'account_name': '@name',
  'sender_id|1-100': 1,
  'receiver_id|1-100': 1,
  'group_id|1-100': 1,
  'reason': '@sentence',
  'type': '@pick(["groupInvitation", "groupApply", "notification"])',
  'status': '@pick(["pending", "accepted", "rejected", "waiting", "notification"])',
  'time': '@datetime',
  }]
});

Mock.mock(`${baseURL}/contactList/friendRequests`, 'get', () => {
  return friendRequests.requests;
});

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

Mock.mock(`${baseURL}/contactList/groupRequests`, 'get', () => {
  return groupRequests.requests;
});

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

Mock.mock(`${baseURL}/contactList/blackList`, 'get', () => {
  return Mock.mock({
        'blackList|5-10': [{
        'id|1': /[0-9]{10}/,
        'name': '@name',
        'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
        'signature': '爱拼才会赢',
      }]
    }).blackList;
});

Mock.mock(`${baseURL}/contactList/blackList/remove`, 'post', (options) => {
  const { account_id } = JSON.parse(options.body);
  const blackList = Mock.mock({
    'blackList|5-10': [{
      'id|1': /[0-9]{10}/,
      'name': '@name',
      'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
    }]
  }).blackList;
  const index = blackList.findIndex(item => item.id === account_id);
  if (index !== -1) {
    blackList.splice(index, 1);
  }
  return {
    status: 200,
    data: blackList,
  };
});

Mock.mock(`${baseURL}/contactList/friends`, 'get', () => {
  return Mock.mock({
      'friends|5-10': [{
      'account_id|1': /[0-9]{10}/,
      'remark': '@name',
      'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
      'status': '@pick(["online", "offline"])',
      'signature': '@sentence',
      'tag': '@pick(["家人", "朋友", "同事"])',
    }]
    }).friends;
});

Mock.mock(`${baseURL}/contactList/groups`, 'get', () => {
  return Mock.mock({
      'groups|5-10': [{
        'account_id|1': /[0-9]{10}/,
        'remark': '@name',
        'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
        'status': '@pick(["online", "offline"])',
        'signature': '@sentence',
        'tag': '@pick(["家人", "朋友", "同事"])',
      }]
    }).groups;
});

Mock.mock(`${baseURL}/contactList/groups`, 'post', (options) => {
  const { name } = JSON.parse(options.body);
  const group = Mock.mock({
    'group': {
      'account_id|1': /[0-9]{10}/,
      'remark': name,
      'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
      'signature': '@sentence',
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
  const groups = Mock.mock({
    'groups|5-10': [{
      'id|1': /[0-9]{10}/,
      'name': '@name',
      'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
      'status': '@pick(["online", "offline"])',
      'signature': '@sentence',
      'tag': '@pick(["家人", "朋友", "同事"])',
    }]
  }).groups;
  const index = groups.findIndex(group => group.id === groupId);
  if (index !== -1) {
    groups.splice(index, 1);
  }
  return {
    status: 200,
    data: groups,
  };
});