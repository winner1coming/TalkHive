import Mock from 'mockjs';
const baseURL = 'http://localhost:8080';

// chatlist
let chats = Mock.mock({
  'chats|14-19': [{
    'id|1': /[0-9]{10}/,
    'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
    'name': '@name',
    'lastMessage': '@sentence',
    'lastMessageTime': function() {
      const now = new Date();
      const daysAgo = Math.floor(Math.random() * 7); // 最近7天内
      const date = new Date(now.setDate(now.getDate() - daysAgo));
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      return `${year}-${month}-${day} ${hours}:${minutes}`;
    },
    'unreadCount|0-10': 1,
    'tags': function() {
      const tags = Mock.mock({
        'array|1-2': ['@pick(["friend", "group", "pinned", "blocked"])']
      }).array;
      if (this.unreadCount > 0) {
        tags.push('unread');
      }
      return tags;
    },
  }]
});
// 按照 lastMessageTime 排序
chats.chats.sort((a, b) => {
  const timeA = new Date(a.lastMessageTime);
  const timeB = new Date(b.lastMessageTime);
  return timeB - timeA;
});
let messages = Mock.mock({
  'messages|10-15': [{
  'message_id|+1': /[0-9]{10}/,
  'send_account_id|1': ['1','2'],
  'content': ()=>Mock.Random.csentence(3, 20),
  'sender': '@name',
  'create_time': '@time("HH:mm")',
  'type': 'text',
  'avatar':'@image("200x200", "#50B347", "#FFF", "Mock.js")',
  }]
});

Mock.mock(`${baseURL}/chatlist`, 'get', () => {

  return chats.chats;   // 失败时返回 reason: 失败原因
});
Mock.mock(`${baseURL}/chatlist/createChat`, 'post', () => {
  const chat = Mock.mock({
    'chat': {
      'id|1': /[0-9]{10}/,
      'avatar': '@image("200x200", "#50B347", "#FFF", "Mock.js")',
      'name': '@name',
      'lastMessage': '@sentence',
      'lastMessageTime': '@time("HH:mm")',
      'unreadCount|0-10': 1,
      'tags': function() {
        return Mock.mock({
          'array|1-3': ['@pick(["friend", "group", "unread", "pinned", "blocked"])']
        }).array;
      },
    }
  });
  return chat.chat;
});
Mock.mock(new RegExp(`${baseURL}/chatlist/search/\\w+`), 'get', (options) => {
  const keyword = options.url.split('/').pop();
  const filteredChats = chats.chats.filter(chat => chat.name.includes(keyword));
  return {
    status: 200,
    data: filteredChats,
  };
});
Mock.mock(`${baseURL}/chatlist/pin`, 'post', (options) => {
  console.log("api pin");
  const { tid, is_pinned } = JSON.parse(options.body);
  const chat = chats.chats.find(chat => chat.id === tid);
  if (chat) {
    if (is_pinned) {
      chat.tags.push('pinned');
    } else {
      chat.tags = chat.tags.filter(tag => tag !== 'pinned');
    }
  }
  return {
    status: 200,
    data: chat,
  };
});
Mock.mock(`${baseURL}/messages/read`, 'post', (options) => {
  const { tid, is_read } = JSON.parse(options.body);
  const chat = chats.chats.find(chat => chat.id === tid);
  if (chat) {
    if (is_read) {
      chat.tags = chat.tags.filter(tag => tag !== 'unread');
    } else {
      chat.tags.push('unread');
    }
  }
  return {
    status: 200,
    data: chat,
  };
});

Mock.mock(new RegExp(`${baseURL}/chatlist/\\d+`), 'delete', (options) => {
  const tid = parseInt(options.url.split('/').pop());
  chats.chats = chats.chats.filter(chat => chat.id !== tid);
  return {
    status: 200,
    data: null,
  };
});

Mock.mock(`${baseURL}/chatlist/mute`, 'post', (options) => {
  const { tid, is_muted } = JSON.parse(options.body);
  const chat = chats.chats.find(chat => chat.id === tid);
  if (chat) {
    if (is_muted) {
      chat.tags.push('mute');
    } else {
      chat.tags = chat.tags.filter(tag => tag !== 'mute');
    }
  }
  return {
    status: 200,
    data: chat,
  };
});

Mock.mock(`${baseURL}/chatlist/block`, 'post', (options) => {
  const { tid, is_blocked } = JSON.parse(options.body);
  const chat = chats.chats.find(chat => chat.id === tid);
  if (chat) {
    if (is_blocked) {
      chat.tags.push('blocked');
    } else {
      chat.tags = chat.tags.filter(tag => tag !== 'blocked');
    }
  }
  return {
    status: 200,
    data: chat,
  };
});

Mock.mock(new RegExp(`${baseURL}/messages/\\d+`), 'get', (options) => {
  const tid = parseInt(options.url.split('/').pop());
  return {
    status: 200,
    messages: messages.messages,
  };
});

Mock.mock(new RegExp(`${baseURL}/messages/\\d+/send`), 'post', (options) => {
  const tid = parseInt(options.url.split('/')[2]);
  const { content } = JSON.parse(options.body);
  const newMessage = Mock.mock({
    'message_id': Mock.Random.guid(),
    'send_account_id': Mock.Random.integer(1, 100),
    'content': content,
    'sender': '@name',
    'create_time': '@time("HH:mm")',
    'type': 'text',
  });
  return {
    status: 200,
    data: newMessage,
  };
});

let history = Mock.mock({
  'history|10-15': [{
    'message_id': Mock.Random.guid(),
    'send_account_id': Mock.Random.integer(1, 100),
    'content': ()=>Mock.Random.csentence(3, 20),
    'sender': '@name',
    'create_time': '@time("HH:mm")',
    'type': 'text',
    'avatar':'@image("200x200", "#50B347", "#FFF", "Mock.js")',
  }]
});
Mock.mock(`${baseURL}/messages/history`, 'post', (options) => {
  const tid = parseInt(JSON.parse(options.body).tid);
  return history.history;
});